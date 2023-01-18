package service

import (
	"fmt"

	"github.com/CRORCR/call/internal/config"
	"github.com/CRORCR/call/internal/grpc"
	"github.com/CRORCR/call/internal/model"
	"github.com/gin-gonic/gin"
)

type UserService struct {
	conf *config.Configuration
	rpc  *grpc.RpcService
}

func NewUserService(conf *config.Configuration, rpcService *grpc.RpcService) *UserService {
	return &UserService{
		conf: conf,
		rpc:  rpcService,
	}
}

// CallPrice 获取主播私聊价格
func (s *UserService) CallPrice(ctx *gin.Context, uid int64) *model.CallPriceResp {
	resp := &model.CallPriceResp{
		PriceCoins: make(map[int64]int64),
	}
	resp.PriceCoins[uid] = 12

	result, err := s.rpc.GetTransferLogResult(ctx, uid)
	fmt.Println("打印结果", result, err)

	return resp
}

// CallPrice 获取主播私聊价格
func (s *UserService) CallPriceList(ctx *gin.Context, uids []int64) *model.CallPriceResp {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("居然有错误", err)
		}
	}()
	resp := &model.CallPriceResp{
		PriceCoins: make(map[int64]int64),
	}
	for _, uid := range uids {
		resp.PriceCoins[uid] = 1234
	}

	return resp
}
