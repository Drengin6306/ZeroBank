package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type conf struct {
	Host     string
	Username string `json:",optional"`
	Password string `json:",optional"`
	DB       int    `json:",optional"`
}

type RedConf *conf

type defaultRedisModel struct {
	rdb *redis.Client
}

type RedisModel interface {
	Close() error
	CheckDailyTransferLimit(accountID string, amount float64, dailyLimit float64) (bool, error)
	CheckDailyWithdrawLimit(accountID string, amount float64, dailyLimit float64) (bool, error)
}

func NewRedisModel(cfg RedConf) RedisModel {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	// ping
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logx.Error("redis ping error:", err)
	}
	return &defaultRedisModel{rdb: rdb}
}

func (m *defaultRedisModel) Close() error {
	return m.rdb.Close()
}
