// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package transaction

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Drengin6306/ZeroBank/pkg/errorx"
	"github.com/Drengin6306/ZeroBank/pkg/idgen"
	"github.com/Drengin6306/ZeroBank/pkg/vars"
	"github.com/Drengin6306/ZeroBank/service/account/rpc/account"
	"github.com/Drengin6306/ZeroBank/service/riskcontrol/rpc/riskcontrol"
	"github.com/Drengin6306/ZeroBank/service/transaction/api/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/transaction/api/internal/types"
	"github.com/Drengin6306/ZeroBank/service/transaction/model"
	"github.com/dtm-labs/client/dtmgrpc"

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

	// 获取账户信息用于风控检查
	info, err := l.svcCtx.AccountRpc.GetAccountInfo(l.ctx, &account.AccountInfoRequest{
		AccountId: accountFrom,
	})
	if err != nil {
		return nil, err
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

	// 使用 DTM SAGA 模式确保分布式事务一致性
	//
	// 注意：DTM 官方镜像不支持 Consul 服务发现驱动
	// 在 Docker 环境中使用容器服务名直接通信
	// Docker Compose 会自动进行 DNS 解析
	accountRpcTarget := "account-rpc:9001"

	// 创建 SAGA 事务
	saga := dtmgrpc.NewSagaGrpc(l.svcCtx.Config.DTM.Server, transactionID)

	// 添加扣款步骤（正向操作 + 补偿操作）
	saga.Add(
		accountRpcTarget+"/account.Account/DeductBalance",    // 正向操作：扣款
		accountRpcTarget+"/account.Account/CompensateDeduct", // 补偿操作：加回扣除的金额
		&account.DeductBalanceRequest{
			AccountId: accountFrom,
			Amount:    req.Amount,
		},
	)

	// 添加收款步骤（正向操作 + 补偿操作）
	saga.Add(
		accountRpcTarget+"/account.Account/AddBalance",    // 正向操作：收款
		accountRpcTarget+"/account.Account/CompensateAdd", // 补偿操作：扣回增加的金额
		&account.AddBalanceRequest{
			AccountId: req.AccountTo,
			Amount:    req.Amount,
		},
	)

	// 设置等待结果并提交
	saga.WaitResult = true
	err = saga.Submit()
	if err != nil {
		l.Logger.Errorf("DTM SAGA 提交失败: %v", err)
		return nil, fmt.Errorf("转账失败: %w", err)
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
