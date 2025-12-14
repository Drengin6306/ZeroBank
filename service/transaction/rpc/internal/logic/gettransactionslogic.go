package logic

import (
	"context"
	"time"

	"github.com/Drengin6306/ZeroBank/pkg/errorx"
	"github.com/Drengin6306/ZeroBank/service/transaction/rpc/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/transaction/rpc/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTransactionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTransactionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTransactionsLogic {
	return &GetTransactionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取交易流水
func (l *GetTransactionsLogic) GetTransactions(in *proto.TransactionRequest) (*proto.TransactionResponse, error) {
	// 解析时间范围
	startTime, err := time.Parse(time.RFC3339, in.StartDate)
	if err != nil {
		return nil, errorx.NewError(errorx.ErrInvalidParams)
	}
	endTime, err := time.Parse(time.RFC3339, in.EndDate)
	if err != nil {
		return nil, errorx.NewError(errorx.ErrInvalidParams)
	}
	// 查询交易流水
	records, err := l.svcCtx.TransactionRecordModel.FindRange(l.ctx, in.AccountId, startTime, endTime)
	// 构造响应
	recordsProto := make([]*proto.Record, 0, len(records))
	for _, record := range records {
		recordsProto = append(recordsProto, &proto.Record{
			TransactionId:   record.TransactionId,
			AccountFrom:     record.AccountFrom,
			AccountTo:       &record.AccountTo.String,
			Amount:          record.Amount,
			TransactionType: int32(record.TransactionType),
			CreatedAt:       record.CreatedAt.Format(time.RFC3339),
		})
	}
	if err != nil {
		return nil, err
	}
	// 返回响应
	return &proto.TransactionResponse{
		Transactions: recordsProto,
	}, nil
}
