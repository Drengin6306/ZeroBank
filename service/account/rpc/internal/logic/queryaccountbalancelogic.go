package logic

import (
	"context"

	"github.com/Drengin6306/ZeroBank/service/account/rpc/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/account/rpc/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryAccountBalanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryAccountBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryAccountBalanceLogic {
	return &QueryAccountBalanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryAccountBalanceLogic) QueryAccountBalance(in *proto.QueryAccountBalanceRequest) (*proto.QueryAccountBalanceResponse, error) {
	account, err := l.svcCtx.AccountModel.FindOne(l.ctx, in.AccountId)
	if err != nil {
		return nil, err
	}
	return &proto.QueryAccountBalanceResponse{
		AccountId: account.AccountId,
		Balance:   account.Balance,
	}, nil
}
