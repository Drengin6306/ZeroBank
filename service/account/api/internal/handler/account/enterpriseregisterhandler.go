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

// 企业用户注册
func EnterpriseRegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EnterpriseRegisterRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Error(w, errorx.NewErrorWithMsg(errorx.ErrInvalidParams, err.Error()))
			return
		}

		l := account.NewEnterpriseRegisterLogic(r.Context(), svcCtx)
		resp, err := l.EnterpriseRegister(&req)
		if err != nil {
			response.Error(w, err)
		} else {
			response.Success(w, resp)
		}
	}
}
