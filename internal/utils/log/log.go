package log

import (
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	Logger   = logrus.New()
	LogEntry *logrus.Entry
)

// 这种方式分文件不报错，一会看看

func GetLogEntry(logpath, loglink string) *logrus.Entry {
	filename := path.Join(logpath, loglink)
	if LogEntry != nil {
		return LogEntry
	}
	src, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil
	}
	Logger.Out = src
	logWriter, _ := rotatelogs.New(
		filename+".%Y%m%d%M.log",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	Logger.AddHook(lfHook)
	LogEntry = logrus.NewEntry(Logger).WithField("service", "")
	return LogEntry
}
