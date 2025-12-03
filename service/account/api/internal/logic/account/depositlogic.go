// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package account

import (
	"context"

	"github.com/Drengin6306/ZeroBank/pkg/vars"
	"github.com/Drengin6306/ZeroBank/service/account/api/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/account/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DepositLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 存款
func NewDepositLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepositLogic {
	return &DepositLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DepositLogic) Deposit(req *types.DepositRequest) (resp *types.DepositResponse, err error) {
	// 存款
	accountID := l.ctx.Value(vars.AccountKey).(string)
	account, err := l.svcCtx.AccountModel.FindOne(l.ctx, accountID)
	if err != nil {
		return nil, err
	}
	account.Balance += req.Amount
	err = l.svcCtx.AccountModel.Update(l.ctx, account)
	if err != nil {
		return nil, err
	}
	resp = &types.DepositResponse{
		Balance: account.Balance,
	}
	return
}
