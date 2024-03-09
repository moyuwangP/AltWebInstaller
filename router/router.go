package router

import (
	"AltWebServer/app/controller"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

type api struct {
	url        string
	method     string
	handler    Handler
	param      interface{}
	middleware []gin.HandlerFunc
}

type group struct {
	apis        []api
	middlewares []gin.HandlerFunc
}

// Setup setup http server router
func Setup(engine *gin.Engine) {
	engine.Use(middlewares...)
	for groupPath, configs := range apis {
		group := engine.Group(groupPath)
		group.Use(configs.middlewares...)
		for _, c := range configs.apis {
			config := c
			group.Handle(config.method, config.url, append(config.middleware, func(c *gin.Context) {
				call(c, config.handler, config.param)
			})...)
		}
	}
}

type Handler func(ctx *gin.Context)

type Middleware Handler
type stackTracer interface {
	StackTrace() errors.StackTrace
}

func call(ctx *gin.Context, function Handler, params interface{}) {
	defer func() {
		if r := recover(); r != nil {
			if r == controller.ErrQuit {
				return
			}

			fmt.Println("Recovered in f", r)
			err, ok := r.(stackTracer) // ok is false if errors doesn't implement stackTracer
			if !ok {
				ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})

			}

			stack := err.StackTrace()
			fmt.Println(stack)
			ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
	}()
	//req := NewReq(r)

	//if params != nil {
	//	t := reflect.TypeOf(params)
	//	valueType := t.Elem()            // 得到结构体对象的类型
	//	newPtr := reflect.New(valueType) // 产生指向此结构体类型的指针
	//	params = newPtr.Interface()
	//
	//	validateParams(ctx, params)
	//}

	function(ctx)
	//if err == nil && response == nil {
	//} else if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	//} else {
	//	ctx.JSON(http.StatusOK, response)
	//}
}

func validateParams(ctx *gin.Context, params interface{}) {
	jsonData, err := io.ReadAll(ctx.Request.Body)
	queries := ctx.Request.URL.Query()
	pMap := map[string]interface{}{}

	for key, query := range queries {
		pMap[key] = query[0]
	}

	for _, p := range ctx.Params {
		pMap[p.Key] = p.Value
	}

	err = json.Unmarshal(jsonData, params)
	j, err := json.Marshal(pMap)
	err = json.Unmarshal(j, params)
	err = validator.New().Struct(params)
	if err != nil {
		// return panic info
	}
}
