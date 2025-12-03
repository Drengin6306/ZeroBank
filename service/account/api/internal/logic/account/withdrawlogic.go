// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package account

import (
	"context"

	"github.com/Drengin6306/ZeroBank/pkg/errorx"
	"github.com/Drengin6306/ZeroBank/pkg/vars"
	"github.com/Drengin6306/ZeroBank/service/account/api/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/account/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WithdrawLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 取款
func NewWithdrawLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawLogic {
	return &WithdrawLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WithdrawLogic) Withdraw(req *types.WithdrawRequest) (resp *types.WithdrawResponse, err error) {
	// 取款
	accountID := l.ctx.Value(vars.AccountKey).(string)
	account, err := l.svcCtx.AccountModel.FindOne(l.ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account.Balance < req.Amount {
		return nil, errorx.NewError(errorx.ErrBalanceNotEnough)
	}
	account.Balance -= req.Amount
	err = l.svcCtx.AccountModel.Update(l.ctx, account)
	if err != nil {
		return nil, err
	}
	resp = &types.WithdrawResponse{
		Balance: account.Balance,
	}
	return
}
