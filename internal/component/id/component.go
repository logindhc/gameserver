package idgener

import (
	cherryFacade "gameserver/cherry/facade"
	"gameserver/internal/cache"
)

// 默认的玩家ID生成器
var PlayerIdGenerator cache.PlayerIDGenerator

// Component 启动时从redis获取自增ID
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
	PlayerIdGenerator = cache.NewIDGenerator()
}

// ParseID 根据生成的 ID 反推区服 ID 和自增值
func ParseID(id int64) (int32, int64) {
	serverID := id & ((1 << cache.ServerIDBits) - 1) // 提取低位区服 ID
	increment := id >> cache.ServerIDBits            // 提取高位自增值
	return int32(serverID), increment
}
