package svc

import (
	"github.com/Drengin6306/ZeroBank/service/transaction/model"
	"github.com/Drengin6306/ZeroBank/service/transaction/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                 config.Config
	TransactionRecordModel model.TransactionRecordModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:                 c,
		TransactionRecordModel: model.NewTransactionRecordModel(conn),
	}
}
