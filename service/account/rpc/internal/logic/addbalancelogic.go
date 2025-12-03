package logic

import (
	"context"

	"github.com/Drengin6306/ZeroBank/service/account/rpc/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/account/rpc/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddBalanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddBalanceLogic {
	return &AddBalanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddBalanceLogic) AddBalance(in *proto.AddBalanceRequest) (*proto.AddBalanceResponse, error) {
	account, err := l.svcCtx.AccountModel.FindOne(l.ctx, in.AccountId)
	if err != nil {
		return nil, err
	}
	account.Balance += in.Amount
	err = l.svcCtx.AccountModel.Update(l.ctx, account)
	if err != nil {
		return nil, err
	}
	resp := &proto.AddBalanceResponse{
		AccountId: account.AccountId,
		Balance:   account.Balance,
	}
	return resp, nil
}
