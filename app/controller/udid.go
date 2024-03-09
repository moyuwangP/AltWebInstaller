package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/tidwall/buntdb"
	"howett.net/plist"
	"io"
	"net/http"
	"strings"
	"time"
)

type UDIDController struct {
	trait
}

var UDID = UDIDController{}

var db *buntdb.DB = nil

func init() {
	var err error
	db, err = buntdb.Open(":memory:")
	if err != nil {
		panic(err)
	}
}

func (receiver *UDIDController) Register(ctx *gin.Context) {
	// TODO: record UDID
	bytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			map[string]string{
				"error": "unable to read plist",
			})
		return
	}

	// Extract XML from data
	plistBegin := `<?xml version="1.0"`
	plistEnd := `</plist>`

	pos1 := strings.Index(string(bytes), plistBegin)
	pos2 := strings.Index(string(bytes), plistEnd)
	if pos1 == -1 || pos2 == -1 {
		receiver.responseFailAndExit(ctx, http.StatusBadRequest, "plist not found")
	}
	plistData := bytes[pos1:pos2]

	dic := map[string]string{}
	_, err = plist.Unmarshal(plistData, &dic)
	if err != nil {
		receiver.responseFailAndExit(ctx, http.StatusBadRequest, "plist unmarshal error")
	}

	udid, ok := dic["UDID"]
	if !ok {
		receiver.responseFailAndExit(ctx, http.StatusBadRequest, "udid not provided")
	}

	addr := ctx.ClientIP()

	err = db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(addr, udid, &buntdb.SetOptions{
			Expires: true,
			TTL:     time.Hour,
		})
		return err
	})

	if err != nil {
		receiver.responseFailAndExit(ctx, http.StatusInternalServerError, "failed to record udid")
	}

	ctx.Redirect(301, "/enrolled/"+udid)
}

func (receiver *UDIDController) GetUDID(ctx *gin.Context) {
	var udid string
	err := db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(ctx.ClientIP())
		if err != nil {
			return err
		}
		udid = val
		return nil
	})
	if errors.Is(err, buntdb.ErrNotFound) {
		receiver.responseFailAndExit(ctx, http.StatusBadRequest, "udid not registered")
	} else if err != nil {
		receiver.responseFailAndExit(ctx, http.StatusInternalServerError, err.Error())
	}
	ctx.JSON(http.StatusOK, map[string]string{"udid": udid})
}
