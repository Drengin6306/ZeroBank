// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"github.com/Drengin6306/ZeroBank/pkg/errorx"
	"github.com/Drengin6306/ZeroBank/pkg/response"
	"github.com/Drengin6306/ZeroBank/service/report/api/internal/logic"
	"github.com/Drengin6306/ZeroBank/service/report/api/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/report/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 生成报表
func getReportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GenerateReportRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Error(w, errorx.NewErrorWithMsg(errorx.ErrInvalidParams, err.Error()))
			return
		}

		l := logic.NewGetReportLogic(r.Context(), svcCtx)
		resp, err := l.GetReport(&req)
		if err != nil {
			response.Error(w, err)
		} else {
			response.Success(w, resp)
		}
	}
}
