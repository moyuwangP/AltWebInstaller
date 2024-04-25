package controller

import (
	db2 "AltWebServer/app/model/db"
	"AltWebServer/app/util"
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/aofei/mimesniffer"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"howett.net/plist"
	"io"
	"net/http"
	"os"
	"regexp"
)

type PackageController struct {
	trait
}

var Package = PackageController{}

const (
	infoPlistPathRegex = `^Payload/.+\.app/Info\.plist$`
	appexPathRegex     = `^Payload/.+\.app/PlugIns/.+\.appex/$`
	appIconRegex       = `^Payload/.+\.app/%s.+\.png$`
	tmpIpa             = "ipa/%s"
	appIcon            = "app_icon/%s"
)

func (p PackageController) Upload(ctx *gin.Context) {
	var (
		fileId string
		err    error
		info   db2.Package
	)

	if fileId, err = saveIPA(ctx); err != nil {
		p.responseFailAndExit(ctx, http.StatusInternalServerError, err.Error())
	}

	if info, err = sniffIPAInfo(fileId); errors.Is(err, zip.ErrFormat) {
		os.Remove(fmt.Sprintf(tmpIpa, fileId))
		p.responseFailAndExit(ctx, http.StatusBadRequest, err.Error())
	} else if err != nil {
		os.Remove(fmt.Sprintf(tmpIpa, fileId))
		p.responseFailAndExit(ctx, http.StatusInternalServerError, err.Error())
	}

	if err = util.DB().
		Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "md5"}}}).
		Create(&info).
		Error; err != nil {
		os.Remove(fmt.Sprintf(tmpIpa, fileId))
		p.responseFailAndExit(ctx, http.StatusInternalServerError, "unable to save ipa")
	}

	ctx.JSON(http.StatusOK, info)
}

func (p PackageController) Get(ctx *gin.Context) {
	hash := ctx.Param("hash")
	ipa := db2.Package{}
	if err := util.DB().Where("md5", hash).First(&ipa).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		p.responseFailAndExit(ctx, http.StatusNotFound, "no such ipa package")
	} else if err != nil {
		p.responseFailAndExit(ctx, http.StatusInternalServerError, "unknown error occurred")
	}
	ctx.JSON(http.StatusOK, ipa)
}

func (p PackageController) GetAppIcon(ctx *gin.Context) {
	var (
		err error
	)
	hash := ctx.Param("hash")
	data, err := os.ReadFile(fmt.Sprintf(appIcon, hash))
	if err != nil {
		p.responseFailAndExit(ctx, http.StatusNotFound, err.Error())
	}

	ctx.Header("cache-control", "max-age=3600")
	ctx.Data(http.StatusOK, mimesniffer.Sniff(data), data)
}

func saveIPA(ctx *gin.Context) (string, error) {
	var (
		ipaBytes []byte
		uid      string
		err      error
	)
	if ipaBytes, err = io.ReadAll(ctx.Request.Body); err != nil {
		return "", errors.Wrap(err, "failed to read body")
	}
	hash := md5.Sum(ipaBytes)
	uid = hex.EncodeToString(hash[:])
	if err = os.MkdirAll("./ipa", os.ModePerm); err != nil {
		return "", err
	}
	if err = os.WriteFile(fmt.Sprintf(tmpIpa, uid), ipaBytes, os.ModePerm); err != nil {
		return "", err
	}

	return uid, nil
}

func sniffIPAInfo(md5 string) (info db2.Package, err error) {
	var (
		reader io.ReadCloser
		bytes  []byte
		infoP  infoPlist
	)

	archive, err := zip.OpenReader(fmt.Sprintf(tmpIpa, md5))
	if err != nil {
		return info, errors.Wrap(err, "unable to read zip")
	}
	defer archive.Close()

	infoFound := false
	for _, f := range archive.File {
		if ok, _ := regexp.MatchString(infoPlistPathRegex, f.Name); !ok {
			continue
		}
		infoFound = true

		if reader, err = f.Open(); err != nil {
			return
		}
		if bytes, err = io.ReadAll(reader); err != nil {
			return
		}

		if _, err = plist.Unmarshal(bytes, &infoP); err != nil {
			return
		}
		break
	}
	if !infoFound {
		return info, errors.New("malformed ipa: info plist not found")
	}

	info = db2.Package{
		CFBundleName:               infoP.CFBundleName,
		CFBundleShortVersionString: infoP.CFBundleShortVersionString,
		CFBundleIdentifier:         infoP.CFBundleIdentifier,
		ContainsPlugIn:             false,
		MD5:                        md5,
	}

	for _, f := range archive.File {
		if ok, _ := regexp.MatchString(appexPathRegex, f.Name); ok {
			info.ContainsPlugIn = true
			break
		}
	}

	if len(infoP.CFBundleIcons.CFBundlePrimaryIcon.CFBundleIconFiles) > 0 {
		iconFile := infoP.CFBundleIcons.CFBundlePrimaryIcon.CFBundleIconFiles[0]
		for _, f := range archive.File {
			if ok, _ := regexp.MatchString(fmt.Sprintf(appIconRegex, iconFile), f.Name); ok {
				if reader, err = f.Open(); err != nil {
					util.LogError(err.Error())
					break
				}
				if bytes, err = io.ReadAll(reader); err != nil {
					util.LogError(err.Error())
					break
				}
				if err = os.MkdirAll("./app_icon", os.ModePerm); err != nil {
					util.LogError(err.Error())
					break
				}
				if err = os.WriteFile(fmt.Sprintf(appIcon, md5), bytes, os.ModePerm); err != nil {
					util.LogError(err.Error())
					break
				}
				break
			}
		}
	}

	return
}

type infoPlist struct {
	CFBundleName               string `json:"cf_bundle_name"`
	CFBundleShortVersionString string `json:"cf_bundle_short_version_string"`
	CFBundleIdentifier         string `json:"cf_bundle_identifier"`
	CFBundleIcons              struct {
		CFBundlePrimaryIcon struct {
			CFBundleIconFiles []string `json:"CFBundleIconFiles"`
			CFBundleIconName  string   `json:"CFBundleIconName"`
		} `json:"CFBundlePrimaryIcon"`
	} `json:"CFBundleIcons"`
}
