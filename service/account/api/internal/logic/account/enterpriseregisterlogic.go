// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package account

import (
	"context"

	"github.com/Drengin6306/ZeroBank/pkg/format"
	"github.com/Drengin6306/ZeroBank/pkg/idgen"
	"github.com/Drengin6306/ZeroBank/pkg/password"
	"github.com/Drengin6306/ZeroBank/pkg/vars"
	"github.com/Drengin6306/ZeroBank/service/account/api/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/account/api/internal/types"
	"github.com/Drengin6306/ZeroBank/service/account/model"
	"github.com/zeromicro/go-zero/core/logx"
)

type EnterpriseRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 企业用户注册
func NewEnterpriseRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnterpriseRegisterLogic {
	return &EnterpriseRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnterpriseRegisterLogic) EnterpriseRegister(req *types.EnterpriseRegisterRequest) (resp *types.RegisterResponse, err error) {
	// 企业用户注册
	phoneNum, err := format.Format(req.Phone)
	if err != nil {
		return nil, err
	}
	enterprise := &model.CustomerEnterprise{
		CompanyName: req.Name,
		CreditCode:  req.CreditCode,
		LegalName:   req.LegalName,
		LegalIdCard: req.LegalIdCard,
		Phone:       phoneNum,
		Email:       req.Email,
	}
	_, err = l.svcCtx.CustomerEnterpriseModel.Insert(l.ctx, enterprise)
	if err != nil {
		return nil, err
	}
	accountID := idgen.GenAccountID()
	account := &model.Account{
		AccountId:   accountID,
		Password:    password.Encrypt(req.Password),
		AccountType: vars.AccountTypeEnterprise,
		Status:      vars.AccountStatusActive,
	}
	_, err = l.svcCtx.AccountModel.Insert(l.ctx, account)
	if err != nil {
		return nil, err
	}
	resp = &types.RegisterResponse{
		AccountID: accountID,
	}
	return
}
