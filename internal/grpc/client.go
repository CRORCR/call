package grpc

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/CRORCR/call/internal/config"
	his_order "github.com/CRORCR/cr-common/proto/engine_order"
	"github.com/CRORCR/cr-common/proto/his_otc"
	"github.com/CRORCR/cr-common/proto/his_transfer"
	"github.com/CRORCR/cr-common/proto/his_wallet"
	callPrice "github.com/CRORCR/duoo-common/proto/call_price"
	"google.golang.org/grpc"
)

/*
import (
	"sync"
	"sync/atomic"
	"unsafe"

	"git.zgtest.club/bhex/df-common/infrastructure/proto/his_api/his_otc"
	"git.zgtest.club/bhex/df-common/infrastructure/proto/his_api/his_transfer"
	"git.zgtest.club/bhex/df-common/infrastructure/proto/his_api/his_wallet"
	"git.zgtest.club/bhex/df-common/infrastructure/proto/org_api/payment"
	"git.zgtest.club/bhex/df-common/infrastructure/proto/org_api/user"
	his_order "github.com/CRORCR/cr-common/proto"
	"google.golang.org/grpc"
)

var (
	globalClientConnOrj      unsafe.Pointer
	globalClientConn         unsafe.Pointer
	lck                      sync.Mutex
	targetConnectionAddress  string
	historyConnectionAddress string
)

func InitGRPCClient() {
	targetConnectionAddress = config.Public.GetString("org_api")
	historyConnectionAddress = config.Public.GetString("his_api")
}

func GetPaymentClient() (payment.PaymentClient, error) {
	conn, err := GetConn(targetConnectionAddress)
	if err != nil {
		return (payment.PaymentClient)(nil), err
	}
	return payment.NewPaymentClient(conn), nil //此处调用pb.go文件中生成的创建client的方法
}

func GetHisTransferClient() (his_transfer.HisTransferClient, error) {
	conn, err := GetConn(historyConnectionAddress)
	if err != nil {
		return (his_transfer.HisTransferClient)(nil), err
	}
	return his_transfer.NewHisTransferClient(conn), nil //此处调用pb.go文件中生成的创建client的方法
}

func GetHisWalletClient() (his_wallet.HisWalletClient, error) {
	conn, err := GetConn(historyConnectionAddress)
	if err != nil {
		return (his_wallet.HisWalletClient)(nil), err
	}
	return his_wallet.NewHisWalletClient(conn), nil //此处调用pb.go文件中生成的创建client的方法
}

func GetHisNormalClient() (his_otc.HisOtcListClient, error) {
	conn, err := GetConn(historyConnectionAddress)
	if err != nil {
		return (his_otc.HisOtcListClient)(nil), err
	}
	return his_otc.NewHisOtcListClient(conn), nil //此处调用pb.go文件中生成的创建client的方法
}

func GetHisOrderClient() (his_order.HisOrderListClient, error) {
	conn, err := GetConn(historyConnectionAddress)
	if err != nil {
		return (his_order.HisOrderListClient)(nil), err
	}
	return his_order.NewHisOrderListClient(conn), nil //此处调用pb.go文件中生成的创建client的方法
}

func GetUserClient() (user.UserClient, error) {
	conn, err := GetConnOrg(targetConnectionAddress)
	if err != nil {
		return (user.UserClient)(nil), err
	}
	return user.NewUserClient(conn), nil //此处调用pb.go文件中生成的创建client的方法
}

func GetConn(target string) (*grpc.ClientConn, error) {
	if atomic.LoadPointer(&globalClientConn) != nil {
		return (*grpc.ClientConn)(globalClientConn), nil
	}
	lck.Lock()
	defer lck.Unlock()
	if atomic.LoadPointer(&globalClientConn) != nil { //double check
		return (*grpc.ClientConn)(globalClientConn), nil
	}
	cli, err := newGrpcConn(target)
	if err != nil {
		return nil, err
	}
	atomic.StorePointer(&globalClientConn, unsafe.Pointer(cli))
	return cli, nil
}

func GetConnOrg(target string) (*grpc.ClientConn, error) {
	if atomic.LoadPointer(&globalClientConnOrj) != nil {
		return (*grpc.ClientConn)(globalClientConnOrj), nil
	}
	lck.Lock()
	defer lck.Unlock()
	if atomic.LoadPointer(&globalClientConnOrj) != nil { //double check
		return (*grpc.ClientConn)(globalClientConnOrj), nil
	}
	cli, err := newGrpcConn(target)
	if err != nil {
		return nil, err
	}
	atomic.StorePointer(&globalClientConnOrj, unsafe.Pointer(cli))
	return cli, nil
}

func newGrpcConn(target string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		target,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
*/

type RpcService struct {
	userClient     unsafe.Pointer
	lck            sync.Mutex
	userServerPort string
	hisConnAddress string
	conf           *config.Configuration
}

func InitRpcClient(conf *config.Configuration) *RpcService {
	r := &RpcService{
		userServerPort: conf.Conf.Rpc.UserServer,
	}
	return r
}

func (r *RpcService) GetUserClient() (callPrice.HisDemoListClient, error) {
	conn, err := r.getUserConn(r.userServerPort)
	if err != nil {
		return (callPrice.HisDemoListClient)(nil), err
	}
	return callPrice.NewHisDemoListClient(conn), nil //此处调用pb.go文件中生成的创建client的方法
}

func (r *RpcService) GetHisTransferClient() (his_transfer.HisTransferClient, error) {
	conn, err := r.getUserConn(r.hisConnAddress)
	if err != nil {
		return (his_transfer.HisTransferClient)(nil), err
	}
	return his_transfer.NewHisTransferClient(conn), nil //此处调用pb.go文件中生成的创建client的方法
}

func (r *RpcService) GetHisWalletClient() (his_wallet.HisWalletClient, error) {
	conn, err := r.getUserConn(r.hisConnAddress)
	if err != nil {
		return (his_wallet.HisWalletClient)(nil), err
	}
	return his_wallet.NewHisWalletClient(conn), nil //此处调用pb.go文件中生成的创建client的方法
}

func (r *RpcService) GetHisNormalClient() (his_otc.HisOtcListClient, error) {
	conn, err := r.getUserConn(r.hisConnAddress)
	if err != nil {
		return (his_otc.HisOtcListClient)(nil), err
	}
	return his_otc.NewHisOtcListClient(conn), nil //此处调用pb.go文件中生成的创建client的方法
}

func (r *RpcService) GetHisOrderClient() (his_order.HisOrderListClient, error) {
	conn, err := r.getUserConn(r.hisConnAddress)
	if err != nil {
		return (his_order.HisOrderListClient)(nil), err
	}
	return his_order.NewHisOrderListClient(conn), nil //此处调用pb.go文件中生成的创建client的方法
}

func (r *RpcService) getUserConn(target string) (*grpc.ClientConn, error) {
	if atomic.LoadPointer(&r.userClient) != nil {
		return (*grpc.ClientConn)(r.userClient), nil
	}
	r.lck.Lock()
	defer r.lck.Unlock()
	if atomic.LoadPointer(&r.userClient) != nil { //double check
		return (*grpc.ClientConn)(r.userClient), nil
	}
	cli, err := r.newGrpcConn(target)
	if err != nil {
		return nil, err
	}
	atomic.StorePointer(&r.userClient, unsafe.Pointer(cli))
	return cli, nil
}

func (r *RpcService) newGrpcConn(target string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		target,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
