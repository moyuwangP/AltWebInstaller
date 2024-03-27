package service

import (
	"AltWebServer/app/command/altserver"
	"AltWebServer/app/command/libimobiledevice"
	"AltWebServer/app/model"
	db2 "AltWebServer/app/model/db"
	"AltWebServer/app/util"
	"context"
	"fmt"
	"howett.net/plist"
	"os"
	"regexp"
	"strings"
	"time"
)

var Device = DeviceService{}

type DeviceService struct {
}

type OnDevicePackage struct {
	libimobiledevice.Package
	IPAHash     string `json:"ipa_hash"`
	RefreshedAt time.Time
}

func (r *DeviceService) AppsInstalled(ctx context.Context, udid string) ([]OnDevicePackage, error) {
	packages, err := libimobiledevice.ListInstalledPackages(ctx, udid)
	if err != nil {
		return nil, err
	}

	var installationRecords []struct {
		db2.Package
		RefreshedAt time.Time
	}
	err = util.DB().Table("installations").Joins("left join packages on packages.md5 = installations.md5").
		Where("udid", udid).
		Select("packages.*, installations.refreshed_at").
		Find(&installationRecords).
		Error
	if err != nil {
		return nil, err
	}

	onDevice := make([]OnDevicePackage, 0)
	for _, record := range installationRecords {
		for _, pack := range packages {
			if strings.HasPrefix(pack.BundleIdentifier, record.CFBundleIdentifier) && pack.Version == record.CFBundleShortVersionString {
				onDevice = append(onDevice, OnDevicePackage{
					Package:     pack,
					RefreshedAt: record.RefreshedAt,
					IPAHash:     record.MD5,
				})
			}
		}
	}

	return onDevice, nil
}

const (
	lockdownDir = "/var/lib/lockdown"
	plistRegex  = "^[0-9A-F]{8}-[0-9A-F]{16}.plist$"
)

func (r *DeviceService) ListPairedDevices() ([]*model.DeviceStatus, error) {
	files, err := os.ReadDir(lockdownDir)
	if err != nil {
		return nil, err
	}

	deviceMap := make([]*model.DeviceStatus, 0)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if ok, _ := regexp.MatchString(plistRegex, file.Name()); !ok {
			continue
		}

		dat, err := os.ReadFile(fmt.Sprintf("%s/%s", lockdownDir, file.Name()))
		if err != nil {
			continue
		}

		var deviceInfo struct{ WiFiMACAddress string }
		if _, err = plist.Unmarshal(dat, &deviceInfo); err != nil {
			util.LogErrorf("failed to read pairing record: %s", file.Name())
		}

		udid := strings.Split(file.Name(), ".")[0]
		deviceMap = append(deviceMap, &model.DeviceStatus{
			UDID:       udid,
			MacAddress: deviceInfo.WiFiMACAddress,
			Online:     false,
		})
	}
	return deviceMap, nil
}

func (r *DeviceService) InstallIPA(ctx context.Context, udid string, IPAHash string, removePlugIns bool) error {
	var (
		err                error
		installationRecord db2.Installation
		info               db2.Package
	)

	if err = util.DB().
		Where("md5", IPAHash).
		First(&info).
		Error; err != nil {
		return err
	}
	installationRecord = db2.Installation{
		UDID:          udid,
		RemovePlugIns: removePlugIns,
		MD5:           IPAHash,
		BundleID:      info.CFBundleIdentifier,
		BundleVersion: info.CFBundleShortVersionString,
	}

	return altserver.InstallIPA(ctx, installationRecord)
}
