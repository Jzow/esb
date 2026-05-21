package util

import (
	"context"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"strings"
	"time"
)

var Rdb *redis.Client

func SetupRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:        setting.RedisSetting.Host,
		Password:    setting.RedisSetting.Password,
		DialTimeout: setting.RedisSetting.IdleTimeout,
		DB:          0,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			logging.Info("[redis]build a redis client")
			return nil
		},
	})
}

func GetRedisCache(key, logTag string) (string, error) {
	value, err := Rdb.Get(context.Background(), key).Result()
	if err == redis.Nil {
		logging.Infof("[%s] cache not exist,when invoke GetRedisCache,key=%s ", logTag, key)
		return "", err
	} else if err != nil {
		logging.Errorf("[%s] failed to GetRedisCache, key=%s ,err: %s", logTag, key, err.Error())
		return "", err
	}
	logging.Infof("[%s] GetRedisCache ,key=%s", logTag, key)
	return value, nil
}

func SetRedisCache(key, value, logTag string, duration time.Duration) error {
	err := Rdb.Set(context.Background(), key, value, duration).Err()
	if err != nil {
		logging.Errorf("[%s] failed to SetRedisCache, key=%s ,value=%s ,err: %s", logTag, value, err.Error())
	}
	logging.Infof("[%s] SetRedisCache ,key=%s", logTag, key)
	return err
}

func GetRedisMutex(mutexName string, expiry time.Duration) *redsync.Mutex {
	pool := goredis.NewPool(Rdb)
	rs := redsync.New(pool)
	mutex := rs.NewMutex(
		strings.ReplaceAll(mutexName, "_", "-"),
		redsync.WithExpiry(expiry),
		redsync.WithTries(2),
	)
	return mutex
}
