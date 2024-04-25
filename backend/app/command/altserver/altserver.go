package altserver

import (
	db2 "AltWebServer/app/model/db"
	"AltWebServer/app/util"
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type InstallCommand struct {
	*exec.Cmd
	InstallParams db2.Installation
	stdout        string
	scanner       Scanner
	path          string
}

type Scanner interface {
	Scan() bool
	Text() string
	Split(split bufio.SplitFunc)
}

var ErrorInstallFailed = errors.New("failed")

func (cmd *InstallCommand) HasNextLine() bool {
	return cmd.scanner.Scan()
}

func (cmd *InstallCommand) NextLine() string {
	line := cmd.scanner.Text()
	util.LogDebug(line)
	cmd.stdout += fmt.Sprintf("%s\n", line)
	return line
}

func (cmd *InstallCommand) Finish() error {
	if cmd.InstallParams.RemovePlugIns {
		defer os.Remove(cmd.path)
	}
	for cmd.HasNextLine() {
		cmd.NextLine()
	}
	if err := cmd.Wait(); err != nil {
		return errors.Wrap(err, "unable to finish command")
	}
	if !strings.Contains(cmd.stdout, "Notify: Installation Succeeded") {
		return ErrorInstallFailed
	}

	return util.DB().
		Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "udid"}, {Name: "bundle_id"}, {Name: "bundle_version"}},
				DoUpdates: clause.AssignmentColumns([]string{"md5", "refreshed_at"}),
			},
		).
		Create(&cmd.InstallParams).
		Error
}

func InstallIPAStream(ctx context.Context, installation db2.Installation) (*InstallCommand, error) {
	var err error
	path := findIPAPath(installation.MD5)
	if installation.RemovePlugIns {
		if path, err = RemovePlugInsFromFile(path); err != nil {
			return nil, errors.Wrap(err, "")
		}
	}
	cmd := exec.CommandContext(
		ctx, util.Config().AltserverPath,
		`--udid`, installation.UDID,
		`--appleID`, util.Config().AppleId,
		`--password`, util.Config().Password,
		path,
	)
	if util.Config().AnisetteUrl != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("ALTSERVER_ANISETTE_SERVER=%s", util.Config().AnisetteUrl))
	}

	installation.RefreshedAt = time.Now()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "unable to open stdout")
	}

	if err = cmd.Start(); err != nil {
		return nil, errors.Wrap(err, "unable to start command")
	}
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)

	return &InstallCommand{Cmd: cmd, InstallParams: installation, scanner: scanner, path: path}, nil
}

func InstallIPA(ctx context.Context, installation db2.Installation) error {
	cmd, err := InstallIPAStream(ctx, installation)
	if err != nil {
		return err
	}
	return cmd.Finish()
}

func findIPAPath(md5 string) string {
	return fmt.Sprintf("./ipa/%s", md5)
}

const appexPathRegex = `^Payload/.+\.app/PlugIns/.+\.appex/$`
const appexContentPathRegex = `^Payload/.+\.app/PlugIns/.+\.appex/.+$`

func RemovePlugInsFromFile(path string) (string, error) {
	id := uuid.New().String()
	if err := unzipFile(path, id); err != nil {
		return "", err
	}

	return fmt.Sprintf("./tmp/%s.ipa", id), nil
}

func unzipFile(path string, id string) error {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return errors.Wrap(err, "")
	}
	if err = os.MkdirAll("./tmp", os.ModePerm); err != nil {
		return err
	}

	newArchive, err := os.Create(fmt.Sprintf("./tmp/%s.ipa", id))
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer archive.Close()
	defer newArchive.Close()

	writer := zip.NewWriter(newArchive)
	defer writer.Close()

	for _, f := range archive.File {
		appex, _ := regexp.MatchString(appexPathRegex, f.Name)
		appexContent, _ := regexp.MatchString(appexContentPathRegex, f.Name)

		if appex || appexContent {
			continue
		}

		if f.FileInfo().IsDir() {
			continue
		}
		header := f.FileHeader
		w, err := writer.CreateHeader(&header)
		if err != nil {
			return errors.Wrap(err, "")
		}
		data, err := f.Open()
		if err != nil {
			return errors.Wrap(err, "")
		}
		if _, err := io.Copy(w, data); err != nil {
			return errors.Wrap(err, "")
		}
	}
	return nil
}
