package svc

import (
	"github.com/Drengin6306/ZeroBank/service/riskcontrol/model/mysql"
	"github.com/Drengin6306/ZeroBank/service/riskcontrol/model/redis"
	"github.com/Drengin6306/ZeroBank/service/riskcontrol/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config           config.Config
	RedisModel       redis.RedisModel
	RiskRecordModel  mysql.RiskRecordModel
	LimitAmountModel mysql.LimitAmountModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:           c,
		RedisModel:       redis.NewRedisModel(c.RiskRedis),
		RiskRecordModel:  mysql.NewRiskRecordModel(conn),
		LimitAmountModel: mysql.NewLimitAmountModel(conn, c.Cache),
	}
}
