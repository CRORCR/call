package contract

import (
	"fmt"
	"github.com/CRORCR/call/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/url"
)

// db连接加上主从 todo

func InitPostgres(conf *config.Configuration) *gorm.DB {
	dsn := url.URL{
		User:     url.UserPassword(conf.Conf.Postgres.DefaultMaster.Username, conf.Conf.Postgres.DefaultMaster.Password),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%d", conf.Conf.Postgres.DefaultMaster.Host, conf.Conf.Postgres.DefaultMaster.Port),
		Path:     "duoo",
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}

	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()

	// 空闲连接池中的最大连接数
	sqlDB.SetMaxIdleConns(10)
	// 数据库的最大打开连接数
	sqlDB.SetMaxOpenConns(100)

	return db
}

func CloseDB(db *gorm.DB) {
	db2, err := db.DB()

	if err != nil {
		panic("failed to close the database connection.")
	}

	db2.Close()
}
