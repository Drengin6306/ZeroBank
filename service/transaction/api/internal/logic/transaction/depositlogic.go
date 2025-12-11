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
	"github.com/Drengin6306/ZeroBank/service/transaction/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type DepositLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDepositLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepositLogic {
	return &DepositLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DepositLogic) Deposit(req *types.DepositRequest) (resp *types.DepositResponse, err error) {
	if req.Amount <= 0 {
		return nil, errorx.NewError(errorx.ErrInvalidParams)
	}
	// 获取accountID
	accountID := l.ctx.Value(vars.AccountKey).(string)

	result, err := l.svcCtx.AccountRpc.AddBalance(l.ctx, &account.AddBalanceRequest{
		AccountId: accountID,
		Amount:    req.Amount,
	})
	if err != nil {
		return nil, err
	}
	transactionID := idgen.GenTransactionID()
	// 记录交易流水
	_, err = l.svcCtx.TransactionRecordModel.Insert(l.ctx, &model.TransactionRecord{
		TransactionId:   transactionID,
		AccountFrom:     accountID,
		Amount:          req.Amount,
		TransactionType: vars.TransactionTypeDeposit,
		Status:          vars.TransactionStatusSuccess,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.DepositResponse{
		TransactionID: transactionID,
		AccountID:     accountID,
		Balance:       result.Balance,
	}
	return
}
