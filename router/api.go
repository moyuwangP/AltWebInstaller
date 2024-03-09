package router

import (
	"AltWebServer/app/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

var middlewares = []gin.HandlerFunc{
	//middleware.CatchPanic,
	//middleware.CORSHandler,
}

var apis = map[string]group{
	"api": {
		apis: []api{
			{
				url:     "devices",
				handler: controller.Device.List,
				method:  http.MethodGet,
			},
		},
	},
	"api/devices/:udid": {
		apis: []api{
			{
				url:     "",
				handler: controller.Device.ListAppsOnDevice,
				method:  http.MethodGet,
			},
			{
				url:     "",
				handler: controller.Device.InstallApp,
				method:  http.MethodPost,
				//param:   model.InstallConfig{},
			},
			//{
			//	url:     ":bundle",
			//	handler: controller.Package.Uninstall,
			//	method:  http.MethodDelete,
			//},
		},
	},
	"api/packages": {
		apis: []api{
			{
				url:     "",
				handler: controller.Package.Upload,
				method:  http.MethodPost,
			},
			{
				url:     ":hash",
				handler: controller.Package.Get,
				method:  http.MethodGet,
			},
		},
	},
	"api/udid": {
		apis: []api{
			{
				url:     "",
				handler: controller.UDID.Register,
				method:  http.MethodPost,
			},
			{
				url:     "",
				handler: controller.UDID.GetUDID,
				method:  http.MethodGet,
			},
		},
		middlewares: []gin.HandlerFunc{},
	},
}
