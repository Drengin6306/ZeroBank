package redis

import (
	"context"
	_ "embed"
	"time"
)

//go:embed quota.lua
var luaScript string

// CheckDailyTransferLimit 每日转账限额检查
func (m *defaultRedisModel) CheckDailyTransferLimit(accountID string, amount int64, dailyLimit int64) (bool, error) {
	key := getDayTransferLimitKey(accountID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 执行Lua脚本
	result, err := m.rdb.Eval(ctx, luaScript, []string{key}, amount, dailyLimit).Result()
	if err != nil {
		return false, err
	}
	allowed, ok := result.(int64)
	if !ok {
		return false, nil
	}
	return allowed == 1, nil
}

// CheckDailyWithdrawLimit 每日取款限额检查
func (m *defaultRedisModel) CheckDailyWithdrawLimit(accountID string, amount int64, dailyLimit int64) (bool, error) {
	key := getDayWithdrawLimitKey(accountID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 执行Lua脚本
	result, err := m.rdb.Eval(ctx, luaScript, []string{key}, amount, dailyLimit).Result()
	if err != nil {
		return false, err
	}
	allowed, ok := result.(int64)
	if !ok {
		return false, nil
	}
	return allowed == 1, nil
}
