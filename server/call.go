package server

import (
	"fmt"
	"strconv"

	"github.com/CRORCR/call/model"
	"github.com/CRORCR/call/service"
	"github.com/CRORCR/duoo-common/code"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserController struct {
	controllerBase
	svc service.Service
}

var UserServer = &UserController{
	svc: service.Service{},
}

// CallPrice get请求参数获取
func (u *UserController) CallPrice(ctx *gin.Context) {
	uid, _ := strconv.ParseInt(ctx.Query("uid"), 10, 64)

	if uid == 0 {
		fmt.Println("参数错误")
		u.ResponseError(ctx, code.RequestParamError)
		return
	}

	// 查询缓存聊天价格
	resp := u.svc.CallPrice(ctx, uid)
	u.ResponseOk(ctx, resp)
}

// get数组获取
func (u *UserController) CallPriceUids(ctx *gin.Context) {
	var req model.CallPriceReq
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		fmt.Println("参数错误", err)
		u.ResponseError(ctx, code.RequestParamError)
		return
	}
	// 查询缓存聊天价格
	resp := u.svc.CallPriceList(ctx, req.Uids)
	u.ResponseOk(ctx, resp)
}
