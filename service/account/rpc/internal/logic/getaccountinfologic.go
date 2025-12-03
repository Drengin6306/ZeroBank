package logic

import (
	"context"

	"github.com/Drengin6306/ZeroBank/pkg/vars"
	"github.com/Drengin6306/ZeroBank/service/account/rpc/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/account/rpc/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAccountInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAccountInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccountInfoLogic {
	return &GetAccountInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAccountInfoLogic) GetAccountInfo(in *proto.AccountInfoRequest) (*proto.AccountInfoResponse, error) {
	// 获取account
	account, err := l.svcCtx.AccountModel.FindOne(l.ctx, in.AccountId)
	if err != nil {
		return nil, err
	}
	// 获取用户信息
	var resp *proto.AccountInfoResponse
	if account.AccountType == vars.AccountTypeIndividual {
		customer, err := l.svcCtx.CustomerIndividualModel.FindOne(l.ctx, account.CustomerId)
		if err != nil {
			return nil, err
		}
		resp = &proto.AccountInfoResponse{
			AccountType: vars.AccountTypeIndividual,
			Balance:     account.Balance,
			UserInfo: &proto.AccountInfoResponse_CustomerUserInfo{
				CustomerUserInfo: &proto.CustomerUserInfo{
					AccountID: account.AccountId,
					Name:      customer.Name,
					Phone:     customer.Phone,
					Email:     customer.Email,
				},
			},
		}
	} else {
		customer, err := l.svcCtx.CustomerEnterpriseModel.FindOne(l.ctx, account.CustomerId)
		if err != nil {
			return nil, err
		}
		resp = &proto.AccountInfoResponse{
			AccountType: vars.AccountTypeEnterprise,
			Balance:     account.Balance,
			UserInfo: &proto.AccountInfoResponse_EnterpriseUserInfo{
				EnterpriseUserInfo: &proto.EnterpriseUserInfo{
					AccountID:  account.AccountId,
					Name:       customer.CompanyName,
					CreditCode: customer.CreditCode,
					LegalName:  customer.LegalName,
					Phone:      customer.Phone,
					Email:      customer.Email,
				},
			},
		}
	}
	return resp, nil
}
