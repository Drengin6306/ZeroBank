package vars

const (
	AccountTypeIndividual = 0
	AccountTypeEnterprise = 1

	AccountStatusActive = 0
	AccountStatusFrozen = 1

	AccountKey = "account"

	TransactionTypeTransfer = 0
	TransactionTypeDeposit  = 1
	TransactionTypeWithdraw = 2

	TransactionStatusFailed  = 0
	TransactionStatusSuccess = 1

	RiskControlAccountFrozen       = 1 // 账户冻结
	RiskControlDailyTransferLimit  = 2 // 日转账限额
	RiskControlSingleTransferLimit = 3 // 单笔转账限额
	RiskControlDailyWithdrawLimit  = 4 // 日提现限额
	RiskControlSingleWithdrawLimit = 5 // 单笔提现限额
)
