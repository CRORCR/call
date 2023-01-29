package dao

import (
	"fmt"
	"github.com/CRORCR/call/internal/contract"
	"github.com/sirupsen/logrus"
	"time"
)

type CallRepository interface {
	GetDialCallById(uid int64) ([]DialCall, error)
	Lock(string) (string, bool)    // 加锁
	UnLock(key string, rid string) // 解锁
}

type callDao struct {
	db    *contract.DataBase
	redis *contract.Redis
}

func CreateUserRepo(db *contract.DataBase, redis *contract.Redis) CallRepository {
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

func (c callDao) GetDialCallById(uid int64) ([]DialCall, error) {
	var dial []DialCall
	if err := c.db.Slave().Table("call_db.dial_call").Where("id in (?)", []int64{1, 2}).Find(&dial).Error; err != nil {
		logrus.Errorf("dao.GetDialCallById.err:%v", err)
		return dial, err
	}

	fmt.Println("输出结果：", dial)
	return dial, nil
}

func (c callDao) Lock(key string) (string, bool) {
	return c.redis.Lock(key, 3*time.Second)
}

func (c callDao) UnLock(key string, rid string) {
	if err := c.redis.ReleaseLock(key, rid); err != nil {
		logrus.Errorf("callDao.UnLock.err:%v", err)
	}
}

type DialCall struct {
	Id      int64 `json:"id" gorm:"primary_key:auto_increment" `
	FromUid int64 `json:"from_uid" gorm:"from_uid"`
	ToUid   int64 `json:"to_uid" gorm:"to_uid"`
	//`create_time` bigint(20) unsigned NOT NULL COMMENT '创建时间',
	//`update_time` bigint(20) unsigned NOT NULL COMMENT '最后更新时间',
}
