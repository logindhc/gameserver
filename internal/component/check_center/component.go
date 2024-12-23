package checkCenter

import (
	rpcCenter "gameserver/internal/rpc/center"
	"time"

	cherryFacade "gameserver/cherry/facade"
	cherryLogger "gameserver/cherry/logger"
)

// Component 启动时,检查center节点是否存活
type Component struct {
	cherryFacade.Component
}

func New() *Component {
	return &Component{}
}

func (c *Component) Name() string {
	return "run_check_component"
}

func (c *Component) OnAfterInit() {
	for {
		if rpcCenter.Ping(c.App()) {
			break
		}
		time.Sleep(2 * time.Second)
		cherryLogger.Warn("center node connect fail. retrying in 2 seconds.")
	}
}
