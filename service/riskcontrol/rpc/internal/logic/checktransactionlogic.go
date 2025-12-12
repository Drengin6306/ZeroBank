package logic

import (
	"context"

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
	//

	return &proto.RiskCheckResponse{}, nil
}
