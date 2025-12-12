-- KEYS[1]: 用户日限额的key
-- ARGV[1]: 本次交易金额
-- ARGV[2]: 日限额阈值
-- 返回值: 1 表示允许交易（未超限），0 表示禁止交易（已超限）

local current = tonumber(redis.call('GET', KEYS[1]) or "0")
local amount = tonumber(ARGV[1])
local limit = tonumber(ARGV[2])

if current + amount > limit then
    return 0
else
    redis.call('INCRBY', KEYS[1], amount)
    -- 如果是新key，设置过期时间到当天23:59:59之后
    if current == 0 then
        -- 获取当前时间戳
        local time = redis.call('TIME')
        local now = tonumber(time[1])
        -- 计算当天剩余秒数：86400 - (当前时间戳 % 86400)
        local expire_seconds = 86400 - (now % 86400)
        redis.call('EXPIRE', KEYS[1], expire_seconds)
    end
    return 1
end
