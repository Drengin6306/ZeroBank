// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package account

import (
	"net/http"

	"github.com/Drengin6306/ZeroBank/pkg/errorx"
	"github.com/Drengin6306/ZeroBank/pkg/response"
	"github.com/Drengin6306/ZeroBank/service/account/api/internal/logic/account"
	"github.com/Drengin6306/ZeroBank/service/account/api/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/account/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 用户登录
func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Error(w, errorx.NewErrorWithMsg(errorx.ErrInvalidParams, err.Error()))
			return
		}

		l := account.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			response.Error(w, err)
		} else {
			response.Success(w, resp)
		}
	}
}
