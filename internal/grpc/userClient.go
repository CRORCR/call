package grpc

import (
	context "context"
	"fmt"
	callPrice "github.com/CRORCR/duoo-common/proto/call_price"
)

func (r *RpcService) GetTransferLogResult(ctx context.Context, uid int64) (*callPrice.GetPriceResp, error) {
	var orderDetail *callPrice.GetPriceResp
	var err error
	if grpcClient, err := r.GetUserClient(); err == nil {
		req := callPrice.GetPriceReq{}
		req.Uid = fmt.Sprintf("%d", uid)
		orderDetail, err = grpcClient.GetDemo(ctx, &req)
		if err != nil {
			return orderDetail, err
		}
	} else {
		return orderDetail, err
	}
	return orderDetail, err
}
