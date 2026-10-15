package coden

type Code = int32

//goland:noinspection GoUnusedConst
const (
	Ok Code = 0 // 成功
)

// HTTP：响应状态码
//
//goland:noinspection GoUnusedConst
const (
	HttpOk                  Code = 200 // 请求成功
	HttpBadRequest          Code = 400 // 错误请求
	HttpUnauthorized        Code = 401 // 未授权
	HttpForbidden           Code = 403 // 禁止访问
	HttpNotFound            Code = 404 // 资源不存在
	HttpMethodNotAllowed    Code = 405 // 方法不允许
	HttpTimeout             Code = 408 // 请求超时
	HttpTooManyRequests     Code = 429 // 请求过多
	HttpInternalServerError Code = 500 // 服务器内部错误
	HttpServiceUnavailable  Code = 503 // 服务不可用

	HttpRequestFailed Code = 599 // http请求异常
)

// Common：通用错误
//
//goland:noinspection GoUnusedConst
const (
	CommonUnknown   Code = 1000 // 未知错误
	CommonInternal  Code = 1001 // 内部错误
	CommonParam     Code = 1002 // 参数错误
	CommonRateLimit Code = 1003 // 限流
)

// Player：玩家相关
//
//goland:noinspection GoUnusedConst
const (
	PlayerNotFound           Code = 1200 // 玩家不存在
	PlayerExist              Code = 1201 // 玩家已存在
	PlayerTokenInvalid       Code = 1202 // 玩家 token 无效
	PlayerTokenCreateFailed  Code = 1203 // 玩家 token 创建失败
	PlayerLockOther          Code = 1204 // 玩家锁在其它游戏中未退出
	PlayerBalanceNotEnough   Code = 1205 // 余额不足
	PlayerBanned             Code = 1206 // 被封禁
	PlayerNumMaxLimit        Code = 1207 // 玩家注册数量达到上限
	PlayerCurrencyNotSupport Code = 1208 // 不支持的货币类型
	PlayerCreateDeny         Code = 1209 // 玩家拒绝创建
)

// Transaction：交易相关
//
//goland:noinspection GoUnusedConst
const (
	TransactionFailed        Code = 1300 // 交易失败
	TransactionCanceled      Code = 1301 // 交易已经撤销
	TransactionDeny          Code = 1302 // 不允许投注，玩家不允许特定游戏
	TransactionLimit         Code = 1303 // 达到投注限额
	TransactionStatusUnknown Code = 1304 // 交易异常，状态未知
	TransactionTimeout       Code = 1305 //运营商交易超时
	TransactionClose         Code = 1306 // 交易主动撤销 已知交易成功
)

// Operator：运营商相关
//
//goland:noinspection GoUnusedConst
const (
	OperatorNotExist Code = 1400 // 运营商不存在
	SignInvalid      Code = 1401 // 运营商签名无效
)

// Server：服务器相关
//
//goland:noinspection GoUnusedConst
const (
	ServerAllClosed  Code = 1500 // 全局停服
	ServerGameClosed Code = 1501 // 游戏停服
)

// Db：数据库相关
//
//goland:noinspection GoUnusedConst
const (
	DbError         Code = 1600 // 数据库错误 // todo 目前表示所有，子项的子项同理
	DbRedisError    Code = 1610 // redis数据库错误
	DbMysqlError    Code = 1620 // mysql数据库错误
	DbMongoError    Code = 1630 // mongo数据库错误
	DbSlsError      Code = 1640 // sls数据库错误
	DbSlsLockFailed Code = 1641 // 获取sls共享锁失败, 稍后再试
)

// OperatorRequestFailed 运营商api相关
//
//goland:noinspection GoUnusedConst
const (
	OperatorRequestErr            Code = 1700 // 请求未发出
	OperatorHttpRequestErr        Code = 1701 // 请求发出失败
	OperatorResponseBodyParseErr  Code = 1702 // 响应报文解析失败
	OperatorHttpTimeoutErr        Code = 1703 // 响应超时
	OperatorHttpServiceFailed     Code = 1704 // 运营商响应业务码为失败
	OperatorHttpServiceUnknown    Code = 1705 // 运营商响应码未知
	OperatorHttpServiceNotSuccess Code = 1706 // 运营商响应码非成功状态
	OperatorHttpNotSuccess        Code = 1707 // 运营商http响应码非200
	OperatorServerClose           Code = 1708 // 运营商停服
)
