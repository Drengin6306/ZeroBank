// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package account

import (
	"context"
	"time"

	"github.com/Drengin6306/ZeroBank/pkg/auth"
	"github.com/Drengin6306/ZeroBank/pkg/errorx"
	"github.com/Drengin6306/ZeroBank/pkg/password"
	"github.com/Drengin6306/ZeroBank/service/account/api/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/account/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// 用户登录
	account, err := l.svcCtx.AccountModel.FindOne(l.ctx, req.AccountID)
	if err != nil {
		return nil, err
	}
	if !password.Verify(req.Password, account.Password) {
		return nil, errorx.NewError(errorx.ErrInvalidAccount)
	}
	token, err := auth.GenerateToken(l.svcCtx.Config.Auth.AccessSecret, time.Now().Unix(), l.svcCtx.Config.Auth.AccessExpire, account.AccountId)
	if err != nil {
		return nil, err
	}
	resp = &types.LoginResponse{
		Token: token,
	}
	return
}
