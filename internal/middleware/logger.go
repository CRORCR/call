package middleware

import (
	"fmt"
	"math"
	"time"

	"github.com/CRORCR/call/internal/utils/log"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	//entry := log.GetLogEntry(config.C.Logger.FilePath, config.C.Logger.FileName)
	entry := log.GetLogEntry("logs", "latest")
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		stopTime := time.Since(startTime)
		spendTime := fmt.Sprintf("%d ms", int(math.Ceil(float64(stopTime.Nanoseconds()/1000000))))
		statusCode := c.Writer.Status()
		dataSize := c.Writer.Size()
		if dataSize < 0 {
			dataSize = 0
		}
		clientIp := c.ClientIP()
		userAgent := c.Request.UserAgent()
		method := c.Request.Method
		url := c.Request.RequestURI
		l := entry.WithFields(logrus.Fields{
			"SpendTime": spendTime,
			"ClientIP":  clientIp,
			"UserAgent": userAgent,
			"Path":      url,
			"Method":    method,
			"Status":    statusCode,
			"DataSize":  dataSize,
		})
		if len(c.Errors) > 0 {
			l.Error(c.Errors.ByType(gin.ErrorTypePrivate))
		}
		if statusCode >= 500 {
			l.Error()
		} else if statusCode >= 400 {
			l.Warn()
		} else {
			l.Info()
		}
	}
}
