package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func CatchPanic(c *gin.Context) {
	fmt.Println("in middleware")
	fmt.Println(c.Request.URL.Query())
	defer func() {
		if err := recover(); err != nil && err != "exit" {
			c.JSON(500,
				gin.H{
					"err_msg":  "internal exception occurred",
					"err_code": "1",
					"data":     err,
				})
		}
		c.Abort()
	}()
	c.Next()
	fmt.Println("finish panic mid")
}

func CORSHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Next()
}
