// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package account

import (
	"net/http"

	"github.com/Drengin6306/ZeroBank/pkg/response"
	"github.com/Drengin6306/ZeroBank/service/account/api/internal/logic/account"
	"github.com/Drengin6306/ZeroBank/service/account/api/internal/svc"
)

// 获取用户余额
func GetBalanceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := account.NewGetBalanceLogic(r.Context(), svcCtx)
		resp, err := l.GetBalance()
		if err != nil {
			response.Error(w, err)
		} else {
			response.Success(w, resp)
		}
	}
}
