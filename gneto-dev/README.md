# gneto

gneto 是一个基于 gnet 封装的高性能网络框架，专为游戏服务器设计。
它提供 Tcp/WebSocket 和 Protobuf/Json 支持，
利用 重用池,零拷贝,线程亲和等特性 针对游戏业务所需的“长连接、高并发、低延迟”特性进行了优化。

## 目录结构

```
gneto/
├── example/            # 示例代码与基准测试
├── codec/              # 编解码层 (分层解耦协议栈)
│   ├── transport/      # 传输层协议 (处理 WS 握手、TCP 字节流)
│   │   ├── websocket/
│   │   └── tcp/
│   ├── packet/         # 数据包定义 (处理自定义包头: MsgId, Seq, Length)
│   └── payload/        # 业务载荷序列化
│       └── protobuf/
├── network/            # 网络核心驱动
│   ├── engine.go       # gnet 服务生命周期管理
│   ├── connection.go   # 封装 gnet.Conn，增强 Agent 绑定能力
│   ├── reactor.go      # IO 事件监听与处理
│   └── dispatcher.go   # 消息初筛与任务分发逻辑
├── router/             # 业务路由系统
│   ├── router.go       # 逻辑分发核心 (Handler 映射)
│   └── group.go        # 广播组/房间管理组件
├── Agent/            # 会话管理
│   └── Agent.go      # Agent 上下文
├── runtime/            # 运行时基础设施
│   ├── pool/           # 高性能协程池 (Ants 或自定义)
│   └── buffer/         # 内存池与缓冲区优化
├── options.go          # 函数式配置选项
├── gneto.go            # 库统一入口 (NewServer/Client)
└── go.mod
```

## 核心模块职责

### 编解码三段论 (codec/)

*   **transport**: (传输层)       决定了数据“怎么过来”  (是 WS 帧还是 TCP 流)
*   **packet**:    (分包层/协议头) 决定了包结构“长什么样”。负责解析出关键的MsgId, 以及后续可扩展的 魔树(快速过滤非法连接),序列号,加密,签名压缩等功能
*   **payload**:   (负载/业务对象) 决定了packet中的Body数据“如何序列化”。(Protobuf/JSON ...)
* 编码顺序 从左到右 -->  |  解码顺序 从右到左 <--
  * payload - packet - transport

### 网络驱动与分发 (network/)

*   **reactor.go**: 直接对接 gnet 的回调，负责最快速的 IO 响应。
*   **dispatcher.go**: 负责将 reactor 抓取到的完整 Packet 投递到 runtime/pool。它可以根据 PlayerId 进行 Hash 分发，确保同一个玩家的消息在同一个 Goroutine 里顺序执行。

### 会话系统 (Agent/)

游戏业务中，Agent 是连接（Conn）与玩家（Player）之间的纽带。独立成包后，可以更方便地处理“玩家掉线、Conn 断开但 Agent 保留 30 秒”这种典型的断线重连需求。

## 游戏场景深度设计

### 逻辑单线程化 (并发模型)

为简化业务开发并避免在业务代码中引入大量锁，本框架在 dispatcher.go 中实现了基于 Agent 的逻辑单线程化调度模型：

调度流程
*   Agent 绑定: 每个连接在建立时绑定唯一的 AgentId。
*   一致性分发: 在 IO 线程中，从连接中获取 AgentId，通过取模或哈希算法映射到固定的逻辑队列索引。。
*   业务消息被投递到对应索引的协程池（或逻辑队列）中执行。

并发保证
*   同一 Agent 的所有逻辑处理始终在同一个逻辑执行单元中串行执行
*   不同 Agent 之间可并行处理，充分利用多核资源
*   业务 Handler 无需显式加锁，也无需关心并发竞争问题

### 零拷贝与广播优化

在 `runtime/buffer` 中引入内存池：

*   **序列化合并**: 广播消息时，在 payload 层只序列化一次。
*   **引用计数**: 在多连接发送时，通过引用计数管理这一块内存，直到所有 AsyncWrite 完成后再回收到池中。
