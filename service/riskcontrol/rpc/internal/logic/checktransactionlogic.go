package logic

import (
	"context"
	"fmt"

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

func (l *CheckTransactionLogic) CheckTransaction(in *proto.RiskCheckRequest) (resp *proto.RiskCheckResponse, err error) {
	// 默认：通过
	passed := true
	reason := vars.RiskControlPassed

	var singleLimitAmount *mysql.LimitAmount
	var dailyLimitAmount *mysql.LimitAmount
	var limit float64

	switch in.TransactionType {
	case vars.TransactionTypeWithdraw:
		singleLimitAmount, err = l.svcCtx.LimitAmountModel.FindOneByAccountTypeLimitType(
			l.ctx, int64(in.AccountType), vars.RiskControlSingleWithdrawLimit,
		)
		if err != nil {
			goto done
		}

		dailyLimitAmount, err = l.svcCtx.LimitAmountModel.FindOneByAccountTypeLimitType(
			l.ctx, int64(in.AccountType), vars.RiskControlDailyWithdrawLimit,
		)
		if err != nil {
			goto done
		}

		logx.Debugf("Withdraw - Single Limit: %v, Daily Limit: %v", singleLimitAmount.Amount, dailyLimitAmount.Amount)

		if in.Amount > singleLimitAmount.Amount {
			passed = false
			reason = vars.RiskControlSingleWithdrawLimit
			limit = singleLimitAmount.Amount
			goto done
		}

		passed, err = l.svcCtx.RedisModel.CheckDailyWithdrawLimit(in.AccountFrom, in.Amount, dailyLimitAmount.Amount)
		if err != nil {
			goto done
		}
		if !passed {
			reason = vars.RiskControlDailyWithdrawLimit
			limit = dailyLimitAmount.Amount
			goto done
		}

	case vars.TransactionTypeTransfer:
		singleLimitAmount, err = l.svcCtx.LimitAmountModel.FindOneByAccountTypeLimitType(
			l.ctx, int64(in.AccountType), vars.RiskControlSingleTransferLimit,
		)
		if err != nil {
			goto done
		}

		dailyLimitAmount, err = l.svcCtx.LimitAmountModel.FindOneByAccountTypeLimitType(
			l.ctx, int64(in.AccountType), vars.RiskControlDailyTransferLimit,
		)
		if err != nil {
			goto done
		}

		if in.Amount > singleLimitAmount.Amount {
			passed = false
			reason = vars.RiskControlSingleTransferLimit
			limit = singleLimitAmount.Amount
			goto done
		}

		passed, err = l.svcCtx.RedisModel.CheckDailyTransferLimit(in.AccountFrom, in.Amount, dailyLimitAmount.Amount)
		if err != nil {
			goto done
		}
		if !passed {
			reason = vars.RiskControlDailyTransferLimit
			limit = dailyLimitAmount.Amount
			goto done
		}
	}

done:
	if err != nil {
		return nil, err
	}
	if !passed {
		_, err = l.svcCtx.RiskRecordModel.Insert(l.ctx, &mysql.RiskRecord{
			AccountId:     in.AccountFrom,
			TransactionId: in.TransactionId,
			RiskType:      int64(reason),
			Amount:        in.Amount,
		})
		if err != nil {
			return nil, err
		}
	}
	resp = &proto.RiskCheckResponse{
		Passed: passed,
		Reason: getReasonString(reason, limit),
	}
	return resp, nil
}

func getReasonString(reason int, amount float64) string {
	switch reason {
	case vars.RiskControlSingleWithdrawLimit:
		return fmt.Sprintf("单笔取款限额超限: %.2f", amount)
	case vars.RiskControlDailyWithdrawLimit:
		return fmt.Sprintf("每日取款限额超限: %.2f", amount)
	case vars.RiskControlSingleTransferLimit:
		return fmt.Sprintf("单笔转账限额超限: %.2f", amount)
	case vars.RiskControlDailyTransferLimit:
		return fmt.Sprintf("每日转账限额超限: %.2f", amount)
	case vars.RiskControlPassed:
		return "风控通过"
	default:
		return "No risk detected"
	}
}
