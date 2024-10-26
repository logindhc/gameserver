package game

import (
	"gameserver/cherry"
	cherryCron "gameserver/cherry/components/cron"
	cherryGops "gameserver/cherry/components/gops"
	cherrySnowflake "gameserver/cherry/extend/snowflake"
	cstring "gameserver/cherry/extend/string"
	cherryUtils "gameserver/cherry/extend/utils"
	checkCenter "gameserver/internal/component/check_center"
	"gameserver/internal/data"
	"gameserver/nodes/game/db"
	"gameserver/nodes/game/module/player"
)

func Run(profileFilePath, nodeId string) {
	if !cherryUtils.IsNumeric(nodeId) {
		panic("node parameter must is number.")
	}

	// snowflake global id
	serverId, _ := cstring.ToInt64(nodeId)
	cherrySnowflake.SetDefaultNode(serverId)

	// 配置cherry引擎
	app := cherry.Configure(profileFilePath, nodeId, false, cherry.Cluster)

	// diagnose
	app.Register(cherryGops.New())
	// 注册调度组件
	app.Register(cherryCron.New())
	// 注册数据配置组件
	app.Register(data.New())
	// 注册检测中心节点组件，确认中心节点启动后，再启动当前节点
	app.Register(checkCenter.New())
	// 注册db组件
	app.Register(db.New())

	app.AddActors(
		&player.ActorPlayers{},
	)

	app.Startup()
}
