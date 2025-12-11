package logic

import (
	"context"
	"errors"

	"github.com/Drengin6306/ZeroBank/service/account/model"
	"github.com/Drengin6306/ZeroBank/service/account/rpc/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/account/rpc/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsAccountExistLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsAccountExistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsAccountExistLogic {
	return &IsAccountExistLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsAccountExistLogic) IsAccountExist(in *proto.AccountInfoRequest) (*proto.IsAccountExistResponse, error) {
	_, err := l.svcCtx.AccountModel.FindOne(l.ctx, in.AccountId)
	if err == nil {
		return &proto.IsAccountExistResponse{
			Exist: true,
		}, nil
	}
	if errors.Is(err, model.ErrNotFound) {
		return &proto.IsAccountExistResponse{
			Exist: false,
		}, nil
	}
	return nil, err
}
