package response

import (
	"net/http"

	"github.com/Drengin6306/ZeroBank/pkg/errorx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func Success(r *http.Request, w http.ResponseWriter, data interface{}) {
	msg := errorx.NewSuccess()
	msg.Data = data
	httpx.WriteJson(w, http.StatusOK, msg)
}

func Error(r *http.Request, w http.ResponseWriter, code errorx.ResCode) {
	msg := errorx.NewError(code)
	httpx.WriteJson(w, http.StatusInternalServerError, msg)
}
