package service

import (
	"fmt"

	"github.com/CRORCR/call/app/model"
	"github.com/gin-gonic/gin"
)

type Service struct {
}

// CallPrice 获取主播私聊价格
func (s *Service) CallPrice(ctx *gin.Context, uid int64) *model.CallPriceResp {
	resp := &model.CallPriceResp{
		PriceCoins: make(map[int64]int64),
	}
	resp.PriceCoins[uid] = 12

	return resp
}

// CallPrice 获取主播私聊价格
func (s *Service) CallPriceList(ctx *gin.Context, uids []int64) *model.CallPriceResp {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("居然有错误", err)
		}
	}()
	fmt.Println("----1-----")
	resp := &model.CallPriceResp{
		PriceCoins: make(map[int64]int64),
	}
	fmt.Println("----2-----")
	for _, uid := range uids {
		resp.PriceCoins[uid] = 1234
		fmt.Println("----3-----")
	}

	return resp
}
