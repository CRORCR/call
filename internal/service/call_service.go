package service

import (
	"fmt"
	"github.com/CRORCR/call/internal/config"
	"github.com/CRORCR/call/internal/dao"
	"github.com/CRORCR/call/internal/grpc"
	"github.com/CRORCR/call/internal/model"
	"github.com/gin-gonic/gin"
	"time"
)

type UserService struct {
	conf    *config.Configuration
	rpc     *grpc.RpcService
	callDao dao.CallRepository
}

func NewUserService(conf *config.Configuration, rpcService *grpc.RpcService, callDao dao.CallRepository) *UserService {
	return &UserService{
		conf:    conf,
		rpc:     rpcService,
		callDao: callDao,
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

	//s.callDao.GetDialCallById(123)

	uuid, ok := s.callDao.Lock("hello")
	if !ok {
		fmt.Println("req limit")
		return resp
	}
	time.Sleep(time.Second * 2)
	defer s.callDao.UnLock("hello", uuid)

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
