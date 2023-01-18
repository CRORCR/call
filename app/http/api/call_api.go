package api

import (
	"fmt"
	"strconv"

	"github.com/CRORCR/call/internal/model"
	"github.com/CRORCR/call/internal/service"
	"github.com/CRORCR/duoo-common/code"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	controllerBase
	svc *service.UserService
}

var UserServer *UserController

func NewUserController(svc *service.UserService) *UserController {
	UserServer = &UserController{svc: svc}
	return UserServer
}

// CallPrice get请求参数获取
func (u *UserController) CallPrice(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithFields(logrus.Fields{
				"service": "UserController.CallPrice",
				"error":   "panic",
			}).Error(err)
		}
	}()

	uid, _ := strconv.ParseInt(ctx.Query("uid"), 10, 64)

	logrus.WithFields(logrus.Fields{
		"level": "2",
		"test":  "test",
	}).Info("hello-trace")

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
