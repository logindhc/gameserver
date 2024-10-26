package master

import "gameserver/cherry"

func Run(profileFilePath, nodeId string) {
	app := cherry.Configure(profileFilePath, nodeId, false, cherry.Cluster)
	app.Startup()
}
