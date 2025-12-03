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

type GetBalanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户余额
func NewGetBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBalanceLogic {
	return &GetBalanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBalanceLogic) GetBalance() (resp *types.BalanceResponse, err error) {
	// 查询用户余额
	accountID := l.ctx.Value(vars.AccountKey).(string)
	account, err := l.svcCtx.AccountModel.FindOne(l.ctx, accountID)
	if err != nil {
		return nil, err
	}
	resp = &types.BalanceResponse{
		Balance: account.Balance,
	}
	return
}
