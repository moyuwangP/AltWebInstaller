package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type trait struct {
}

var ErrQuit = errors.New("not found")

func (receiver trait) responseFailAndExit(ctx *gin.Context, code int, msg string) {
	ctx.AbortWithStatusJSON(code, map[string]string{"error_msg": msg})
	// todo: log err
	panic(ErrQuit)
}
