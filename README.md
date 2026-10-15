# golang

积累的一些 Golang 工具和库。

## 索引

| 项目 | 说明 |
| --- | --- |
| [nova-dev/unittest](./nova-dev/unittest) | Golang 单元测试工具、规范与示例（testify/mock、bytedance/mockey、Convey、覆盖率脚本） |
| [gneto-dev](./gneto-dev) | 基于 gnet 封装的高性能游戏服务器网络框架，支持 TCP/WebSocket 与 Protobuf/JSON |
| [gateway-dev](./gateway-dev) | 有状态游戏服务接入 K8s 集群的网关，基于 Nginx(OpenResty) 反向代理做连接路由 |

### [nova-dev/unittest](./nova-dev/unittest)

单元测试相关的工具说明、规范与可运行示例：

- [单元测试工具.md](./nova-dev/unittest/单元测试工具.md) — 测试初始化（`framen.TestMainRun`）、执行脚本、Mock 注入框架、Convey 用法
- [单元测试规范.md](./nova-dev/unittest/单元测试规范.md) — 覆盖率要求（≥80%）、命名与目录规范、表格驱动测试、CI/CD 集成
- [example/](./nova-dev/unittest/example) — 示例代码
  - [mockey/](./nova-dev/unittest/example/mockey) — bytedance/mockey 示例：函数/类注入、Convey 结合用法、Redis Mock
  - [testify/](./nova-dev/unittest/example/testify) — testify/mock 示例：函数 Mock、service/logic 分层测试

### [gneto-dev](./gneto-dev)

专为游戏业务“长连接、高并发、低延迟”场景优化的网络框架，基于对象重用池、零拷贝、线程亲和等特性：

- 编解码三层解耦：`transport`（TCP/WS 传输）→ `packet`（自定义包头）→ `payload`（Protobuf/JSON）
- 基于 Agent 的逻辑单线程化调度，同一玩家消息串行处理，业务 Handler 无需加锁
- 详见 [gneto-dev/README.md](./gneto-dev/README.md)，快速上手参考 [example/main.go](./gneto-dev/example/main.go)

### [gateway-dev](./gateway-dev)

服务于有状态游戏服务接入 K8s 集群的接入网关，基于 Nginx(OpenResty) 反向代理实现长连接路由：

- 通过 Lua 与 Nacos 服务发现感知游戏服实例上下线，连接建立时定向路由到目标游戏服
- 配合滚动发布、缩容场景的优雅退出，保证已有连接平滑迁移
- 详见 [gateway-dev/README.md](./gateway-dev/README.md)（含 OpenResty 安装、opm 依赖包、启动与运维命令）

