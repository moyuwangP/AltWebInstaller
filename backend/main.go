package main

import (
	"AltWebServer/app/controller"
	"AltWebServer/app/util"
	"AltWebServer/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {

	r := gin.Default()
	router.Setup(r)

	if util.Config().AutoRefresh != "" {
		c := cron.New()
		if _, err := c.AddFunc("* * * * *", controller.DoRefresh); err != nil {
			panic(err)
		}
		if _, err := c.AddFunc(
			util.Config().AutoRefresh,
			controller.ScheduleRefresh,
		); err != nil {
			panic(err)
		}
		go c.Run()
	}

	_ = r.Run(fmt.Sprintf(":%s", util.Config().Port))
}
