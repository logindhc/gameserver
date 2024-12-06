package player

import (
	cstring "gameserver/cherry/extend/string"
	clog "gameserver/cherry/logger"
	"gameserver/cherry/net/parser/pomelo"
	cproto "gameserver/cherry/net/proto"
	"gameserver/internal/code"
	"gameserver/internal/event"
	"gameserver/internal/pb"
	sessionKey "gameserver/internal/session_key"
	"gameserver/nodes/game/db"
	"gameserver/nodes/game/module/online"
)

type (
	// ActorPlayer 每位登录的玩家对应一个子actor
	ActorPlayer struct {
		pomelo.ActorBase
		IsOnline bool // 玩家是否在线
		Id       int64
		Item
	}
)

func (p *ActorPlayer) OnInit() {
	clog.Debugf("[ActorPlayer] path = %s init!", p.PathString())
	p.Item = Item{p}
	p.Item.OnInit()
	// 注册 session关闭的remote函数(网关触发连接断开后，会调用RPC发送该消息)
	p.Remote().Register("sessionClose", p.sessionClose)
	p.Local().Register("enter", p.playerEnter) // 注册 进入角色

}

func (p *ActorPlayer) OnStop() {
	clog.Debugf("[ActorPlayer] path = %s exit!", p.PathString())
}

// sessionClose 接收角色session关闭处理
func (p *ActorPlayer) sessionClose() {
	online.UnBindPlayer(p.Id)
	p.IsOnline = false

	logoutEvent := event.NewPlayerLogout(p.ActorID(), p.Id)
	p.PostEvent(&logoutEvent)
	p.Exit()
	clog.Debugf("[ActorPlayer] exit! id = %d", p.Id)
}

// PlayerEnter 玩家进入游戏
func (p *ActorPlayer) playerEnter(session *cproto.Session, req *pb.Int64) {
	playerId := req.Value
	if playerId < 1 {
		p.ResponseCode(session, code.PlayerIdError)
		return
	}

	// 检查并查找该用户下的该角色
	playerTable := db.PlayerRepository.Get(playerId)
	if playerTable == nil {
		// 创建角色
		playerTable = p.playerCreate(session)
	}

	if playerTable == nil {
		p.ResponseCode(session, code.PlayerIdError)
		return
	}

	// 保存进入游戏的玩家对应的agentPath
	online.BindPlayer(playerId, session.AgentPath)

	// 设置网关节点session的PlayerID属性
	p.Call(session.ActorPath(), "setSession", &pb.StringKeyValue{
		Key:   sessionKey.PlayerID,
		Value: cstring.ToString(playerId),
	})

	p.Id = playerId
	p.IsOnline = true // 设置为在线状态

	// 这里改为客户端主动请求更佳
	// [01]推送角色 道具数据
	//module.Item.ListPush(session, playerId)
	// [02]推送角色 英雄数据
	//module.Hero.ListPush(session, playerId)
	// [03]推送角色 成就数据
	//module.Achieve.CheckNewAndPush(playerId, true, true)
	// [04]推送角色 邮件数据
	//module.Mail.ListPush(session, playerId)

	// [99]最后推送 角色进入游戏响应结果
	response := &pb.PlayerEnterResponse{}
	response.GuideMaps = map[int32]int32{}
	response.Player = buildPBPlayer(playerTable)
	p.Response(session, response)

	// 角色登录事件
	loginEvent := event.NewPlayerLogin(p.ActorID(), playerId)
	p.PostEvent(&loginEvent)
}

// playerCreate 玩家创角
func (p *ActorPlayer) playerCreate(session *cproto.Session) *db.PlayerTable {
	// 获取创角初始化配置
	// 创建角色&添加角色初始的资产
	newPlayerTable, errCode := db.CreatePlayer(session)
	if code.IsFail(errCode) {
		p.ResponseCode(session, errCode)
		return nil
	}
	// TODO 更新最后一次登陆的角色信息到中心节点

	// 抛出角色创建事件
	playerCreateEvent := event.NewPlayerCreate(newPlayerTable.ID, newPlayerTable.Nickname, newPlayerTable.Gender)
	p.PostEvent(&playerCreateEvent)
	return newPlayerTable
}

func buildPBPlayer(playerTable *db.PlayerTable) *pb.Player {
	return &pb.Player{
		PlayerId:   playerTable.ID,
		PlayerName: playerTable.Nickname,
		Level:      playerTable.Level,
		CreateTime: playerTable.CreateTime,
		Exp:        playerTable.Exp,
		Gender:     playerTable.Gender,
	}
}
