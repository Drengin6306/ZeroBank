package redis

import (
	"time"
)

// redis key prefix
const (
	RedisKeyPrefix              = "zerobank:"
	RedisDayTransferLimitPrefix = RedisKeyPrefix + "daily_quota:transfer:"
	RedisDayWithdrawLimitPrefix = RedisKeyPrefix + "daily_quota:withdraw:"
)

func getDayTransferLimitKey(accountID string) string {
	// key: zerobank:daily_quota:transfer:<accountID>:<date>
	date := time.Now().Format("20060102")
	return RedisDayTransferLimitPrefix + accountID + ":" + date
}

func getDayWithdrawLimitKey(accountID string) string {
	// key: zerobank:daily_quota:withdraw:<accountID>:<date>
	date := time.Now().Format("20060102")
	return RedisDayWithdrawLimitPrefix + accountID + ":" + date
}
