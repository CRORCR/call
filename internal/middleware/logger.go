package middleware

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	retalog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	fileName := "./log/trace"
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(fmt.Sprintf("Loading failure：%v", err))
	}

	logger := logrus.New()
	//设置日志级别
	logger.SetLevel(logrus.TraceLevel)
	logger.Out = src

	// 设置 rotatelogs
	logWriter, err := retalog.New(
		// 分割后的文件名称
		fileName+"-%Y%m%d%H.log",

		// 生成软链，指向最新日志文件
		retalog.WithLinkName(fileName),

		// 设置最大保存时间(7天)
		retalog.WithMaxAge(7*24*time.Hour),

		// 设置日志切割时间间隔(1h)
		retalog.WithRotationTime(time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.TraceLevel: logWriter,
	}

	logger.AddHook(lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))

	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		stopTime := time.Since(startTime)
		cost := fmt.Sprintf("%d ms", int(math.Ceil(float64(stopTime.Nanoseconds()/1000000))))

		reqMethod := c.Request.Method
		reqUrl := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		// 日志格式
		logger.WithFields(logrus.Fields{
			"status_code": statusCode,
			"cost":        cost,
			"client_ip":   clientIP,
			"req_method":  reqMethod,
			"req_uri":     reqUrl,
		}).Trace()
	}
}
