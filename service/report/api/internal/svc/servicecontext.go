// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/Drengin6306/ZeroBank/service/report/api/internal/config"
	"github.com/Drengin6306/ZeroBank/service/transaction/rpc/transaction"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	TransactionRpc transaction.Transaction
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		TransactionRpc: transaction.NewTransaction(zrpc.MustNewClient(c.RPC.Transaction)),
	}
}
