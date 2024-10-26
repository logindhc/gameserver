package center

import (
	"gameserver/cherry"
	"gameserver/cherry/components/cron"
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

	app.Register(cherryCron.New())
	app.Register(data.New())
	app.Register(db.New())

	app.AddActors(
		&account.ActorAccount{},
		&ops.ActorOps{},
	)

	app.Startup()
}
