// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package account

import (
	"context"
	"errors"

	"github.com/Drengin6306/ZeroBank/pkg/errorx"
	"github.com/Drengin6306/ZeroBank/pkg/format"
	"github.com/Drengin6306/ZeroBank/pkg/idgen"
	"github.com/Drengin6306/ZeroBank/pkg/password"
	"github.com/Drengin6306/ZeroBank/pkg/vars"
	"github.com/Drengin6306/ZeroBank/service/account/api/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/account/api/internal/types"
	"github.com/Drengin6306/ZeroBank/service/account/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 个人用户注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.CustomerRegisterRequest) (resp *types.RegisterResponse, err error) {
	// 用户是否已存在
	_, err = l.svcCtx.AccountModel.FindOne(l.ctx, req.IdCard)
	if err == nil {
		return nil, errorx.NewErrorWithMsg(errorx.ErrInvalidAccount, "用户已存在")
	}
	if !errors.Is(err, model.ErrNotFound) {
		return nil, err
	}
	// 个人用户注册
	phoneNum, err := format.Format(req.Phone)
	if err != nil {
		return nil, err
	}
	if IsValidEmail(req.Email) == false {
		return nil, errorx.NewErrorWithMsg(errorx.ErrInvalidParams, "邮箱格式错误")
	}
	customer := &model.CustomerIndividual{
		IdCard: req.IdCard,
		Name:   req.Name,
		Email:  req.Email,
		Phone:  phoneNum,
	}

	// 用户不存在则创建用户
	_, err = l.svcCtx.CustomerIndividualModel.FindOne(l.ctx, req.IdCard)
	if err == model.ErrNotFound {
		_, err = l.svcCtx.CustomerIndividualModel.Insert(l.ctx, customer)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		return nil, errorx.NewErrorWithMsg(errorx.ErrCustomerExists, "该身份证号已被注册")
	}

	// 创建账户
	accountID := idgen.GenAccountID()
	account := &model.Account{
		AccountId:   accountID,
		Password:    password.Encrypt(req.Password),
		CustomerId:  req.IdCard,
		AccountType: vars.AccountTypeIndividual,
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
