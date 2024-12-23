package player

import (
	cstring "gameserver/cherry/extend/string"
	ctime "gameserver/cherry/extend/time"
	cfacade "gameserver/cherry/facade"
	clog "gameserver/cherry/logger"
	"gameserver/cherry/net/parser/pomelo"
	cproto "gameserver/cherry/net/proto"
	"gameserver/internal/cache"
	"gameserver/internal/code"
	"gameserver/internal/constant"
	"gameserver/internal/data"
	"gameserver/internal/event"
	"gameserver/internal/pb"
	sessionKey "gameserver/internal/session_key"
	"gameserver/internal/utils"
	"gameserver/nodes/game/db"
	"gameserver/nodes/game/module/online"
	"gameserver/nodes/game/opcode"
	"gameserver/nodes/game/res/resmgr"
)

type (
	// ActorPlayer 每位登录的玩家对应一个子actor
	ActorPlayer struct {
		pomelo.ActorBase
		IsOnline bool  // 玩家是否在线
		Id       int64 //enter之后才有值,只是用来处理sessionClose来解绑
		currency *ActorCurrency
		item     *ActorItem
		hero     *ActorHero
	}
)

func (p *ActorPlayer) OnInit() {
	clog.Debugf("[ActorPlayer] path = %s init!", p.PathString())
	// 注册 session关闭的remote函数(网关触发连接断开后，会调用RPC发送该消息)
	p.Remote().Register("sessionClose", p.sessionClose)
	p.Local().Register("enter", p.playerEnter) // 注册 进入角色
	p.Local().Register("gm", p.gm)
	//货币模块
	p.Local().Register("currencyInfo", p.currency.getInfo)
	//道具模块
	p.Local().Register("itemInfo", p.item.getInfo)
	p.Local().Register("itemUse", p.item.use)
	//英雄模块
	p.Local().Register("heroInfo", p.hero.getInfo)
	p.Local().Register("heroUp", p.hero.up)
}

func (p *ActorPlayer) OnStop() {
	clog.Debugf("[ActorPlayer] path = %s exit!", p.PathString())
}

// sessionClose 接收角色session关闭处理
func (p *ActorPlayer) sessionClose() {
	if p.Id < 1 {
		//处理游戏服重启的情况
		clog.Debug("[ActorPlayer] session close, but player not exist!")
		p.Exit()
		return
	}
	online.UnBindPlayer(p.Id)
	p.IsOnline = false
	logoutEvent := event.NewPlayerLogout(p.ActorID(), p.Id)
	p.PostEvent(&logoutEvent)
	p.Exit()
	clog.Debugf("[ActorPlayer] exit! id = %d", p.Id)
}

// PlayerEnter 玩家进入游戏
func (p *ActorPlayer) playerEnter(session *cproto.Session, _ *pb.C2SPlayerEnter) {
	playerId := session.Uid
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
	serverId := session.GetInt32(sessionKey.ServerID)
	server, err := cache.GetServerInfo(serverId)
	if err != nil {
		return
	}

	if server.Status == 0 && playerTable.White == 0 {
		//维护中，不是白名单不允许进游戏
		p.ResponseCode(session, code.PlayerDenyLogin)
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
	p.currency.pushInfo(session)
	p.item.pushInfo(session)
	p.hero.pushInfo(session)
	// [99]最后推送 角色进入游戏响应结果
	response := &pb.S2CPlayerEnter{}
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
	gameRow, b := data.GameConfig.Get(constant.PlayerInitResID)
	if b {
		_, err := resmgr.Instance.Add(newPlayerTable.ID, gameRow.Arr2Val, opcode.PlayerInit, "player_create")
		if err != nil {
			p.ResponseCode(session, code.PlayerCreateFail)
			clog.Errorf("[playerCreate] add diamond fail! err = %v", err)
			return nil
		}
	}
	//db.LogRegisterRepository.Add(&db.LogRegister{
	//	Device:   newPlayerTable.OpenId,
	//	Channel:  newPlayerTable.Channel,
	//	Platform: newPlayerTable.Platform,
	//	Time:     newPlayerTable.CreateTime,
	//})
	// 抛出角色创建事件
	playerCreateEvent := event.NewPlayerCreate(newPlayerTable.ID, newPlayerTable.Nickname, newPlayerTable.Gender)
	p.PostEvent(&playerCreateEvent)
	return newPlayerTable
}

func buildPBPlayer(playerTable *db.PlayerTable) *pb.Player {
	return &pb.Player{
		PlayerId:   playerTable.ID,
		PlayerName: playerTable.Nickname,
		CreateTime: playerTable.CreateTime,
		Gender:     playerTable.Gender,
	}
}
func (p *ActorPlayer) gm(session *cproto.Session, req *pb.C2SPlayerGM) {
	player := db.PlayerRepository.GetOrCreate(session.Uid)
	if player.GmLevel == 0 {
		//return
	}
	cmd := req.Cmd
	args := req.Args
	switch cmd {
	case "addres":
		array, ok := cstring.Split2ArrayInt(args, ",")
		if !ok || len(array) != 3 {
			return
		}
		addOne, err := resmgr.Instance.AddOne(player.ID, array[0], array[1], array[2], opcode.Gm, "gm")
		if err != nil {
			clog.Warnf("[gm] add res fail! err = %v", err)
			p.Response(session, &pb.S2CPlayerGM{
				Result: false,
			})
			return
		}
		p.PushResUpdateInfo(session, addOne)
	}
	p.Response(session, &pb.S2CPlayerGM{
		Result: true,
	})
}

// onLoginEvent 玩家登陆事件处理
func (p *ActorPlayer) onLoginEvent(e cfacade.IEventData) {
	evt, ok := e.(*event.PlayerLogin)
	if ok == false {
		return
	}
	cherryTime := ctime.Now()
	second := cherryTime.ToSecond()

	player := db.PlayerRepository.GetOrCreate(evt.PlayerId)
	player.LastLoginTime = second
	db.PlayerRepository.Update(player)

	dotLogin := db.LogLogin{
		ID:         evt.PlayerId,
		FirstTime:  &second,
		LastTime:   &second,
		DayIndex:   cherryTime.ToShortInt32DateFormat(),
		TotalCount: 1,
	}
	db.LogLoginRepository.Add(&dotLogin)

	clog.Infof("[PlayerLoginEvent] [playerId = %d, onlineCount = %d]",
		evt.PlayerId,
		online.Count(),
	)
}

// onLoginEvent 玩家登出事件处理
func (p *ActorPlayer) onLogoutEvent(e cfacade.IEventData) {
	evt, ok := e.(*event.PlayerLogout)
	if !ok {
		return
	}
	if evt.PlayerId <= 0 {
		return
	}
	cherryTime := ctime.Now()
	second := cherryTime.ToSecond()
	player := db.PlayerRepository.GetOrCreate(evt.PlayerId)
	player.LastLogoutTime = second
	db.PlayerRepository.Update(player)
	clog.Infof("[PlayerLogoutEvent] [playerId = %d, onlineCount = %d]",
		evt.PlayerId,
		online.Count(),
	)
}

// onPlayerCreateEvent 玩家创建事件
func (p *ActorPlayer) onPlayerCreateEvent(e cfacade.IEventData) {
	evt, ok := e.(*event.PlayerCreate)
	if !ok {
		return
	}
	clog.Infof("[onPlayerCreateEvent] [%+v]", evt)
}

func (p *ActorPlayer) PushResUpdateInfo(session *cproto.Session, ress [][]int) {
	s2CResUpdate := &pb.S2CResUpdate{ResInfos: make([]*pb.ResInfo, 0, len(ress))}
	s2CResUpdate.ResInfos = append(s2CResUpdate.ResInfos, utils.ConverterArray2ResInfo(ress)...)
	p.Push(session, "resUpdate", s2CResUpdate)
}
