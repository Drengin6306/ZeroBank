// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package account

import (
	"context"

	"github.com/Drengin6306/ZeroBank/pkg/vars"
	"github.com/Drengin6306/ZeroBank/service/account/api/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/account/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户信息
func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo() (resp *types.UserInfoResponse, err error) {
	// 获取account
	accountID := l.ctx.Value(vars.AccountKey).(string)
	account, err := l.svcCtx.AccountModel.FindOne(l.ctx, accountID)
	if err != nil {
		return nil, err
	}
	// 获取用户信息
	if account.AccountType == vars.AccountTypeIndividual {
		customer, err := l.svcCtx.CustomerIndividualModel.FindOne(l.ctx, account.CustomerId)
		if err != nil {
			return nil, err
		}
		resp = &types.UserInfoResponse{
			AccountType: vars.AccountTypeIndividual,
			UserInfo: &types.CustomerUserInfo{
				AccountID: account.AccountId,
				Name:      customer.Name,
				Phone:     customer.Phone,
				Email:     customer.Email,
			},
		}
	} else {
		customer, err := l.svcCtx.CustomerEnterpriseModel.FindOne(l.ctx, account.CustomerId)
		if err != nil {
			return nil, err
		}
		resp = &types.UserInfoResponse{
			AccountType: vars.AccountTypeEnterprise,
			UserInfo: &types.EnterpriseUserInfo{
				AccountID:  account.AccountId,
				Name:       customer.CompanyName,
				CreditCode: customer.CreditCode,
				LegalName:  customer.LegalName,
				Phone:      customer.Phone,
				Email:      customer.Email,
			},
		}
	}
	return
}
