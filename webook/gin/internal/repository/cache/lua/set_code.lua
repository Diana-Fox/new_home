local key=KEYS[1]
local cntKey=key..":cnt"
local val=ARGV[1]
local ttl=tonumber(redis.call("ttl",key))
if ttl == -1 then
    --key存在，但是没有过期时间
    return -2
elseif ttl == -2 or ttl < 540 then
    --存在，没有过期，且属于允许可以再次发送验证码的状态
    redis.call("set",key,val)
    redis.call("expire",key,600)--10分钟
    redis.call("set",cntKey,3)
    redis.call("expire",cntKey,600)--10分钟
    return 0
else
    --发送太频繁
    return -1
end