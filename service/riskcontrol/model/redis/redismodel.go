package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type RedisConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       int
}

type defaultRedisModel struct {
	rdb *redis.Client
}

type RedisModel interface {
	Close() error
	CheckDailyTransferLimit(accountID string, amount int64, dailyLimit int64) (bool, error)
	CheckDailyWithdrawLimit(accountID string, amount int64, dailyLimit int64) (bool, error)
}

func NewRedisModel(cfg *RedisConfig) RedisModel {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + strconv.Itoa(cfg.Port),
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
