package response

import (
	"net/http"

	"github.com/Drengin6306/ZeroBank/pkg/errorx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func Success(w http.ResponseWriter, data interface{}) {
	msg := errorx.NewSuccess()
	msg.Data = data
	httpx.WriteJson(w, http.StatusOK, msg)
}

func Error(w http.ResponseWriter, err error) {
	if e, ok := err.(*errorx.ResponseError); ok {
		httpx.WriteJson(w, http.StatusOK, e)
		return
	}
	err = errorx.NewError(errorx.ErrServerBusy)
	httpx.WriteJson(w, http.StatusOK, err)
}
