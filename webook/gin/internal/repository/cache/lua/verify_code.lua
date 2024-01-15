local key=KEYS[1]
local expectedCode=ARGV[1]
local code=redis.call("get",key)
local cntKey=key..":cnt"
local cnt=tonumber(redis.call("get",cntKey))
if cnt<=0 then
    return -1--多次输入错误
elseif expectedCode==code then
    redis.call("del",key)
    redis.call("del",cntKey)
    return 0
else
    redis.call("decr",cntKey)
    return -2--输错了
end