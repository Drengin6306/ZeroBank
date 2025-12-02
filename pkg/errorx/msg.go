package errorx

var resMap = map[ResCode]string{
	Success:             "成功",
	ErrInvalidParams:    "参数错误",
	ErrNotLogin:         "未登录",
	ErrInvalidAccount:   "账号或密码错误",
	ErrForbidden:        "权限不足禁止访问",
	ErrNotFound:         "资源不存在",
	ErrServerBusy:       "服务器繁忙",
	ErrUnknown:          "未知错误",
	ErrAccountNotFound:  "账户不存在",
	ErrAccountFrozen:    "账户已冻结",
	ErrBalanceNotEnough: "余额不足",
	ErrTargetInvalid:    "目标账户异常",
	ErrAccountLimit:     "账户限额",
	ErrRiskControl:      "风控拒绝",
}

type ResponseSuccess struct {
	Code    ResCode `json:"code"`
	Message string  `json:"message"`
	Data    any     `json:"data"`
}

type ResponseError struct {
	Code    ResCode `json:"code"`
	Message string  `json:"message"`
}
