package config

import (
	"github.com/Drengin6306/ZeroBank/service/riskcontrol/model/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type Config struct {
	zrpc.RpcServerConf
	Consul consul.Conf
	Redis  redis.RedisConfig
}
