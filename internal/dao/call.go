package dao

import (
	"context"
	"fmt"
	"github.com/CRORCR/call/internal/contract"
	"github.com/sirupsen/logrus"
	"time"
)

type CallRepository interface {
	GetDialCallById(ctx context.Context, uid int64) ([]DialCall, error)
	Lock(ctx context.Context, key string) (string, bool) // 加锁
	UnLock(ctx context.Context, key string, rid string)  // 解锁
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

func (c callDao) GetDialCallById(ctx context.Context, uid int64) ([]DialCall, error) {
	var dial []DialCall
	if err := c.db.Slave().Table("call_db.dial_call").Where("id in (?)", []int64{1, 2}).Find(&dial).Error; err != nil {
		logrus.Errorf("dao.GetDialCallById.err:%v", err)
		return dial, err
	}

	fmt.Println("输出结果：", dial)
	return dial, nil
}

func (c callDao) Lock(ctx context.Context, key string) (string, bool) {
	logger := logrus.WithFields(logrus.Fields{"request_id": ctx.Value("requestId"), "key": key})

	logger.Infof("dao.GetDialCallById.Lock,key:%v", key)
	return c.redis.Lock(key, 3*time.Second)
}

func (c callDao) UnLock(ctx context.Context, key string, rid string) {
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

/*
测试批量发送，一次执行
	pipeline := c.redis.Client.Pipeline()
	pipeline.Set("hello", 1, time.Minute)
	pipeline.Set("hello2", 2, time.Minute)
	pipeline.Set("hello3", 3, time.Second*30)
	cmders, err := pipeline.Exec()
	if err != nil {
		fmt.Println("错误了")
		return nil, nil
	}
	data := make([]string, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.StatusCmd).Val()
		data = append(data, v)
	}

	fmt.Println("-------jieguo", data) //  [OK OK OK]
	return nil, nil
*/