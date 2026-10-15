local _M = {}

local shd_dict = ngx.shared["dyUpsMap"]
local json = require("cjson")
local printLog = 120

-- dict key
_M.dictUpCache = "upCache"
_M.dictCacheVer = "cacheVer"
-- 刷新相关变量
_M.upIdxMap = {}
_M.refreshLogSign = 9999                -- 首次一定全量打印
_M.refreshVer = 0
-- 当前work中(线程), up 信息缓存, 优化请求效率
-- upsCacheMap key: gameId; val: []{node:val,url:val,loadScore:val}
-- nodeCacheMap key: node; val: []ts (暂时只用到 ts, 没有放其他信息, 以后需要再来扩展)
_M.upsCacheMap = {}
_M.nodeCacheMap = {}
_M.cacheVer = 0
_M.workNum = 0

shd_dict:set(_M.dictCacheVer, 0)        -- 初始设置共享版本号为0



-- 初始化 -----------------------------------------
local refreshHandler = function (premature)
    if premature then
        return
    end
    _M.RefreshSvrList()
end

function _M.Init()
    -- 启动刷新服务器列表定时器
    -- 只需第一个worker跑的逻辑
    if ngx.worker.id() == 0 then
        -- 启动定时刷新
        local ok, err = ngx.timer.every(5, refreshHandler)
        if not ok then
            ngx.log(ngx.EMERG, "failed to create the timer: ", err)
            return
        end
    end

    _M.workNum = ngx.worker.count()
end

-- 分配服务器 -----------------------------------------
function _M.SelectPeer(gameId)
    -- 缓存 更新检测
    _M.CheckCache()

    -- 获取对应服务器列表
    local ups = _M.upsCacheMap[gameId]
    local upIdx = _M.upIdxMap[gameId] or 0
    if not ups or #ups == 0 then
        return nil, nil
    else
        -- 当前 worker 轮询
        if upIdx >= #ups then
            upIdx = 0
        end

        upIdx = upIdx + 1
        _M.upIdxMap[gameId] = upIdx
        local up = ups[upIdx]
        return up.url, up.node
    end
end

function _M.GetNodeTs(nodeName)
    -- 缓存 更新检测
    _M.CheckCache()

    local ts = _M.nodeCacheMap[nodeName]
    if ts then
        -- (暂时只用到 ts, 没有放其他信息, 以后需要再来扩展)
        return tonumber(ts)
    end

    return nil
end

-- 线程私有缓存 检测+更新
function _M.CheckCache()
    if _M.workNum == 1 then
        -- 单线程,直接赋值,不走共享缓存拷贝
        return
    end

    -- ups 更新检测
    local cacheVer = shd_dict:get(_M.dictCacheVer)
    if cacheVer ~= _M.cacheVer then
        -- 更新 ups 缓存
        local svrStr = shd_dict:get(_M.dictUpCache)
        _M.cacheVer = cacheVer
        local svrMap = json.decode(svrStr)
        if type(svrMap) == "table" then
            _M.upsCacheMap = svrMap[1]
            _M.nodeCacheMap = svrMap[2]
        else
           ngx.log(ngx.ERR, "illegal svrStr:", svrStr)
        end
        -- ngx.log(ngx.DEBUG, "update Ups; workId(", ngx.worker.id(),") ver(", _M.cacheVer ,")-------------",svrStr)
    end
end

-- 刷新服务器列表 -----------------------------------------
function _M.RefreshSvrList()
    -- 循环刷新列表中的所有GameId 相关服务器信息
    local svrMap = _M.RefreshSvrAll()
    if not svrMap then
        return
    end

    local svrStr = json.encode(svrMap)

    if _M.workNum == 1 then
        -- 单线程,直接赋值,不走共享缓存拷贝
        _M.upsCacheMap = svrMap[1]
        _M.nodeCacheMap = svrMap[2]
        return
    end

    -- 更新 共享缓存
    _M.refreshVer = _M.refreshVer + 1
    local ok, err, forcible = shd_dict:set(_M.dictUpCache, svrStr)
    if not ok or forcible then
        ngx.log(ngx.CRIT, "tp:shared; ", "up Cache err:[", ok, ",", err, ",", forcible)
        return
    end
    local ok, err, forcible = shd_dict:set(_M.dictCacheVer, _M.refreshVer)
    if not ok or forcible then
        ngx.log(ngx.CRIT, "tp:shared; ", "up Version err:[", ok, ",", err, ",", forcible)
    end
end

function _M.RefreshSvrAll()
    -- 打印 (10分钟 或 有更新时(logSign标记) 打印)
    local logAll = false
    _M.refreshLogSign = _M.refreshLogSign + 1
    if _M.refreshLogSign >= printLog then
        logAll = true
        _M.refreshLogSign = 0
    end

    -- 向 nacos 请求服务列表
    local nacos = require("resty.nacos")
    --local ip = require("dns").getIp(nacosCfg.domain, nacosCfg.domain)
    local ip = nacosCfg.domain
    local result = nacos.get_instance_list(ip, "Game", nacosCfg.namespaceId, nil, nil, true)
    if not result then
        ngx.log(ngx.CRIT, "nacos ret nil; ip:", ip, " ns:", nacosCfg.namespaceId)
        return nil
    end
    -- ngx.log(ngx.NOTICE, "RefreshSvrAll Suc:", result)


    -- json 解析
    local ok, body = pcall(json.decode, result)
    if not ok then
        ngx.log(ngx.CRIT, "解析 JSON 失败: ", result, " 错误: ", body)  -- body 此时为错误信息
        return nil  -- 终止本次刷新，避免程序中断
    end
    local host = body["hosts"]

    -- 数据组织
    local upsMap = {}
    local nodeMap = {}
    local changeMap = {}
    for _, v in ipairs(host) do
        -- 跳过不健康节点
        if not v.healthy then
            ngx.log(ngx.INFO, "skip unhealthy node: ", v.ip)
            goto continue
        end

        local gameId = v["metadata"]["gameId"]
        local node = v["metadata"]["node"]
        local ts = tonumber(v["metadata"]["ts"])
        local url = v["ip"] .. ":" .. v["port"]

        if gameId and node and ts then
            -- ups map
            local ups = upsMap[gameId]
            if not ups then
                -- 创建
                ups = {}
                upsMap[gameId] = ups
            end

            if not _M.nodeCacheMap[node] then
                -- 新增的
                changeMap[gameId] = gameId
            end

            -- 去重处理,应用ts 最新的
            local lastTs = nodeMap[node]
            if lastTs then
                if ts <= lastTs then
                    -- ts 更旧，跳过
                    goto continue
                else
                    -- 删除重复的
                    nodeMap[node] = nil
                    for i = 1, #ups do
                        if ups[i].node == node then
                            table.remove(ups, i)
                            break
                        end
                    end
                end
            end

            -- 数据插入
            table.insert(ups, {node=node, url=url})

            -- node map  (暂时只用到 ts, 没有放其他信息, 以后需要再来扩展)
            nodeMap[node] = ts
        else
            ngx.log(ngx.ERR, "illegal nacos host [gameId, node, ts]:[",gameId,", ",node,", ",ts,"]")
        end

        ::continue::
    end

    if logAll then
        -- 全量打印
        _M.PrintLogAll(upsMap, nodeMap)
    end

    for k, v in pairs(changeMap) do
        _M.PrintLog(v, upsMap, nodeMap)
    end

    local svrMap = {upsMap, nodeMap}
    return svrMap
end

function _M.PrintLog(gameId, upsMap, nodeMap)
    local ups = upsMap[gameId]

    if not ups then
        ngx.log(ngx.ERR, "PrintLog param ups nil; gameId:", gameId)
        return
    end

    -- 组织打印
    local nodeStr = {}
    for i, v in ipairs(ups) do
        local nodeTs = nodeMap[v.node]
        if not nodeTs then
            ngx.log(ngx.ERR, "PrintLog param nodeTs nil; gameId:", gameId, "; node:", v.node)
            return
        end
        table.insert(nodeStr, "(" .. v.node .. ":" .. nodeTs .. ")")
    end

    ngx.log(ngx.NOTICE, "tp:RefreshUpChange,", gameId, ", workNum:", _M.workNum, "; ", table.concat(nodeStr))
end

function _M.PrintLogAll(upsMap, nodeMap)
    local nodeStr = {}
    local logCount = 0

    for gameId, _ in pairs(upsMap) do
        local ups = upsMap[gameId]

        if type(ups) == "table" and #ups > 0 then 
            local currentNodes = {}
            for i, v in ipairs(ups) do
                local nodeTs = nodeMap[v.node]
                if nodeTs then
                    -- 记录节点总数，用于分段判断
                    logCount = logCount + 1
                    table.insert(currentNodes, "(" .. v.node .. ":" .. nodeTs .. ")")
                else
                    ngx.log(ngx.ERR, "PrintLog param nodeTs nil; gameId:", gameId, "; node:", v.node)
                end
            end

            -- 将当前游戏的所有节点拼成一个字符串单元
            -- 格式：gameId|(n1:ts1)(n2:ts2)
            if #currentNodes > 0 then
                local item = gameId .. "|" .. table.concat(currentNodes, "")
                table.insert(nodeStr, item)
            end

            -- 只有在一个 gameId 处理完后，才判断是否需要分段打印
            -- 这样可以保证同一个 gameId 的所有节点都在同一行日志中
            if logCount > 100 then
                ngx.log(ngx.NOTICE, "tp:RefreshUps, workNum:", _M.workNum, "; ", table.concat(nodeStr, ","))
                nodeStr = {}
                logCount = 0
            end
        end
    end

    -- 打印最后剩余的缓冲区数据
    if #nodeStr > 0 then
        ngx.log(ngx.NOTICE, "tp:RefreshUps, workNum:", _M.workNum, "; ", table.concat(nodeStr, ","))
    end
end



return _M