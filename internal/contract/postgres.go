package contract

import (
	"fmt"
	"github.com/CRORCR/call/internal/config"
	"github.com/CRORCR/call/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/url"
	"sync/atomic"
)

type DataBase struct {
	Ins *DbGroup
}

// db连接加上主从 todo
type DbGroup struct {
	name    string     // 哪个逻辑库
	master  *gorm.DB   // 哪个主
	replica []*gorm.DB // 哪个从
	total   uint64     // 一共几个从库
	next    uint64     // 下一个从库
}

func InitPostgres(conf *config.Configuration) *DataBase {
	masterDb := initPostgres(conf.Conf.Postgres.DefaultMaster)
	slaveDb := initPostgres(conf.Conf.Postgres.DefaultSlave)

	// 现在不区分主从，以后再拆开
	ins := &DbGroup{name: "call_db", master: masterDb, replica: []*gorm.DB{slaveDb}}
	ins.total = uint64(len(ins.replica))
	return &DataBase{Ins: ins}
}

func initPostgres(conf model.DefaultDbConfig) *gorm.DB {
	dsn := url.URL{
		User:     url.UserPassword(conf.Username, conf.Password),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Path:     "duoo",
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}

	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()

	// 空闲连接池中的最大连接数
	sqlDB.SetMaxIdleConns(int(conf.MaxIdle))
	// 数据库的最大打开连接数
	sqlDB.SetMaxOpenConns(int(conf.MaxOpen))
	return db
}

// CloseDB 主从全部都要关闭
func CloseDB(dbList ...*DataBase) {
	for _, db := range dbList {
		for _, db2 := range db.Ins.replica {
			db3, err := db2.DB()
			if err != nil {
				logrus.Error("failed to close slave the database connection.")
			}
			db3.Close()
		}

		if db2, err := db.Master().DB(); err != nil {
			logrus.Error("failed to close master the database connection.")
		} else {
			db2.Close()
		}
	}
}

func (d DataBase) Master() *gorm.DB {
	return d.Ins.master
}

func (d DataBase) Slave() *gorm.DB {
	if d.Ins.replica == nil {
		return d.Ins.master
	}
	next := atomic.AddUint64(&d.Ins.next, 1)
	return d.Ins.replica[next%d.Ins.total]
}
