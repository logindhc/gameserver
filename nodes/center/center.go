package center

import (
	"gameserver/cherry"
	"gameserver/cherry/components/cron"
	cherryredis "gameserver/cherry/components/redis"
	cherrySnowflake "gameserver/cherry/extend/snowflake"
	idgener "gameserver/internal/component/id"
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
	// 在DB组件之前注册玩家id生成组件,需要指定表来初始化自增ID
	app.Register(idgener.New())
	// 注册数据库组件
	app.Register(db.New())
	// 注册redis组件
	app.Register(cherryredis.New())

	app.AddActors(
		&account.ActorAccounts{},
		&ops.ActorOps{},
	)

	app.Startup()

}
