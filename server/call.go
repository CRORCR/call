package server

import (
	"fmt"
	"strconv"

	"github.com/CRORCR/call/app/model"
	"github.com/CRORCR/call/service"
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
func (u *UserController) CallPrice(c *gin.Context) {
	uid, _ := strconv.ParseInt(c.Query("uid"), 10, 64)

	if uid == 0 {
		fmt.Println("参数错误")
		u.ResponseError(c, 10001)
		return
	}

	// 查询缓存聊天价格
	resp := u.svc.CallPrice(c, uid)
	u.ResponseOk(c, resp)
}

// get数组获取
func (u *UserController) CallPriceUids(c *gin.Context) {
	var req model.CallPriceReq
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		fmt.Println("参数错误", err)
		u.ResponseError(c, 10001)
		return
	}
	// 查询缓存聊天价格
	resp := u.svc.CallPriceList(c, req.Uids)
	u.ResponseOk(c, resp)
}
