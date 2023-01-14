package server

import (
	"net/http"

	"github.com/CRORCR/call/model/base"
	"github.com/CRORCR/duoo-common/code"
	"github.com/gin-gonic/gin"
)

// todo 后面需要支持多语言

type controllerBase struct {
}

// response: {"error_code":100101,"error_message":"Params error","succeed":false,"data":null}
func (*controllerBase) ResponseError(ctx *gin.Context, err error) {
	// 根据code查询对应msg
	e := code.Cause(err)

	result := base.Response{
		Succeed:      false,
		ErrorCode:    e.Code(),
		ErrorMessage: e.Message(),
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
