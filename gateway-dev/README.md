# gateway

## 说明

服务于**有状态游戏服务接入 K8s 集群**的接入网关。

游戏服务多为长连接且有状态，客户端需要稳定路由到指定的游戏服实例，无法直接使用普通的无状态负载均衡。本网关基于 **Nginx(OpenResty) 反向代理**实现连接路由：通过 Lua 与 Nacos 服务发现感知游戏服实例的上下线，在连接建立阶段将客户端连接定向路由到目标游戏服，并在游戏服滚动发布、缩容等场景下配合优雅退出，保证已有连接平滑迁移。

## 目录结构

```
gateway-dev/
├── conf/                       # Nginx 配置
│   ├── nginx.conf              # 主配置（server / location / 反向代理入口）
│   └── wss_proxy_params        # WebSocket(WSS) 反向代理参数
├── lua/                        # 业务 Lua 代码
│   ├── gameRouter.lua          # 游戏服路由入口：校验 accessToken(JWT)、playerId、gameId 并选择目标游戏服
│   ├── dyUps.lua               # 动态 upstream：基于 shared dict 与 worker 缓存，定时刷新游戏服节点列表
│   ├── redisTool.lua           # Redis 业务查询封装（requestId、游戏锁等，使用 eval 脚本）
│   ├── redisHelper.lua         # Redis 连接初始化与连接池管理
│   ├── tool.lua                # 通用工具函数
│   └── discard/                # 已废弃脚本（dns / redis_iresty / stopWatch）
├── lualib/
│   └── httpUtils.lua           # 修改版 lua-resty-nacos HTTP 工具（见“特殊修改”）
├── Dockerfile                  # 容器镜像构建
├── build.sh                    # 构建脚本
├── entrypoint.sh               # 容器启动入口
├── start.sh                    # 本地启动脚本（nginx -p run -c ../conf/nginx.conf）
└── nginxmd.sh                  # 调试辅助：列出 resty 已加载的包
```

## 安装
    基于 openresty
    mac 安装 brew install openresty

## 引入包
    opm get kuweiguge/lua-resty-nacos
    opm get SkyLothar/lua-resty-jwt
    opm get openresty/lua-resty-dns

### 引入包结果
    SkyLothar/lua-resty-jwt                                      0.1.11
    openresty/lua-resty-dns                                      0.21
    openresty/lua-resty-string                                   0.11
    kuweiguge/lua-resty-nacos                                    0.0.2
    jkeys089/lua-resty-hmac                                      0.06
    agentzh/lua-resty-http                                       0.09


## 特殊修改
    kuweiguge/lua-resty-nacos 的 httpUtils.lua 需修改 为 lualib/httpUtils.lua (http.new() 修改, 否则无法支持定时请求)

## 启动脚本
    start.sh

## 操作命令
    nginx -s reload
    nginx -s quit
    ps -ef | grep nginx | grep -v grep

## 优雅退出相关设置
    terminationGracePeriodSeconds(.yaml) 比 worker_shutdown_timeout(nginx.conf) 大10~20s,
    terminationGracePeriodSeconds 与 game服 保持一致
    lifecycle:
        preStop:
          exec:
            command: ["/bin/sh", "-c", "sleep 5"]  # 预留时间, 先断k8s导流,再启动Nginx终止