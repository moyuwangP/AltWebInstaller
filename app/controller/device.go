package controller

import (
	"AltWebServer/app/command/libimobiledevice"
	"AltWebServer/app/model"
	"AltWebServer/app/service"
	"AltWebServer/app/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
	"io"
	"net/http"
)

type DeviceController struct {
	trait
}

var Device = DeviceController{}

// List all available device
func (d DeviceController) List(ctx *gin.Context) {
	devices, err := service.Device.ListPairedDevices()
	if err != nil {
		d.responseFailAndExit(ctx, http.StatusInternalServerError, err.Error())
	}
	deviceMap := util.KeyBy(devices, func(v *model.DeviceStatus) string {
		return v.UDID
	})

	onlineDeviceList, err := libimobiledevice.ListOnlineDevices(ctx.Request.Context())
	if err != nil {
		util.LogError("unable to find online devices, set all device's online status to offline")
	}

	for _, onlineDevice := range onlineDeviceList {
		if device, ok := deviceMap[onlineDevice]; ok {
			device.Online = true
		}
	}

	ctx.JSON(http.StatusOK, maps.Values(deviceMap))
}

func (d DeviceController) ListAppsOnDevice(ctx *gin.Context) {
	udid := ctx.Param("udid")
	packages, err := service.Device.AppsInstalled(ctx.Request.Context(), udid)
	if errors.Is(err, libimobiledevice.DeviceNotFound) {
		d.responseFailAndExit(ctx, http.StatusNotFound, "device not found, check if device is online: "+err.Error())
	} else if err != nil {
		d.responseFailAndExit(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, packages)
}

type InstallationRequest struct {
	Package       string `json:"package"`
	RemovePlugIns bool   `json:"remove_plug_ins"`
}

func (d DeviceController) InstallApp(ctx *gin.Context) {
	var (
		err     error
		udid    string
		body    []byte
		request InstallationRequest
	)

	udid = ctx.Param("udid")
	if body, err = io.ReadAll(ctx.Request.Body); err != nil {
		d.responseFailAndExit(ctx, http.StatusInternalServerError, err.Error())
	}
	if err = json.Unmarshal(body, &request); err != nil {
		d.responseFailAndExit(ctx, http.StatusInternalServerError, err.Error())
	}

	if err = service.Device.InstallIPA(
		ctx.Request.Context(), udid,
		request.Package, request.RemovePlugIns,
	); errors.Is(err, gorm.ErrRecordNotFound) {
		d.responseFailAndExit(ctx, http.StatusNotFound, fmt.Sprintf("ipa with hash '%s' not found", request.Package))
	} else if err != nil {
		d.responseFailAndExit(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (d DeviceController) UninstallPackage(ctx *gin.Context) {
	// TODO remove package from device
	ctx.JSON(http.StatusNotImplemented, map[string]string{})
}
