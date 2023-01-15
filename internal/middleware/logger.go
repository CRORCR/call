package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	retalog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func LoggerV1() gin.HandlerFunc {
	fileName := "./logs/trace.log"
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(fmt.Sprintf("Loading failure：%v", err))
	}

	logger := logrus.New()
	//设置日志级别
	logger.SetLevel(logrus.TraceLevel)
	logger.Out = src
	//logrus.SetOutput(ioutil.Discard) //控制台不输出

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

		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // Write body back
		response, _ := c.Get("response")
		// 日志格式
		logger.WithFields(logrus.Fields{
			"status_code": statusCode,
			"cost":        cost,
			"client_ip":   clientIP,
			"req_method":  reqMethod,
			"req_uri":     reqUrl,
			"request":     string(bodyBytes),
			"response":    response,
		}).Trace()
	}
}

func Logger() gin.HandlerFunc {
	path := "./logs/%Y%m%d%H%M.log"
	/* 日志轮转相关函数
	`WithLinkName` 为最新的日志建立软连接
	`WithRotationTime` 设置日志分割的时间，隔多久分割一次
	WithMaxAge 和 WithRotationCount二者只能设置一个
	 `WithMaxAge` 设置文件清理前的最长保存时间
	 `WithRotationCount` 设置文件清理前最多保存的个数
	*/
	// 下面配置日志每隔 1 分钟轮转一个新文件，保留最近 3 分钟的日志文件，多余的自动清理掉。
	writer, _ := retalog.New(
		path,
		//retalog.WithLinkName(path),
		retalog.WithMaxAge(time.Duration(180)*time.Second),
		retalog.WithRotationTime(time.Duration(60)*time.Second),
	)
	logrus.SetOutput(writer)
	logrus.SetLevel(logrus.TraceLevel)

	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		stopTime := time.Since(startTime)
		cost := fmt.Sprintf("%d ms", int(math.Ceil(float64(stopTime.Nanoseconds()/1000000))))

		reqMethod := c.Request.Method
		reqUrl := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // Write body back
		response, _ := c.Get("response")
		// 日志格式
		logrus.WithFields(logrus.Fields{
			"status_code": statusCode,
			"cost":        cost,
			"client_ip":   clientIP,
			"req_method":  reqMethod,
			"req_uri":     reqUrl,
			"request":     string(bodyBytes),
			"response":    response,
		}).Trace()
	}
}
