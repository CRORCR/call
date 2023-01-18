package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CRORCR/call/app/http/api"
	"github.com/CRORCR/call/app/http/middleware"
	"github.com/CRORCR/call/app/http/router"
	"github.com/CRORCR/call/internal/config"
	"github.com/CRORCR/call/internal/grpc"
	"github.com/CRORCR/call/internal/service"
)

func main() {
	// 加载配置
	config := config.InitConfig()
	middleware.NewLogger(config.Conf.Log)

	// 初始化rpc
	rpcService := grpc.InitRpcClient(config)
	//初始化 repo

	// 初始化service
	userService := service.NewUserService(config, rpcService)
	api.NewUserController(userService)
	appHandler := router.InitRouter()
	server := &http.Server{
		Handler: appHandler,
		Addr:    config.Conf.App.Port,
	}

	fmt.Printf("\nstart http server [%s] on [%s] \n", config.Conf.App.ServiceName, server.Addr)

	// 这个goroutine是启动服务的goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 当前的goroutine等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit

	// 调用Server.Shutdown graceful结束
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Printf("Server %s exiting \n", config.Conf.App.ServiceName)
}
