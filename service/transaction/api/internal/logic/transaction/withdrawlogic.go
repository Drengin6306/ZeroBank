// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package transaction

import (
	"context"

	"github.com/Drengin6306/ZeroBank/pkg/errorx"
	"github.com/Drengin6306/ZeroBank/pkg/idgen"
	"github.com/Drengin6306/ZeroBank/pkg/vars"
	"github.com/Drengin6306/ZeroBank/service/account/rpc/account"
	"github.com/Drengin6306/ZeroBank/service/riskcontrol/rpc/riskcontrol"
	"github.com/Drengin6306/ZeroBank/service/transaction/api/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/transaction/api/internal/types"
	"github.com/Drengin6306/ZeroBank/service/transaction/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type WithdrawLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWithdrawLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawLogic {
	return &WithdrawLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WithdrawLogic) Withdraw(req *types.WithdrawRequest) (resp *types.WithdrawResponse, err error) {
	if req.Amount <= 0 {
		return nil, errorx.NewError(errorx.ErrInvalidParams)
	}
	accountID := l.ctx.Value(vars.AccountKey).(string)
	info, err := l.svcCtx.AccountRpc.GetAccountInfo(l.ctx, &account.AccountInfoRequest{
		AccountId: accountID,
	})
	if err != nil {
		return nil, err
	}
	if info.GetBalance() < req.Amount {
		return nil, errorx.NewError(errorx.ErrBalanceNotEnough)
	}
	transactionID := idgen.GenTransactionID()
	// 风控检查
	riskResp, err := l.svcCtx.RiskControlRpc.CheckTransaction(l.ctx, &riskcontrol.RiskCheckRequest{
		AccountFrom:     accountID,
		AccountType:     int32(info.GetAccountType()),
		TransactionType: vars.TransactionTypeWithdraw,
		TransactionId:   transactionID,
		Amount:          req.Amount,
	})
	if err != nil {
		return nil, err
	}
	if !riskResp.Passed {
		// 交易单号加拒绝原因
		msg := "TransactionID: " + transactionID + ", Reason: " + riskResp.Reason
		return nil, errorx.NewErrorWithMsg(errorx.ErrRiskControl, msg)
	}
	result, err := l.svcCtx.AccountRpc.DeductBalance(l.ctx, &account.DeductBalanceRequest{
		AccountId: accountID,
		Amount:    req.Amount,
	})
	if err != nil {
		return nil, err
	}
	// 记录交易流水
	_, err = l.svcCtx.TransactionRecordModel.Insert(l.ctx, &model.TransactionRecord{
		TransactionId:   transactionID,
		AccountFrom:     accountID,
		Amount:          req.Amount,
		TransactionType: vars.TransactionTypeWithdraw,
		Status:          vars.TransactionStatusSuccess,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.WithdrawResponse{
		TransactionID: transactionID,
		AccountID:     accountID,
		Balance:       result.Balance,
	}
	return
}
