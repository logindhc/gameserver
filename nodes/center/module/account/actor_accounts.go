package account

import (
	cfacade "gameserver/cherry/facade"
	cactor "gameserver/cherry/net/actor"
)

type (
	ActorAccounts struct {
		cactor.Base
	}
)

func (p *ActorAccounts) AliasID() string {
	return "account"
}

func (p *ActorAccounts) OnFindChild(msg *cfacade.Message) (cfacade.IActor, bool) {
	// 动态创建 child actor
	childID := msg.TargetPath().ChildID
	childActor, err := p.Child().Create(childID, &ActorAccount{})

	if err != nil {
		return nil, false
	}

	return childActor, true
}
