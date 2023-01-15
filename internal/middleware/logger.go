package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	retalog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// 第一版本
func Logger() gin.HandlerFunc {
	filepath := "./logs/log.log"
	src, err := os.OpenFile(filepath, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0755)

	if err != nil {
		fmt.Println("err:", err)
	}

	logger := logrus.New()
	logger.Out = src
	logger.SetLevel(logrus.DebugLevel)

	logInfoWriter, _ := retalog.New(
		"./logs/info-%Y%m%d%H.log",
		retalog.WithMaxAge(7*24*time.Hour),
		retalog.WithRotationTime(time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logInfoWriter,
		logrus.DebugLevel: logInfoWriter,
		logrus.WarnLevel:  logInfoWriter,
		logrus.ErrorLevel: logInfoWriter,
		logrus.PanicLevel: logInfoWriter,
	}
	Hook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.AddHook(Hook)

	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		stopTime := time.Since(startTime).Microseconds()
		spendTime := fmt.Sprintf("%d ms", stopTime)
		clientIp := ctx.ClientIP()
		//userAgent := ctx.Request.UserAgent()

		method := ctx.Request.Method
		path := ctx.Request.RequestURI
		hostName, err := os.Hostname()
		statusCode := ctx.Writer.Status()
		if err != nil {
			hostName = "unknown"
		}
		entry := logger.WithFields(logrus.Fields{
			"host_name": hostName,
			"status":    statusCode,
			"cost":      spendTime,
			"ip":        clientIp,
			"method":    method,
			"path":      path,
			//"Agent":     userAgent,
		})

		if len(ctx.Errors) > 0 {
			entry.Error(ctx.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		if statusCode >= 500 {
			entry.Error()
		} else if statusCode >= 400 {
			entry.Info()
		} else {
			entry.Info()
		}
	}
}
