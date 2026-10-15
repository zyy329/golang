package gn

// 事件广播
//
//goland:noinspection GoUnusedConst
const (
	RedisChannelBasicConfigServer = "basic:config:server" // 全服状态
	RedisChannelBasicOperator     = "basic:operator"      // 运营商【operatorId】
	RedisChannelBasicCurrency     = "basic:currency"      // 币种【currency】
	RedisChannelBasicGame         = "basic:game"          // 游戏【gameId,gameId】

	RedisChannelGameConfig  = "game:conf"     // 游戏:配置【table】【table:key】(key表示约定filter字段)
	RedisChannelGameControl = "game:cont"     // 游戏:控制【table】【table:key】
	RedisChannelGameKick    = "game:kick"     // 游戏:T人
	RedisChannelGameRtp     = "game:cont:rtp" // 游戏:RTP更新，格式："9527:0.95"

	RedisChannelWebBaseDomain  = "web:base:domain"  // 基础域名
	RedisChannelWebSystemEvent = "web:system:event" // 系统事件【id,id】

	RedisChannelLobbySortDomain = "lobby:base:domain" // 基础域名权重变更事件

	RedisChannelWebClientLogSwitch        = "web:client:log-witch"          // 游戏:客户端:日志开关
	RedisChannelWebClientBundle           = "web:client:bundle"             // 游戏:客户端:bundle包【code,code】
	RedisChannelWebClientForced           = "web:client:forced"             // 游戏:客户端:游戏强更【code,code】
	RedisChannelWebClientSkinOperator     = "web:client:skin:operator"      // 游戏:客户端:皮肤:运营商【json】
	RedisChannelWebClientSkinOperatorGame = "web:client:skin:operator:game" // 游戏:客户端:皮肤:运营商:游戏【json】

	RedisChannelWebDocAccess    = "web:doc:access"    // 技术文档访问码【id】
	RedisChannelWebOperatorUser = "web:operator:user" // 运营商后台用户【id】
)
