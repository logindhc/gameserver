package code

var (
	OK                       int32 = 0   // is ok
	Error                    int32 = 1   // error
	ParamError               int32 = 2   // 参数错误
	ConfigError              int32 = 99  // 配置错误
	ChannelIDError           int32 = 100 // channel错误
	ServerError              int32 = 101 // 服务器异常
	VersionError             int32 = 102 // version异常
	PlatformIDError          int32 = 103 // channel错误
	SDKError                 int32 = 201 // sdk验证异常
	AccountAuthFail          int32 = 202 // 帐号授权失败
	AccountBindFail          int32 = 203 // 帐号绑定失败
	AccountTokenValidateFail int32 = 204 // token验证失败
	AccountNameIsExist       int32 = 205 // 帐号已存在
	AccountRegisterError     int32 = 206 //
	AccountGetFail           int32 = 207 //
	PlayerDenyLogin          int32 = 301 // 玩家禁止登录
	PlayerDuplicateLogin     int32 = 302 // 玩家重复登录
	PlayerNameExist          int32 = 303 // 玩家角色名已存在
	PlayerCreateFail         int32 = 304 // 玩家创建角色失败
	PlayerNotLogin           int32 = 305 // 玩家未登录
	PlayerIdError            int32 = 306 // 玩家id错误
	GoldNotEnough            int32 = 401 // 金币不足
	MoneyNotEnough           int32 = 402 // 银币不足
	DiamondNotEnough         int32 = 403 // 钻石不足
	ItemNotEnough            int32 = 501 // 道具不足
	ItemNotAvailable         int32 = 502 // 道具不可用
	HeroNotEnough            int32 = 601 // 英雄不存在
	HeroLevelError           int32 = 602 // 英雄等级异常
	HeroMaxLevel             int32 = 603 // 英雄已经是最大等级
)
