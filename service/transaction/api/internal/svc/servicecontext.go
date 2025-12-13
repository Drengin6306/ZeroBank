// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/Drengin6306/ZeroBank/service/account/rpc/account"
	"github.com/Drengin6306/ZeroBank/service/riskcontrol/rpc/riskcontrol"
	"github.com/Drengin6306/ZeroBank/service/transaction/api/internal/config"
	"github.com/Drengin6306/ZeroBank/service/transaction/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                 config.Config
	AccountRpc             account.Account
	RiskControlRpc         riskcontrol.RiskControl
	TransactionRecordModel model.TransactionRecordModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:                 c,
		AccountRpc:             account.NewAccount(zrpc.MustNewClient(c.RPC.Account)),
		RiskControlRpc:         riskcontrol.NewRiskControl(zrpc.MustNewClient(c.RPC.RiskControl)),
		TransactionRecordModel: model.NewTransactionRecordModel(conn),
	}
}
