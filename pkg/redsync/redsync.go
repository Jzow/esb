// Package redsync 实现了基于redis的分布式锁
package redsync

import (
	"errors"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/redigo"
	"github.com/gomodule/redigo/redis"
	"time"
)

var instance *redMutex

type redMutex struct {
	pool redsyncredis.Pool
	sync *redsync.Redsync
}

// RedLock struct名修改为大写，以前的小写导致无法在外部引用，非常蛋疼
type RedLock struct {
	mutex *redsync.Mutex
}

func NewPool(host, password string) (*redis.Pool, error) {
	pool := &redis.Pool{ //redis连接池
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", host,
				redis.DialPassword(password))
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}

	conn := pool.Get()
	defer conn.Close()
	if _, err := conn.Do("ping"); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

// Setup 初始化
func Setup() {
	if instance != nil {
		return
	}
	pool, err := NewPool(setting.RedisSetting.Host, setting.RedisSetting.Password)
	if err != nil {
		util.Log("failed yo NewPool for redis lock", err)
		return
	}
	instance = &redMutex{}
	instance.pool = redigo.NewPool(pool)
	instance.sync = redsync.New(instance.pool)
}

// GetLock 获取锁key
func GetLock(key string, expiry time.Duration) (lock *RedLock, err error) {
	if instance == nil {
		err = errors.New("not init")
		return
	}
	mu := instance.sync.NewMutex(key, redsync.WithExpiry(expiry))
	err = mu.Lock()
	if err == nil {
		lock = &RedLock{mutex: mu}
	}
	return
}

// Unlock 释放获得的锁
func (m *RedLock) Unlock() (err error) {
	if m == nil || m.mutex == nil {
		err = errors.New("invalid lock")
		return
	}
	_, err = m.mutex.Unlock()
	return
}
