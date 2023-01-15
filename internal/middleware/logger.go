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

func Logger() gin.HandlerFunc {
	// 本地测试
	//path := "./logs/%Y%m%d%H%M.log"
	//writer, _ := retalog.New(
	//	path,
	//	//retalog.WithLinkName(path),
	//	retalog.WithMaxAge(time.Duration(180)*time.Second),
	//	retalog.WithRotationTime(time.Duration(60)*time.Second),
	//)

	path := "./logs/%Y%m%d%H.log"
	writer, _ := retalog.New(
		path,
		retalog.WithMaxAge(time.Hour*24*7),
		retalog.WithRotationTime(time.Hour),
	)

	logrus.SetOutput(writer)

	// todo
	// 如果env是本地，控制台输出日志
	// 如果env是prod，只能输出info级别日志

	//logrus.SetOutput(os.Stderr) // 控制台输出
	//logrus.SetOutput(ioutil.Discard) //控制台不输出

	logrus.SetLevel(logrus.TraceLevel)
	//logrus.SetReportCaller(true) // 行号是否输出

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
		// 日志格式
		logrus.WithFields(logrus.Fields{
			"status_code": statusCode,
			"cost":        cost,
			"client_ip":   clientIP,
			"req_method":  reqMethod,
			"req_uri":     reqUrl,
			"request":     string(bodyBytes),
		}).Trace()
	}
}

func LoggerV2() gin.HandlerFunc {
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
