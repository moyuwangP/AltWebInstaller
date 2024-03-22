package altserver

import (
	"AltWebServer/app/util"
	"archive/zip"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func InstallIPA(ctx context.Context, ipaId string, udid string, removePlugIns bool) error {
	var err error
	path := findIPAPath(ipaId)
	if removePlugIns {
		if path, err = RemovePlugInsFromFile(path); err != nil {
			return errors.Wrap(err, "")
		}
		defer os.Remove(path)
	}
	cmd := exec.CommandContext(
		ctx, util.Config().AltserverPath,
		`--udid`, udid,
		`--appleID`, util.Config().AppleId,
		`--password`, util.Config().Password,
		path,
	)
	if util.Config().AnisetteUrl != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("ALTSERVER_ANISETTE_SERVER=%s", util.Config().AnisetteUrl))
	}

	out, err := cmd.Output()
	fmt.Println(string(out))
	if !strings.Contains(string(out), "Notify: Installation Succeeded") {
		return errors.New("failed")
	}
	return errors.Wrap(err, "")
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
