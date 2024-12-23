package player

import (
	"gameserver/nodes/game/module/online"
	"time"

	cfacade "gameserver/cherry/facade"
	clog "gameserver/cherry/logger"
	"gameserver/cherry/net/parser/pomelo"
)

type (
	// ActorPlayers 玩家总管理actor
	ActorPlayers struct {
		pomelo.ActorBase
		childExitTime time.Duration
	}
)

func (p *ActorPlayers) AliasID() string {
	return "player"
}

func (p *ActorPlayers) OnInit() {
	p.childExitTime = time.Minute * 30
	// 注册角色登陆事件
	//p.Event().Register(event.PlayerLoginKey, p.onLoginEvent)
	//p.Event().Register(event.PlayerLogoutKey, p.onLogoutEvent)
	//p.Event().Register(event.PlayerCreateKey, p.onPlayerCreateEvent)
}

func (p *ActorPlayers) OnFindChild(msg *cfacade.Message) (cfacade.IActor, bool) {
	// 动态创建 player child actor
	childID := msg.TargetPath().ChildID
	actorPlayer := &ActorPlayer{IsOnline: false}
	childActor, err := p.Child().Create(childID, actorPlayer)
	actorPlayer.item = NewActorItem(actorPlayer)
	actorPlayer.hero = NewActorHero(actorPlayer)
	actorPlayer.currency = NewActorCurrency(actorPlayer)
	if err != nil {
		return nil, false
	}

	return childActor, true
}

//
//// onLoginEvent 玩家登陆事件处理
//func (p *ActorPlayers) onLoginEvent(e cfacade.IEventData) {
//	evt, ok := e.(*event.PlayerLogin)
//	if ok == false {
//		return
//	}
//	cherryTime := ctime.Now()
//	second := cherryTime.ToSecond()
//
//	player := db.PlayerRepository.GetOrCreate(evt.PlayerId)
//	player.LastLoginTime = second
//	db.PlayerRepository.Update(player)
//
//	dotLogin := db.LogLogin{
//		ID:         evt.PlayerId,
//		FirstTime:  &second,
//		LastTime:   &second,
//		DayIndex:   cherryTime.ToShortInt32DateFormat(),
//		TotalCount: 1,
//	}
//	db.LogLoginRepository.Add(&dotLogin)
//
//	clog.Infof("[PlayerLoginEvent] [playerId = %d, onlineCount = %d]",
//		evt.PlayerId,
//		online.Count(),
//	)
//}
//
//// onLoginEvent 玩家登出事件处理
//func (p *ActorPlayers) onLogoutEvent(e cfacade.IEventData) {
//	evt, ok := e.(*event.PlayerLogout)
//	if !ok {
//		return
//	}
//	if evt.PlayerId <= 0 {
//		return
//	}
//	cherryTime := ctime.Now()
//	second := cherryTime.ToSecond()
//	player := db.PlayerRepository.GetOrCreate(evt.PlayerId)
//	player.LastLogoutTime = second
//	db.PlayerRepository.Update(player)
//	clog.Infof("[PlayerLogoutEvent] [playerId = %d, onlineCount = %d]",
//		evt.PlayerId,
//		online.Count(),
//	)
//}
//
//// onPlayerCreateEvent 玩家创建事件
//func (p *ActorPlayers) onPlayerCreateEvent(e cfacade.IEventData) {
//	evt, ok := e.(*event.PlayerCreate)
//	if !ok {
//		return
//	}
//
//	clog.Infof("[initItem] [%+v]", evt)
//}

func (p *ActorPlayers) OnStop() {
	clog.Infof("onlineCount = %d", online.Count())
}
