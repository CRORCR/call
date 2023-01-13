package contract

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var t = time.NewTicker(24 * time.Hour)
var file *os.File
var err error

func Logger() {
	go func() {
		for {
			select {
			case <-t.C:
				for path, v := range loggerMap {
					func() {
						//文件名: 日期+模块名 组合
						file, err = os.OpenFile(fmt.Sprintf("../loggers/%v_%v.log", timeToStringEx(), path), os.O_CREATE|os.O_APPEND, 0666)
						if err != nil {
							fmt.Printf("open %s failed, err:%v\n", path, err)
							return
						}
						v.SetOutput(file)
						v.WithFields(logrus.Fields{
							"test": "demo",
						}).Info("success")
					}()
				}
			}
		}
	}()
}

func timeToStringEx() string {
	return time.Now().Format("2006-01-02")
}

var loggerMap = make(map[string]*logrus.Logger, 0)

func newLogrus(path string) (log *logrus.Logger) {
	log = logrus.New()
	loggerMap[path] = log
	file, err = os.OpenFile(fmt.Sprintf("./loggers/%v_%v.txt", timeToStringEx(), path), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("open filed", err)
		return nil
	}
	log.SetOutput(file)
	return
}

func GetLogrus(path string) (log *logrus.Logger) {
	if loggerMap[path] != nil {
		return loggerMap[path]
	}
	log = newLogrus(path)
	if log != nil {
		loggerMap[path] = log
	}
	return log
}
