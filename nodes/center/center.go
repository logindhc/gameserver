package center

import (
	"gameserver/cherry"
	"gameserver/cherry/components/cron"
	cherryredis "gameserver/cherry/components/redis"
	cherrySnowflake "gameserver/cherry/extend/snowflake"
	"gameserver/internal/data"
	"gameserver/nodes/center/db"
	"gameserver/nodes/center/module/account"
	"gameserver/nodes/center/module/ops"
)

func Run(profileFilePath, nodeId string) {
	app := cherry.Configure(
		profileFilePath,
		nodeId,
		false,
		cherry.Cluster,
	)

	cherrySnowflake.InitDefaultNode(nodeId)

	app.Register(cherryCron.New())
	app.Register(data.New())
	app.Register(db.New())
	// 注册redis组件
	app.Register(cherryredis.New())

	app.AddActors(
		&account.ActorAccounts{},
		&ops.ActorOps{},
	)

	app.Startup()

}
