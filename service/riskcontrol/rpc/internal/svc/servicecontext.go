package svc

import (
	"github.com/Drengin6306/ZeroBank/service/riskcontrol/model/redis"
	"github.com/Drengin6306/ZeroBank/service/riskcontrol/rpc/internal/config"
)

type ServiceContext struct {
	Config     config.Config
	RedisModel redis.RedisModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		RedisModel: redis.NewRedisModel(&c.Redis),
		Config:     c,
	}
}
