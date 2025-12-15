package errorx

type ResCode uint32

const unknownMsg string = "未知错误"

const (
	Success ResCode = 200

	// 通用错误
	ErrInvalidParams  ResCode = 1000 + iota // 参数错误
	ErrNotLogin                             // 未登录
	ErrInvalidAccount                       // 账号或密码错误
	ErrForbidden                            // 权限不足禁止访问
	ErrNotFound                             // 资源不存在
	ErrServerBusy                           // 服务器繁忙
	ErrCustomerExists                       // 客户已存在
	ErrUnknown                              // 未知错误

	// 业务错误
	ErrAccountNotFound  ResCode = 2000 + iota // 账户不存在
	ErrAccountFrozen                          // 账户已冻结
	ErrBalanceNotEnough                       // 余额不足
	ErrTargetInvalid                          // 目标账户异常
	ErrAccountLimit                           // 账户限额
	ErrRiskControl                            // 风控拒绝
)
