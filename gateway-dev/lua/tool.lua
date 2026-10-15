local _M = {}

local EnvProd = "prod"

function _M.init_by_lua_block()
    -- 全局配置 ------------
    nameserver = os.getenv("API_DNS")
    nacosCfg = {
        env = os.getenv("API_ENV"),
        id = os.getenv("API_ID"),
        domain = "http://" .. os.getenv("API_NACOS"),
        namespaceId = os.getenv("API_NACOS_NAMESPACE"),
    }
    openTest = os.getenv("OPEN_TEST")
    jwtSecret = os.getenv("JWT_SECRET")

    nacosCfg.namespaceId = _M.NacosNamespace(nacosCfg.env, nacosCfg.namespaceId, nacosCfg.id)
end

function _M.init_worker_by_lua_block()
    ngx.log(ngx.NOTICE, "init_worker(", ngx.worker.id(),")--- ",
        "[dns:", nameserver,
        "][env: ", nacosCfg.env,
        "][id: ", nacosCfg.id,
        "][nacos domain: ", nacosCfg.domain,
        "][nacos namespaceId: ", nacosCfg.namespaceId,
        "][openTest:", openTest, "]"
    )


    -- dyUps 初始化
    require("dyUps").Init()
end

function _M.NacosNamespace(env, namespace, id)
    if env == EnvProd then
        return "api"
    end
    if namespace then
        return namespace
    end
    if not id or id == "0" or id == 0 then
        return env .. "_api"
    else
        return env .. "_api_" .. id
    end
end

return _M