package logic

import (
	"context"

	"github.com/Drengin6306/ZeroBank/service/account/rpc/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/account/rpc/proto"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CompensateAddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCompensateAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompensateAddLogic {
	return &CompensateAddLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CompensateAdd 收款补偿：将之前增加的金额扣回去
func (l *CompensateAddLogic) CompensateAdd(in *proto.CompensateAddRequest) (*proto.CompensateAddResponse, error) {
	var resp *proto.CompensateAddResponse

	// 使用事务和悲观锁保证并发安全
	err := l.svcCtx.AccountModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 使用悲观锁查询账户信息
		account, err := l.svcCtx.AccountModel.FindOneForUpdate(ctx, session, in.AccountId)
		if err != nil {
			return err
		}

		// 补偿：扣回之前增加的金额
		account.Balance -= in.Amount

		// 使用事务内的 session 更新账户信息
		err = l.svcCtx.AccountModel.WithSession(session).Update(ctx, account)
		if err != nil {
			return err
		}

		resp = &proto.CompensateAddResponse{
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
