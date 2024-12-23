package resmgr

import (
	cfacade "gameserver/cherry/facade"
	"gameserver/nodes/game/facade"
)

var Instance facade.IResManager

type Component struct {
	cfacade.Component
	instance facade.IResManager
}

func (c *Component) Name() string {
	return "game_res_component"
}

// Init 组件初始化函数
func (c *Component) Init() {
	c.instance.Register(&currencyRes{})
	c.instance.Register(&itemRes{})
	c.instance.Register(&heroRes{})
	c.instance.Register(&boxRes{})
}

func New() *Component {
	c := &Component{instance: &resManager{
		resMap: make(map[int]facade.IRes),
	}}
	Instance = c.instance
	return c
}
