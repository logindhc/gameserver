package controller

import (
	cherryGin "gameserver/cherry/components/gin"
	cstring "gameserver/cherry/extend/string"
	cherryLogger "gameserver/cherry/logger"
	"gameserver/internal/code"
	"gameserver/internal/data"
	sessionKey "gameserver/internal/session_key"
	"gameserver/internal/token"
	"gameserver/nodes/web/cache"
	"gameserver/nodes/web/sdk"
)

type Controller struct {
	cherryGin.BaseController
}

func (p *Controller) Init() {
	group := p.Group("/")
	group.GET("/update", p.updateServer)
	group.GET("/serverInfo/:channel", p.serverInfo)
}

// 后台更新服务器状态
// http://127.0.0.1/update?server_id=10001&status=1&server_host=3&server_port=123&sign=123&time=123
func (p *Controller) updateServer(c *cherryGin.Context) {
	serverId := c.GetInt("server_id", 0, true)
	if serverId < 1 {
		cherryLogger.Warnf("if serverId < 1 . params=%s", c.GetParams())
		code.RenderResult(c, code.Error)
		return
	}
	status := c.GetInt("status", 0, true)
	if status < 1 {
		cherryLogger.Warnf("if status < 1 . params=%s", c.GetParams())
		code.RenderResult(c, code.Error)
		return
	}
	serverHost := c.GetString("server_host", "", true)
	if cstring.IsBlank(serverHost) {
		cherryLogger.Warnf("if serverHost is blank . params=%s", c.GetParams())
		code.RenderResult(c, code.Error)
		return
	}
	serverPort := c.GetString("server_port", "", true)
	if cstring.IsBlank(serverPort) {
		cherryLogger.Warnf("if serverPort is blank . params=%s", c.GetParams())
		code.RenderResult(c, code.Error)
		return
	}
	server := &cache.Server{
		ServerId:   serverId,
		Status:     status,
		ServerHost: serverHost,
		ServerPort: serverPort,
	}
	err := cache.UpdateServerStatus(server)
	if err != nil {
		cherryLogger.Warnf("update server status error. server=%v, error=%s", server, err)
		code.RenderResult(c, code.Error)
		return
	}

	code.RenderResult(c, code.OK)
}

// login 根据channel获取sdkConfig，与第三方进行帐号登陆效验,验证完毕后，返回token和连接地址
// http://127.0.0.1/serverInfo?channel=101&code=test11&platform=3&time=123&sign=123
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

		uidStr, found := result.GetString(sessionKey.PlayerID)
		if found == false {
			cherryLogger.Warnf("callback result map not found `uid`. result = %s", result)
			code.RenderResult(c, code.LoginError)
			return
		}
		uid := cstring.ToInt64D(uidStr)
		serverId := result.GetInt(sessionKey.ServerID)
		// get server status
		server, err2 := cache.GetServerStatus(serverId)
		if err2 != nil {
			cherryLogger.Warnf("get server status error. serverId=%d, error=%s", serverId, err2)
			code.RenderResult(c, code.ServerError)
			return
		}
		base64Token := token.New(uid, openId, channel, platform, int32(serverId), config.Salt).ToBase64()
		res := map[string]interface{}{
			"token":  base64Token,
			"server": server,
		}
		code.RenderResult(c, code.OK, res)
	})
}
