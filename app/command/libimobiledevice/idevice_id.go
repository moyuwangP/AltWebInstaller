package libimobiledevice

import (
	"AltWebServer/app/command"
	"context"
	"github.com/pkg/errors"
	"os/exec"
	"strings"
)

const (
	ideviceId     = "idevice_id"
	networkSuffix = " (Network)"
)

var ideviceIdPath string

func init() {
	path, err := command.FindBinary(context.Background(), ideviceId)
	if err != nil {
		panic(errors.Errorf("%s: not found", ideviceId))
	}
	ideviceIdPath = path
}

func ListOnlineDevices(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, ideviceIdPath, "-l", "-n")
	onlineDevices, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute cmd")
	}
	onlineDeviceList := strings.Split(
		strings.ReplaceAll(
			strings.TrimSpace(string(onlineDevices)),
			networkSuffix, ""),
		"\n")
	return onlineDeviceList, nil
}
