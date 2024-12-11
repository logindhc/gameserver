package controller

import (
	cherryGin "gameserver/cherry/components/gin"
	cstring "gameserver/cherry/extend/string"
	cherryLogger "gameserver/cherry/logger"
	"gameserver/internal/code"
	"gameserver/nodes/web/cache"
)

type GMController struct {
	cherryGin.BaseController
}

func (gc *GMController) Init() {
	group := gc.Group("/")
	group.GET("/gm/server/update", gc.updateServer)
	group.GET("/gm/ipwhite/update", gc.updateIpWhite)
	group.GET("/gm/ipwhite/del", gc.delIpWhite)
	group.GET("/gm/channel/update", gc.updateChannel)
	group.GET("/gm/channel/del", gc.delChannel)
	group.GET("/gm/appcheck/update", gc.updateAppCheck)
	group.GET("/gm/appcheck/del", gc.delAppCheck)
}

// 后台更新服务器状态
// http://127.0.0.1//gm/server/update?server_id=10001&status=1&server_host=3&server_port=123&sign=123&time=123
func (gc *GMController) updateServer(c *cherryGin.Context) {
	serverId := c.GetInt32("server_id", 0, true)
	if serverId < 1 {
		cherryLogger.Warnf("if serverId < 1 . params=%s", c.GetParams())
		code.RenderResult(c, code.Error)
		return
	}
	status := c.GetInt32("status", 0, true)
	if status < 0 {
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
	checkHost := c.GetString("check_host", "", true)
	checkPort := c.GetString("check_port", "", true)

	server := &cache.Server{
		ServerId:   serverId,
		Status:     status,
		ServerHost: serverHost,
		ServerPort: serverPort,
		CheckHost:  checkHost,
		CheckPort:  checkPort,
	}
	err := cache.UpdateServerInfo(server)
	if err != nil {
		cherryLogger.Warnf("update server status error. server=%v, error=%s", server, err)
		code.RenderResult(c, code.Error)
		return
	}

	code.RenderResult(c, code.OK)
}

func (gc *GMController) updateIpWhite(c *cherryGin.Context) {
	ip := c.GetInt32("ip", 0, true)
	if ip <= 0 {
		cherryLogger.Warnf("if ip <=0 . params=%s", c.GetParams())
		code.RenderResult(c, code.Error)
		return
	}
	ipWhite := new(cache.IpWhite)
	ipWhite.Ip = ip
	remark := c.GetString("remark", "", true)
	if cstring.IsNotBlank(remark) {
		ipWhite.Remark = remark
	}
	ipWhite.Enable = c.GetInt("enable", 0, true)
	err := cache.UpdateIpWhite(ipWhite)
	if err != nil {
		cherryLogger.Warnf("update ip white error. ip=%v, error=%s", ip, err)
		code.RenderResult(c, code.Error)
		return
	}
	code.RenderResult(c, code.OK)
}

func (gc *GMController) delIpWhite(c *cherryGin.Context) {
	ip := c.GetInt32("ip", 0, true)
	if ip <= 0 {
		cherryLogger.Warnf("if ip <=0 . params=%s", c.GetParams())
		code.RenderResult(c, code.Error)
		return
	}
	err := cache.DelIpWhite(ip)
	if err != nil {
		cherryLogger.Warnf("del ip white error. ip=%v, error=%s", ip, err)
		code.RenderResult(c, code.Error)
		return
	}
	code.RenderResult(c, code.OK)
}

func (gc *GMController) updateChannel(c *cherryGin.Context) {
	channelId := c.GetInt32("id", 0, true)
	if channelId <= 0 {
		cherryLogger.Warnf("if channelId <=0 . params=%s", c.GetParams())
		code.RenderResult(c, code.Error)
		return
	}
	serverIds := c.GetString("server_ids", "", true)
	if cstring.IsBlank(serverIds) {
		cherryLogger.Warnf("if serverIds is blank . params=%s", c.GetParams())
		code.RenderResult(c, code.Error)
		return
	}
	serverIdMap := cstring.SplitMapString(serverIds, ",")
	channelInfo := new(cache.ChannelInfo)
	channelInfo.Id = channelId
	channelInfo.ServerIds = serverIdMap
	name := c.GetString("name", "", true)
	if cstring.IsNotBlank(name) {
		channelInfo.Name = name
	}
	err := cache.UpdateChannelInfo(channelInfo)
	if err != nil {
		cherryLogger.Warnf("update channel error. channel=%v, error=%s", channelInfo, err)
		code.RenderResult(c, code.Error)
		return
	}
	code.RenderResult(c, code.OK)
}

func (gc *GMController) delChannel(c *cherryGin.Context) {
	id := c.GetInt32("id", 0, true)
	if id <= 0 {
		cherryLogger.Warnf("if id <=0 . params=%s", c.GetParams())
		code.RenderResult(c, code.Error)
		return
	}
	err := cache.DelChannelInfo(id)
	if err != nil {
		cherryLogger.Warnf("del channel error. id=%v, error=%s", id, err)
		code.RenderResult(c, code.Error)
		return
	}
	code.RenderResult(c, code.OK)
}

func (gc *GMController) updateAppCheck(c *cherryGin.Context) {
	channelId := c.GetInt32("id", 0, true)
	if channelId <= 0 {
		cherryLogger.Warnf("if channelId <=0 . params=%s", c.GetParams())
		code.RenderResult(c, code.Error)
		return
	}
	version := c.GetString("version", "", true)
	tipVersion := c.GetString("tipVersion", "", true)
	forceVersion := c.GetString("forceVersion", "", true)
	err := cache.UpdateAppCheckInfo(channelId, version, tipVersion, forceVersion)
	if err != nil {
		cherryLogger.Warnf("update appcheck error. channelId=%v, error=%s", channelId, err)
		code.RenderResult(c, code.Error)
		return
	}
	code.RenderResult(c, code.OK)
}

func (gc *GMController) delAppCheck(c *cherryGin.Context) {
	id := c.GetInt32("id", 0, true)
	if id <= 0 {
		cherryLogger.Warnf("if id <=0 . params=%s", c.GetParams())
		code.RenderResult(c, code.Error)
		return
	}
	err := cache.DelAppCheckInfo(id)
	if err != nil {
		cherryLogger.Warnf("del appcheck error. id=%v, error=%s", id, err)
		code.RenderResult(c, code.Error)
		return
	}
	code.RenderResult(c, code.OK)
}
