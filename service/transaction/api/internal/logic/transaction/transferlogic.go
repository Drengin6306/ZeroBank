// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package transaction

import (
	"context"
	"database/sql"

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

	exists, err := l.svcCtx.AccountRpc.IsAccountExist(l.ctx, &account.AccountInfoRequest{
		AccountId: req.AccountTo,
	})
	if err != nil {
		return nil, err
	}
	if !exists.Exist {
		return nil, errorx.NewError(errorx.ErrAccountNotFound)
	}

	accountFrom := l.ctx.Value(vars.AccountKey).(string)
	info, err := l.svcCtx.AccountRpc.GetAccountInfo(l.ctx, &account.AccountInfoRequest{
		AccountId: accountFrom,
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
		AccountFrom:     accountFrom,
		AccountTo:       req.AccountTo,
		AccountType:     int32(info.GetAccountType()),
		TransactionType: vars.TransactionTypeTransfer,
		TransactionId:   transactionID,
		Amount:          req.Amount,
	})
	if err != nil {
		return nil, err
	}
	if !riskResp.Passed {
		// 交易单号加拒绝原因
		msg := riskResp.Reason + " (交易流水号: " + transactionID + ")"
		return nil, errorx.NewErrorWithMsg(errorx.ErrRiskControl, msg)
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
	// 记录交易流水
	_, err = l.svcCtx.TransactionRecordModel.Insert(l.ctx, &model.TransactionRecord{
		TransactionId:   transactionID,
		AccountFrom:     accountFrom,
		AccountTo:       sql.NullString{String: req.AccountTo, Valid: true},
		Amount:          req.Amount,
		TransactionType: vars.TransactionTypeTransfer,
		Status:          vars.TransactionStatusSuccess,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.TransferResponse{
		TransactionID: transactionID,
		AccountFrom:   accountFrom,
		AccountTo:     req.AccountTo,
		Amount:        req.Amount,
	}
	return
}
