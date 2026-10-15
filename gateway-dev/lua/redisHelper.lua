local _M = {}

local redis_c = require "resty.redis"
local rdsCfg = {}

function _M.Init(cfg)
    cfg = cfg or {}
    
    
	local ip = cfg.ip or '127.0.0.1'
	local port = cfg.port or 6379
	local password = cfg.password

    -- local db_index= opts.db_index or 0
    local timeout = (cfg.timeout and cfg.timeout * 1000) or 1000
	local pool_max_idle_time = cfg.pool_max_idle_time or 60000
	local pool_size = cfg.pool_size or 100


    rdsCfg = {
	    ip = ip,
	    port = port,
	    password = password,
        timeout = timeout,                          -- 操作响应超时时间
	    pool_max_idle_time = pool_max_idle_time,    -- 连接池空闲保留时间
	    pool_size = pool_size                       -- 连接池大小
    }
end

-- 创建或从连接池获取链接
-- return redis
function _M.Connect()
    local redis, err = redis_c:new()
    if not redis then
        ngx.log(ngx.ERR, "new redis_c error : ", err)
        return nil
    end

    local ok, err = redis:connect(rdsCfg.ip, rdsCfg.port)
    if not ok then  
		ngx.log(ngx.ERR, "connect redis error : ", err)
        _M.Close(redis)
		return nil
	end

    ----密码认证
    if rdsCfg.password then
		local count, err = redis:get_reused_times()
		if 0 == count then ----新建连接，需要认证密码
			ok, err = redis:auth(rdsCfg.password)
			if not ok then
                ngx.log(ngx.ERR, "failed to auth redis: ", err)
                return nil
			end
		elseif err then
            ngx.log(ngx.ERR, "failed to get redis reused times : ", err)
			return nil
        else
            ----从连接池中获取连接，无需再次认证密码
		end
	end

    return redis
end

-- 将连接还回连接池
function _M.ReturnPool(redis)
    if not redis then  
        return  
    end  
    --释放连接(连接池实现)
    local ok, err = redis:set_keepalive(rdsCfg.pool_max_idle_time, rdsCfg.pool_size)  
    if not ok then
        ngx.log(ngx.ERR, "redis et keepalive error : ", err)
    end  
end

-- 关闭连接
function _M.Close(redis)
    if not redis then
        return
    end
    local ok, err = redis:close()
    if not ok then
        ngx.log(ngx.ERR, "close redis error : ", err)
    end
end

function _M.Call(callback)
    local redis = _M.Connect()
    if not redis then
        return
    end
    
    -- redis:select(rdsCfg.db_index)
    local ret = callback(redis)

    _M.ReturnPool(redis)
    return ret
end


return _M