// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package transaction

import (
	"context"

	"github.com/Drengin6306/ZeroBank/pkg/errorx"
	"github.com/Drengin6306/ZeroBank/pkg/idgen"
	"github.com/Drengin6306/ZeroBank/pkg/vars"
	"github.com/Drengin6306/ZeroBank/service/account/rpc/account"
	"github.com/Drengin6306/ZeroBank/service/transaction/api/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/transaction/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransferLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransferLogic {
	return &TransferLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransferLogic) Transfer(req *types.TransferRequest) (resp *types.TransferResponse, err error) {
	if req.Amount <= 0 {
		return nil, errorx.NewError(errorx.ErrInvalidParams)
	}
	accountFrom := l.ctx.Value(vars.AccountKey).(string)
	balanceFrom, err := l.svcCtx.AccountRpc.QueryAccountBalance(l.ctx, &account.QueryAccountBalanceRequest{
		AccountId: accountFrom,
	})
	if err != nil {
		return nil, err
	}
	if balanceFrom.Balance < req.Amount {
		return nil, errorx.NewError(errorx.ErrBalanceNotEnough)
	}
	_, err = l.svcCtx.AccountRpc.DeductBalance(l.ctx, &account.DeductBalanceRequest{
		AccountId: accountFrom,
		Amount:    req.Amount,
	})
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.AccountRpc.AddBalance(l.ctx, &account.AddBalanceRequest{
		AccountId: req.AccountTo,
		Amount:    req.Amount,
	})
	if err != nil {
		return nil, err
	}
	transactionID := idgen.GenTransactionID()
	resp = &types.TransferResponse{
		TransactionID: transactionID,
		AccountFrom:   accountFrom,
		AccountTo:     req.AccountTo,
		Amount:        req.Amount,
	}
	return
}
