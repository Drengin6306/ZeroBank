// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package transaction

import (
	"net/http"

	"github.com/Drengin6306/ZeroBank/pkg/errorx"
	"github.com/Drengin6306/ZeroBank/pkg/response"
	"github.com/Drengin6306/ZeroBank/service/transaction/api/internal/logic/transaction"
	"github.com/Drengin6306/ZeroBank/service/transaction/api/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/transaction/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func TransferHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TransferRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Error(w, errorx.NewErrorWithMsg(errorx.ErrInvalidParams, err.Error()))
			return
		}

		l := transaction.NewTransferLogic(r.Context(), svcCtx)
		resp, err := l.Transfer(&req)
		if err != nil {
			response.Error(w, err)
		} else {
			response.Success(w, resp)
		}
	}
}
