package server_init

import (
	cherryFacade "gameserver/cherry/facade"
	"gameserver/internal/cache"
)

// Component 启动时 缓存server信息
type Component struct {
	cherryFacade.Component
}

func New() *Component {
	return &Component{}
}

func (c *Component) Name() string {
	return "run_init_server_component"
}

func (c *Component) OnAfterInit() {
	cache.InitServer()
}
