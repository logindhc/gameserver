package item

import (
	cfacade "gameserver/cherry/facade"
	clog "gameserver/cherry/logger"
	"gameserver/cherry/net/parser/pomelo"
)

type (
	// ActorItems 每位玩家对应一个子actor
	ActorItems struct {
		pomelo.ActorBase
	}
)

func (p *ActorItems) AliasID() string {
	return "item"
}

func (p *ActorItems) OnInit() {
	clog.Debugf("[ActorItems] path = %s init!", p.PathString())
}

func (p *ActorItems) OnFindChild(msg *cfacade.Message) (cfacade.IActor, bool) {
	// 动态创建 player child actor
	childID := msg.TargetPath().ChildID
	childActor, err := p.Child().Create(childID, &ActorItem{})
	if err != nil {
		return nil, false
	}

	return childActor, true
}

func (p *ActorItems) OnStop() {
	clog.Debugf("[ActorItems] path = %s exit!", p.PathString())
}
