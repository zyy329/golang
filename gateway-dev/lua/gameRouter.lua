local _M = {}

local BASE_ERR = 700
local ERR_PARAM = BASE_ERR           -- 请求参数错误
local ERR_JWT = BASE_ERR + 1         -- JWT错误
local ERR_JWT_PARAM = BASE_ERR + 2   -- JWT载荷参数错误; 缺失 playerId 或 requestId
local ERR_UID = BASE_ERR + 3         -- JWT载荷 和 url 中的 playerId 不一致
local ERR_REQUEST_ID = BASE_ERR + 4  -- requestId 已失效
local ERR_INNER = BASE_ERR + 5       -- 服务器内部错误
local ERR_SHUTDOWN = BASE_ERR + 6    -- 停服维护
local ERR_LOCK = BASE_ERR + 7        -- 金币锁锁定中, 不可进入其他游戏

local json = require("cjson.safe")
local dyUps = require("dyUps")
local jwt = require "resty.jwt"
local rdsT = require("redisTool")

-- 游戏服路由
function _M.Router()
    -- 参数 检查
    -- 获取请求参数
    --local wsParam = ngx.var.http_sec_websocket_protocol
    --local tb = _M.ParasWsParam(wsParam)

    local token = ngx.var.arg_accessToken
    local playerId = ngx.var.arg_playerId
    local gameId = ngx.var.arg_gameId

    -- 进入打印, 以后不需要了可考虑屏蔽掉
    ngx.log(ngx.NOTICE, "gameRouter in")

    if not token or not playerId or not gameId then
        ngx.log(ngx.WARN, "gameRouter wsParam err:["
        , token ~= nil, "|"
        , playerId ~= nil, "|"
        , gameId ~= nil, "]"
        )
        ngx.exit(ERR_PARAM)
    end


    -- AccessToken 解密
    local requestId = _M.ParseToken(playerId, token)
    if not requestId then
        -- 错误信息已打印, 不用重复打印
        ngx.exit(ERR_JWT)
    end

    --ngx.log(ngx.DEBUG, string.format("jwt param: [playerId:%s],[requestId:%s],[gameId:%s]", playerId, requestId, gameId))

    if not rdsT.InitRds() then
        -- redis 初始化失败; 已打印错误信息, 这里不用再进行打印
        ngx.exit(ERR_INNER)
        return
    end

    -- redis requestId, GameLock 获取
    local rds_requestId, rds_gameLockUrl, rds_node = rdsT.GetRequestIdAndGameLockUrl(playerId, gameId)

    -- requestId 校验
    if requestId ~= rds_requestId then
        ngx.log(ngx.WARN, "player(", playerId ,") requestId failed:", requestId, ":", rds_requestId)
        ngx.exit(ERR_REQUEST_ID)
    end

    -- 断线重连地址
    if rds_gameLockUrl then
        ngx.log(ngx.NOTICE, "player(", playerId ,"), game(", gameId,"), node(", rds_node,") reconnect:", rds_gameLockUrl)
        return rds_gameLockUrl
    end

	-- 利用 动态upstream 进行负载均衡分配
    local selectUrl, node = dyUps.SelectPeer(gameId)
    if not selectUrl then
        ngx.log(ngx.ERR, "no active server(", gameId ,"); player(", playerId ,")")
        ngx.exit(ERR_SHUTDOWN)
    end
    ngx.log(ngx.NOTICE, "tp:RouterSuc;  node:", node, "; url:", selectUrl)
    return selectUrl
end

function _M.ParseToken(playerId, token)
    --ngx.log(ngx.DEBUG, "jwtSecret:", jwtSecret, "; type:", type(jwtSecret))
    local jwt_obj = jwt:verify(jwtSecret, token)
    --ngx.log(ngx.DEBUG, "jwt:", json.encode(jwt_obj))

    if jwt_obj["verified"] and jwt_obj["valid"] and jwt_obj["payload"]["content"] then
        local jwtC = json.decode(jwt_obj["payload"]["content"])
        return _M.CheckJwt(playerId, jwtC)
    else
        ngx.log(ngx.WARN, "jwt err; [verified:", jwt_obj["verified"], "][valid:", jwt_obj["valid"], "][reason:",jwt_obj["reason"], "]")
        --ngx.log(ngx.WARN, "jwt:", json.encode(jwt_obj))
        return nil
    end
end

function _M.CheckJwt(srcplayerId, jwtC)
    local playerId = jwtC["playerId"]
    local requestId = jwtC["requestId"]
    if not playerId or not requestId then
        ngx.log(ngx.WARN, string.format("jwt param err: [playerId:%s],[requestId:%s]",
                playerId, requestId))
        ngx.exit(ERR_JWT_PARAM)
    end

    if playerId ~= srcplayerId then
        ngx.log(ngx.WARN, string.format("jwt playerId err: [playerId:%s/%s],[requestId:%s]",
                playerId, srcplayerId, requestId))
        ngx.exit(ERR_UID)
    end

    return requestId
end

function _M.Test()
    if not rdsT.InitRds() then
        -- redis 初始化失败; 已打印错误信息, 这里不用再进行打印
        ngx.exit(ERR_INNER)
        return
    end

    local playerId = "4780308078592000"
    local requestId, addr, gameId = rdsT.GetRequestIdAndGameLockUrl(playerId, 5324)
    ngx.log(ngx.ERR, "id:", playerId, "; requestId:", requestId, "; addr:", addr, "; gameId:", gameId)
end

return _M
