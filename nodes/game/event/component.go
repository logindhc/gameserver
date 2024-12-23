package eventmgr

import (
	cfacade "gameserver/cherry/facade"
)

var Instance *EventBus

type Component struct {
	cfacade.Component
	instance *EventBus
}

func (c *Component) Name() string {
	return "player_event_component"
}

// Init 组件初始化函数
func (c *Component) Init() {
}

func New() *Component {
	c := &Component{instance: NewEventBus()}
	Instance = c.instance
	return c
}
