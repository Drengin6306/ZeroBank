package logic

import (
	"context"

	"github.com/Drengin6306/ZeroBank/pkg/vars"
	"github.com/Drengin6306/ZeroBank/service/riskcontrol/model/mysql"
	"github.com/Drengin6306/ZeroBank/service/riskcontrol/rpc/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/riskcontrol/rpc/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckTransactionLogic {
	return &CheckTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckTransactionLogic) CheckTransaction(in *proto.RiskCheckRequest) (*proto.RiskCheckResponse, error) {
	var (
		singleLimitAmount *mysql.LimitAmount
		dailyLimitAmount  *mysql.LimitAmount
		err               error
		passed            bool = true
	)
	if in.TransactionType == vars.TransactionTypeWithdraw {
		// 提现交易类型的限额检查逻辑
		singleLimitAmount, err = l.svcCtx.LimitAmountModel.FindOneByAccountTypeLimitType(l.ctx, int64(in.AccountType), vars.RiskControlSingleWithdrawLimit)
		if err != nil {
			return nil, err
		}
		dailyLimitAmount, err = l.svcCtx.LimitAmountModel.FindOneByAccountTypeLimitType(l.ctx, int64(in.AccountType), vars.RiskControlDailyWithdrawLimit)
		if err != nil {
			return nil, err
		}
		logx.Debugf("Withdraw - Single Limit: %v, Daily Limit: %v", singleLimitAmount.Amount, dailyLimitAmount.Amount)
		if in.Amount > singleLimitAmount.Amount {
			passed = false
			return &proto.RiskCheckResponse{
				Passed: passed,
				Reason: vars.RiskControlSingleWithdrawLimit,
			}, nil
		}
		passed, err = l.svcCtx.RedisModel.CheckDailyWithdrawLimit(in.AccountFrom, in.Amount, dailyLimitAmount.Amount)
		if err != nil {
			return nil, err
		}
		if !passed {
			return &proto.RiskCheckResponse{
				Passed: passed,
				Reason: vars.RiskControlDailyWithdrawLimit,
			}, nil
		}
	} else if in.TransactionType == vars.TransactionTypeTransfer {
		// 转账交易类型的限额检查逻辑
		singleLimitAmount, err = l.svcCtx.LimitAmountModel.FindOneByAccountTypeLimitType(l.ctx, int64(in.AccountType), vars.RiskControlSingleTransferLimit)
		if err != nil {
			return nil, err
		}
		dailyLimitAmount, err = l.svcCtx.LimitAmountModel.FindOneByAccountTypeLimitType(l.ctx, int64(in.AccountType), vars.RiskControlDailyTransferLimit)
		if err != nil {
			return nil, err
		}
		if in.Amount > singleLimitAmount.Amount {
			passed = false
			return &proto.RiskCheckResponse{
				Passed: passed,
				Reason: vars.RiskControlSingleTransferLimit,
			}, nil
		}
		passed, err = l.svcCtx.RedisModel.CheckDailyTransferLimit(in.AccountFrom, in.Amount, dailyLimitAmount.Amount)
		if err != nil {
			return nil, err
		}
		if !passed {
			return &proto.RiskCheckResponse{
				Passed: passed,
				Reason: vars.RiskControlDailyTransferLimit,
			}, nil
		}
	}
	return &proto.RiskCheckResponse{
		Passed: true,
		Reason: vars.RiskControlPassed,
	}, nil
}
