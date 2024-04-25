package controller

import (
	"AltWebServer/app/command/altserver"
	"AltWebServer/app/command/libimobiledevice"
	"AltWebServer/app/model"
	db2 "AltWebServer/app/model/db"
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
		err                error
		udid               string
		body               []byte
		request            InstallationRequest
		installationRecord db2.Installation
		info               db2.Package
	)

	udid = ctx.Param("udid")
	if body, err = io.ReadAll(ctx.Request.Body); err != nil {
		d.responseFailAndExit(ctx, http.StatusInternalServerError, err.Error())
	}
	if err = json.Unmarshal(body, &request); err != nil {
		d.responseFailAndExit(ctx, http.StatusInternalServerError, err.Error())
	}

	if err = util.DB().
		Where("md5", request.Package).
		First(&info).
		Error; errors.Is(err, gorm.ErrRecordNotFound) {
		d.responseFailAndExit(ctx, http.StatusNotFound, fmt.Sprintf("ipa with hash '%s' not found", request.Package))
	} else if err != nil {
		d.responseFailAndExit(ctx, http.StatusInternalServerError, err.Error())

	}
	installationRecord = db2.Installation{
		UDID:          udid,
		RemovePlugIns: request.RemovePlugIns,
		MD5:           request.Package,
		BundleID:      info.CFBundleIdentifier,
		BundleVersion: info.CFBundleShortVersionString,
	}

	cmd, err := altserver.InstallIPAStream(ctx, installationRecord)
	if err != nil {
		d.responseFailAndExit(ctx, http.StatusInternalServerError, err.Error())
	}
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")

	ctx.Stream(func(w io.Writer) bool {
		if cmd.HasNextLine() {
			ctx.SSEvent("message", cmd.NextLine())
			return true
		}

		if err := cmd.Finish(); err != nil {
			ctx.SSEvent("message", "failed")
			ctx.SSEvent("message", err.Error())
		} else {
			ctx.SSEvent("message", "ok")
		}
		return false
	})
}

func (d DeviceController) UninstallPackage(ctx *gin.Context) {
	// TODO remove package from device
	ctx.JSON(http.StatusNotImplemented, map[string]string{})
}
