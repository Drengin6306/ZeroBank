-- KEYS[1]: 用户日限额的key
-- ARGV[1]: 本次交易金额（两位小数，如 "12.34"）
-- ARGV[2]: 日限额阈值（两位小数，如 "1000.00"）
-- 返回值: 1 表示允许交易（未超限），0 表示禁止交易（已超限）


local function money_to_cents(v)
    if v == nil then
        return nil
    end

    local n
    if type(v) == "number" then
        n = v
    else
        n = tonumber(v)
    end

    if n == nil then
        return nil
    end
    
    if (n ~= n) or (n == math.huge) or (n == -math.huge) then
        return nil
    end

    if n >= 0 then
        return math.floor(n * 100 + 0.5)
    else
        return math.ceil(n * 100 - 0.5)
    end
end

local existed = redis.call('EXISTS', KEYS[1]) == 1
local current_str = redis.call('GET', KEYS[1])
local current
if current_str == false or current_str == nil then
    current = 0
else
    current = tonumber(current_str)
    if current == nil then
        return redis.error_reply("invalid current cents: current=" .. tostring(current_str))
    end
end

local amount = money_to_cents(ARGV[1])
local limit = money_to_cents(ARGV[2])

if (amount == nil) or (limit == nil) then
    return redis.error_reply(
        "invalid money: amount=" .. tostring(ARGV[1]) ..
        ", limit=" .. tostring(ARGV[2]) ..
        ", current=" .. tostring(current_str)
    )
end

if current + amount > limit then
    return 0
else
    redis.call('INCRBY', KEYS[1], amount)

    -- 如果是新key，设置过期时间到当天23:59:59之后
    if not existed then
        local time = redis.call('TIME')
        local now = tonumber(time[1])
        local expire_seconds = 86400 - (now % 86400)
        redis.call('EXPIRE', KEYS[1], expire_seconds)
    end
    return 1
end
