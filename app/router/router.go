package router

import (
	"time"

	"github.com/CRORCR/call/app/middleware"
	"github.com/CRORCR/call/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	// 默认带有Logger 和 Recovery 两个中间件
	//gin.SetMode(gin.ReleaseMode) // 输出调试信息
	gin.SetMode(gin.DebugMode) // 输出调试信息

	router := gin.Default()
	//中间件 Use设置中间件
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*", "lang", "json-token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowAllOrigins:  true,
	}))

	//加载自定义中间件
	//router.Use(contract.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.Cost())
	router.Use(middleware.Timeout(3 * time.Second))

	return router
}

func InitRouter() *gin.Engine {
	router := initRouter()

	userRouter := router.Group("/api/v1/call")
	{
		// 聊价查询
		userRouter.GET("/price", server.UserServer.CallPrice)
		userRouter.GET("/price/v2", server.UserServer.CallPriceUids)
		userRouter.DELETE("/someDelete", server.UserServer.CallPrice)
		userRouter.POST("/users/update", server.UserServer.CallPrice)
		//此规则能够匹配/user/lcq/30这种格式，但不能匹配/user/李长全/30 不支持中文，而且也不能为空，否则404
		userRouter.GET("/users/:name/:age", server.UserServer.CallPrice)
		//v1.Use(lib.JWTAuth())
	}

	return router
}
