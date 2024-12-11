package controller

import (
	cherryGin "gameserver/cherry/components/gin"
	cstring "gameserver/cherry/extend/string"
	cherryLogger "gameserver/cherry/logger"
	"gameserver/internal/code"
	"gameserver/internal/data"
	sessionKey "gameserver/internal/session_key"
	"gameserver/internal/token"
	"gameserver/internal/utils"
	"gameserver/nodes/web/cache"
	"gameserver/nodes/web/sdk"
)

type Controller struct {
	cherryGin.BaseController
}

func (p *Controller) Init() {
	cache.InitServer()
	group := p.Group("/")
	group.GET("/api/serverInfo/:channel", p.serverInfo)
}

// login 根据channel获取sdkConfig，与第三方进行帐号登陆效验,验证完毕后，返回token和连接地址
// http://127.0.0.1/serverInfo?channel=101&code=test11&platform=3&time=123&sign=123&version=1.1
func (p *Controller) serverInfo(c *cherryGin.Context) {
	channel := c.GetInt32("channel", 0, true)
	if channel < 1 {
		cherryLogger.Warnf("if channel < 1 . params=%s", c.GetParams())
		code.RenderResult(c, code.ChannelIDError)
		return
	}
	platform := c.GetInt32("platform", 0, true)
	if platform < 1 {
		cherryLogger.Warnf("if platform < 1 . params=%s", c.GetParams())
		code.RenderResult(c, code.ChannelIDError)
		return
	}

	config := data.SdkConfig.Get(channel)
	if config == nil {
		cherryLogger.Warnf("if platformConfig == nil . params=%s", c.GetParams())
		code.RenderResult(c, code.LoginError)
		return
	}

	version := c.GetString("version", "", true)
	if cstring.IsBlank(version) {
		cherryLogger.Warnf("if version is blank . params=%s", c.GetParams())
		code.RenderResult(c, code.VersionError)
		return
	}

	sdkInvoke, err := sdk.GetInvoke(config.SdkId)
	if err != nil {
		cherryLogger.Warnf("[channel = %d] get invoke error. params=%s", channel, c.GetParams())
		code.RenderResult(c, code.ChannelIDError)
		return
	}

	params := c.GetParams(true)
	// invoke login
	sdkInvoke.Login(config, params, func(statusCode int32, result sdk.Params, error ...error) {
		if code.IsFail(statusCode) {
			cherryLogger.Warnf("login validate fail. code = %d, params = %s", statusCode, c.GetParams())
			if len(error) > 0 {
				cherryLogger.Warnf("code = %d, error = %s", statusCode, error[0])
			}
			code.RenderResult(c, statusCode)
			return
		}

		if result == nil {
			cherryLogger.Warnf("callback result map is nil. params= %s", c.GetParams())
			code.RenderResult(c, code.LoginError)
			return
		}

		openId, found := result.GetString(sessionKey.OpenID)
		if found == false {
			cherryLogger.Warnf("callback result map not found `openId`. result = %s", result)
			code.RenderResult(c, code.LoginError)
			return
		}

		channelInfo, err := cache.GetChannelInfo(channel)
		if err != nil {
			cherryLogger.Warnf("get channel info error. channel=%d, error=%s", channel, err)
			code.RenderResult(c, code.ServerError)
			return
		}
		isCheck := false
		serverId := int32(0)
		if cstring.IsNotBlank(channelInfo.Version) && version == channelInfo.Version {
			//提审版本
			isCheck = true
			serverId = channelInfo.ServerIds[0] //提服直接获取第一个serverId，所有相同渠道的提审服应该是一样的
		} else {
			//根据最小负载的game节点
			nodeIds, ok := utils.GetAllNodeIdByRank()
			if ok != nil {
				cherryLogger.Warnf("get game node id error. error=%s", ok)
				code.RenderResult(c, code.ServerError)
				return
			}
			for _, sId := range nodeIds {
				for _, cId := range channelInfo.ServerIds { //渠道指定serverId，获取最小的负载
					iId := cstring.ToInt32D(sId)
					if iId == cId {
						serverId = iId
						break
					}
				}
			}
			if serverId == 0 { //当前渠道拿不到服id，后台配置有问题
				cherryLogger.Warnf("get serverId error. channel=%d, nodeIds=%s,channelInfoServerIds=%v", channel, nodeIds, channelInfo.ServerIds)
				code.RenderResult(c, code.ServerError)
				return
			}
		}

		server, err2 := cache.GetServerInfo(serverId)
		if err2 != nil {
			cherryLogger.Warnf("get server status error. serverId=%d, error=%s", serverId, err2)
			code.RenderResult(c, code.ServerError)
			return
		}
		res := map[string]any{}
		base64Token := token.New(openId, channel, platform, config.Salt).ToBase64()
		if isCheck {
			//提审版本直接进提审服
			res["server_host"] = server.CheckHost
			res["server_port"] = server.CheckPort
		} else {
			//正常进服
			res["server_host"] = server.ServerHost
			res["server_port"] = server.ServerPort
			if cstring.IsNotBlank(channelInfo.TipVersion) {
				res["tip_version"] = channelInfo.TipVersion
			}
			if cstring.IsNotBlank(channelInfo.ForceVersion) {
				res["force_version"] = channelInfo.ForceVersion
			}
		}
		res["token"] = base64Token
		res["is_check"] = isCheck
		res["server_state"] = server.Status
		ip := c.ClientIP()
		res["ip"] = ip
		ipWhite, err3 := cache.GetIpWhite(utils.IP2Long(ip))
		if err3 != nil {
			res["isWhite"] = false
		} else {
			res["isWhite"] = ipWhite.Enable == 1
		}
		code.RenderResult(c, code.OK, res)
	})
}
