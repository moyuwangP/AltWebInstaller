package libimobiledevice

import (
	"AltWebServer/app/command"
	"bytes"
	"context"
	"encoding/csv"
	"github.com/pkg/errors"
	"os/exec"
	"strings"
)

const (
	ideviceInstaller = "ideviceinstaller"
	notFoundPrefix   = "No device found with udid"
)

var ideviceInstallerPath string

func init() {
	path, err := command.FindBinary(context.Background(), ideviceInstaller)
	if err != nil {
		panic(errors.Errorf("%s: not found", ideviceInstaller))
	}
	ideviceInstallerPath = path
}

type Package struct {
	BundleIdentifier string
	Version          string
	DisplayName      string
}

func ListInstalledPackages(ctx context.Context, udid string) ([]Package, error) {
	cmd := exec.CommandContext(ctx, ideviceInstallerPath, "-n", "-u", udid, "list")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		if strings.HasPrefix(stderr.String(), notFoundPrefix) {
			return nil, DeviceNotFound
		}
		return nil, errors.Wrapf(err, "failed to execute cmd, %d\n", stderr.String())
	}

	out := strings.NewReader(strings.ReplaceAll(stdout.String(), ", \"", ",\""))
	reader := csv.NewReader(out)
	data, err := reader.ReadAll()
	packages := make([]Package, len(data)-1)
	for i := 1; i < len(data); i++ {
		packages[i-1] = Package{
			BundleIdentifier: data[i][0],
			Version:          data[i][1],
			DisplayName:      data[i][2],
		}
	}
	return packages, nil
}
