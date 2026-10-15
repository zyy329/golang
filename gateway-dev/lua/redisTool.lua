local _M = {}

local dyUps = require("dyUps")
local mRdsH = require("redisHelper")
local json = require("cjson.safe")

local isInit = false
local scriptMap = {}

-- 获取指定玩家 requestId, gameLock 信息
-- 注意 requestId, game_lock:player, 需在同一个实例中, 不支持分片(eval 命令)
local IDX_requestId = 1
local IDX_lockJson = 2
local Key_requestId_and_gameLock = "requestId_and_gameLock"
local script_requestId_and_gameLock = "\
    local playerId = KEYS[1]\
    local gameId = ARGV[1]\
    local k1 = 'lobby:login:requestId:{' .. playerId .. '}:' .. gameId\
    local k2 = 'lobby:login:game_lock:player:{' .. playerId .. '}'\
    local requestId = redis.pcall('get', k1)\
    local lockJson = redis.pcall('hget', k2, gameId)\
    return {requestId, lockJson}\
"



function _M.InitRds()
    if isInit then
        return true
    end

    -- 通过nacos 获取redis配置信息
    local nacos = require("resty.nacos")
    local result = nacos.get_config(nacosCfg.domain, 'redis', nacosCfg.namespaceId, nil)

    local body = json.decode(result)
    if not body then
        ngx.log(ngx.CRIT, "get rdsCfg err: ", result)
        return false
    end

    local redisCfg = body["Common"]
    if not redisCfg or type(redisCfg) ~= "table" or not redisCfg["Addr"] then
        ngx.log(ngx.CRIT, "get rdsCfg info err; ", result)
        return false
    end

    -- 初始化
    local rdsCfg = {
        password = redisCfg["Password"],
        db_index = 0,
        timeout = 10
    }
    local arr = _M.SplitStr(redisCfg["Addr"], ":")
    rdsCfg.ip, rdsCfg.port = arr[1], arr[2]
    ngx.log(ngx.NOTICE, "rdsCfg: ip[", rdsCfg.ip, "]  port:", rdsCfg.port)
    mRdsH.Init(rdsCfg)

    -- 脚本初始化
    local ok = _M.InitScript()
    if ok then
        isInit = true
    end
    return ok
end

function _M.InitScript()
    return mRdsH.Call(function(redis)
        local hash = redis:script("load", script_requestId_and_gameLock)
        scriptMap[Key_requestId_and_gameLock] = hash
        ngx.log(ngx.NOTICE, "Script Hash:", hash)
        return hash ~= nil
    end)
end

function _M.GetRequestIdAndGameLockUrl(playerId, gameId)
    --ngx.log(ngx.DEBUG, "Script Hash:", scriptMap[Key_requestId_and_gameLock], " k1:", k1, " k2:", k2)
    local retArr = mRdsH.Call(function(redis)
        return redis:evalsha(scriptMap[Key_requestId_and_gameLock], 1, playerId, gameId)
    end)

    -- ngx.log(ngx.NOTICE, "retArr:", json.encode(retArr))
    
    local requestId, lockJson, addr, node, ts, svrStart = nil, nil, nil, nil, nil, nil
    if retArr then
        requestId = retArr[IDX_requestId]
        lockJson = retArr[IDX_lockJson]
        -- ngx.log(ngx.NOTICE, "requestId:", requestId)
        -- ngx.log(ngx.NOTICE, "lockJson:", lockJson)
        -- json 解析
        local lockInfo = json.decode(lockJson)
        if not lockInfo then
            -- ngx.log(ngx.CRIT, "decode err; lockJson:{", lockJson, "} player:{", playerId,"}")
            return requestId, nil, nil
        end
        addr = lockInfo["addr"]
        node = lockInfo["node"]
        ts = lockInfo["ts"]
        -- ngx.log(ngx.NOTICE, "addr:", addr)
        -- ngx.log(ngx.NOTICE, "node:", node)
        -- ngx.log(ngx.NOTICE, "ts:", ts)
        
        if addr then 
            if not node or not ts then
                addr = nil
            end

            svrStart = dyUps.GetNodeTs(node)
            -- ngx.log(ngx.NOTICE, "ts:", ts, "; s:", svrStart)
            if not svrStart or tonumber(ts) < svrStart then
                -- ngx.log(ngx.NOTICE, "svrStart:",svrStart,"; ts:", ts)
                -- 移除已失效的 金币锁; (Game服务已不再 或 GameLock创建时间早于对应Game启动时间)
                mRdsH.Call(function(redis)
                    local num, err = redis:hdel('lobby:login:game_lock:player:{' .. playerId .. '}', gameId)
                    ngx.log(ngx.NOTICE, "del lose efficacy gameLock; result:[", num, ";", err, "]player:{", playerId,"}, node:{", node,"} ts:{", ts,"/", svrStart,"}")
                end)
                addr = nil
            end
        end
    end

    return requestId, addr, node
end

function _M.SplitStr(str, delimiter)
    local arr = {}

    for v in string.gmatch(str, "([^" .. delimiter .. "]+)") do
        arr[#arr + 1] = v
    end

    return arr
end

function _M.ParasWsParam(param)
    local tb = {}

    -- 去除空格
    param = param:gsub("%s+", "")

    -- 分割解析
    local arr = _M.SplitStr(param, "|,")
    for i=1,#arr-1,2 do
        --ngx.log(ngx.DEBUG, "loop", i, "|", arr[i], "|", arr[i+1])
        tb[arr[i]] = arr[i+1]
    end
    --ngx.log(ngx.DEBUG, #arr, "  ", table.concat(arr))
    --for k, v in pairs(tb) do
    --    ngx.log(ngx.DEBUG, "tb_", k, ":", v)
    --end

    return tb
end

return _M