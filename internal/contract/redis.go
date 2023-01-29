package contract

import (
	"fmt"
	"github.com/CRORCR/call/internal/config"
	"github.com/CRORCR/call/internal/model"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/go-redis/redis"
)

// redis 搭建

type Redis struct {
	Client *redis.Client
}

var (
	RainKey = "call"
)

func InitRedisClient(conf *config.Configuration) *Redis {
	return &Redis{
		Client: CreateRedisConnection(conf.Conf.Redis),
	}
}

func CreateRedisConnection(conf model.RedisConfig) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       int(conf.DB),

		//连接池容量及闲置连接数量
		PoolSize:     10,                  // 连接池最大socket连接数
		MinIdleConns: int(conf.MaxActive), // 在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量

		//超时
		DialTimeout:  time.Duration(conf.ConnectTimeout) * time.Millisecond, //连接超时时间
		ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,    //读超时，-1表示取消读超时
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,   //写超时
		PoolTimeout:  4 * time.Second,                                       //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长

		//闲置连接检查
		IdleCheckFrequency: 60 * time.Second,                              //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查
		IdleTimeout:        time.Duration(conf.IdleTimeout) * time.Minute, //闲置超时，默认3分钟，-1表示取消闲置超时检查
		MaxConnAge:         0 * time.Second,                               //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

		//命令执行失败时的重试策略
		MaxRetries:      int(conf.Retry),        // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
	})
	if _, err := redisClient.Ping().Result(); err != nil {
		panic(err)
	}

	if redisClient == nil {
		panic(fmt.Sprintf("InitRedis Error  host:%s, port:%d", conf.Host, conf.Port))
	}

	return redisClient
}

func (r *Redis) RedisClose() {
	_ = r.Client.Close()
}

// ------- 通用操作 -------

func (r *Redis) Expire(key string, t time.Duration) error {
	return r.Client.Expire(key, t).Err()
}

// ------- 字符串操作 -------

func (r *Redis) Set(key string, value string, expiration time.Duration) error {
	if err := r.Client.Set(key, value, expiration).Err(); err != nil {
		return fmt.Errorf("redis set error:%v", err)
	}
	return nil
}

func (r *Redis) Get(key string) string {
	result, err := r.Client.Get(key).Result()
	if err != nil {
		//if err.Error() == `redis: nil` {
		//	return ``
		//}
		if err == redis.Nil {
			fmt.Println("todo 是空的么？")
		}
	}
	return result
}

// ------- hash操作 -------

func (r *Redis) SetHash(key string, fields map[string]interface{}, expiration time.Duration) error {
	if err := r.Client.HMSet(key, fields).Err(); err != nil {
		return err
	}
	if expiration > 0 {
		r.Client.Expire(key, expiration)
	}
	return nil
}

func (r *Redis) GetHash(key string) (map[string]string, error) {
	return r.Client.HGetAll(key).Result()
}

func (r *Redis) HGet(key, field string) string {
	result, err := r.Client.HGet(key, field).Result()
	if err != nil {
		if err.Error() == `redis: nil` {
			return ""
		}
		logrus.Error(`redis hget error`, err)
	}
	return result
}

func (r *Redis) GetHLength(id string) int64 {
	return r.Client.HLen(fmt.Sprintf("%v_%v", RainKey, id)).Val()
}
