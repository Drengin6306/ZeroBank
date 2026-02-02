package logic

import (
	"context"

	"github.com/Drengin6306/ZeroBank/pkg/errorx"
	"github.com/Drengin6306/ZeroBank/service/account/rpc/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/account/rpc/proto"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DeductBalanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductBalanceLogic {
	return &DeductBalanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeductBalanceLogic) DeductBalance(in *proto.DeductBalanceRequest) (*proto.DeductBalanceResponse, error) {
	var resp *proto.DeductBalanceResponse

	// 使用事务和悲观锁保证并发安全
	err := l.svcCtx.AccountModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 使用悲观锁查询账户信息
		account, err := l.svcCtx.AccountModel.FindOneForUpdate(ctx, session, in.AccountId)
		if err != nil {
			return err
		}

		// 检查余额是否充足
		if account.Balance < in.Amount {
			return errorx.NewError(errorx.ErrBalanceNotEnough)
		}

		// 扣减余额
		account.Balance -= in.Amount

		// 使用事务内的 session 更新账户信息
		err = l.svcCtx.AccountModel.WithSession(session).Update(ctx, account)
		if err != nil {
			return err
		}

		resp = &proto.DeductBalanceResponse{
			AccountId: account.AccountId,
			Balance:   account.Balance,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}
