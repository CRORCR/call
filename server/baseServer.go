package server

import (
	"net/http"

	"github.com/CRORCR/call/app/model/base"
	"github.com/gin-gonic/gin"
)

// todo 后面需要支持多语言

type controllerBase struct {
}

// response: {"error_code":10001,"error_message":"参数错误","succeed":false,"data":null}
func (*controllerBase) ResponseError(ctx *gin.Context, errorCode int64) {
	// 根据code查询对应msg
	msg := ""
	result := base.Response{
		Succeed:      false,
		ErrorCode:    errorCode,
		ErrorMessage: msg,
	}
	ctx.JSON(http.StatusOK, result)
}

// response: {"error_code":0,"error_message":"","succeed":true,"data":{"price_coins":{"11":1234,"22":1234}}}
func (*controllerBase) ResponseOk(ctx *gin.Context, data interface{}) {
	result := base.Response{
		Succeed:   true,
		ErrorCode: 0,
		Data:      data,
	}
	ctx.JSON(http.StatusOK, result)
}
