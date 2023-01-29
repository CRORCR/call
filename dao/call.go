package dao

import (
	"fmt"
	"github.com/CRORCR/call/internal/contract"
	"gorm.io/gorm"
)

type CallRepository interface {
	GetDialCallById(uid int64) DialCall
}

type callDao struct {
	db    *gorm.DB
	redis *contract.Redis
}

func CreateUserRepo(db *gorm.DB, redis *contract.Redis) CallRepository {
	return &callDao{
		db:    db,
		redis: redis,
	}
}

//func (c callDao) GetDialCallById(uid int64) DialCall {
//	var dial []DialCall
//	err := c.db.Table("call_db.dial_call").Raw("select * from call_db.dial_call limit 1").Scan(&dial).Error
//	fmt.Println("有错误么？", err)
//	fmt.Println("有错误么-2？", dial)
//	return DialCall{}
//}

func (c callDao) GetDialCallById(uid int64) DialCall {
	var dial []DialCall
	err := c.db.Table("call_db.dial_call").Where("id in (?)", []int64{1, 2}).Find(&dial).Error

	fmt.Println("有错误么？", err)
	fmt.Println("有错误么-2？", dial)
	return DialCall{}
}

type DialCall struct {
	Id      int64 `json:"id" gorm:"primary_key:auto_increment" `
	FromUid int64 `json:"from_uid" gorm:"from_uid"`
	ToUid   int64 `json:"to_uid" gorm:"to_uid"`
	//`create_time` bigint(20) unsigned NOT NULL COMMENT '创建时间',
	//`update_time` bigint(20) unsigned NOT NULL COMMENT '最后更新时间',
}
