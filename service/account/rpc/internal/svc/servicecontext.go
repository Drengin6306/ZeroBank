package svc

import (
	"github.com/Drengin6306/ZeroBank/service/account/model"
	"github.com/Drengin6306/ZeroBank/service/account/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                  config.Config
	AccountModel            model.AccountModel
	CustomerIndividualModel model.CustomerIndividualModel
	CustomerEnterpriseModel model.CustomerEnterpriseModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:                  c,
		AccountModel:            model.NewAccountModel(conn),
		CustomerIndividualModel: model.NewCustomerIndividualModel(conn),
		CustomerEnterpriseModel: model.NewCustomerEnterpriseModel(conn),
	}
}
