// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"
	"net/url"
	"os"

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
		}
		filePath := "temp/" + resp.FileName

		defer func() {
			if _, err := os.Stat(filePath); err == nil {
				os.Remove(filePath)
			}
		}()

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			response.Error(w, errorx.NewError(errorx.ErrServerBusy))
			return
		}

		// 使用 QueryEscape 防止文件名包含中文或特殊字符导致乱码
		encodedFilename := url.QueryEscape(resp.FileName)
		w.Header().Set("Content-Disposition", "attachment; filename*=UTF-8''"+encodedFilename)
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

		http.ServeFile(w, r, filePath)
	}
}
