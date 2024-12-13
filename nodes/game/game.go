package game

import (
	"bufio"
	"fmt"
	"gameserver/cherry"
	cherryCron "gameserver/cherry/components/cron"
	cherryGops "gameserver/cherry/components/gops"
	cherryredis "gameserver/cherry/components/redis"
	cherrySnowflake "gameserver/cherry/extend/snowflake"
	cstring "gameserver/cherry/extend/string"
	cherryUtils "gameserver/cherry/extend/utils"
	clog "gameserver/cherry/logger"
	"gameserver/hotfix"
	"gameserver/hotfix/symbols"
	checkCenter "gameserver/internal/component/check_center"
	serverinit "gameserver/internal/component/init_server"
	"gameserver/internal/data"
	"gameserver/nodes/game/db"
	"gameserver/nodes/game/job"
	"gameserver/nodes/game/module/player"
	"os"
	"strings"
)

func Run(profileFilePath, nodeId string) {
	if !cherryUtils.IsNumeric(nodeId) {
		panic("node parameter must is number.")
	}
	if cstring.ToIntD(nodeId) >= 1024 || cstring.ToIntD(nodeId) < 1 {
		panic("node parameter nodeId err.")
	}
	// 配置cherry引擎
	app := cherry.Configure(profileFilePath, nodeId, false, cherry.Cluster)

	// snowflake global id
	serverId, _ := cstring.ToInt64(nodeId)
	cherrySnowflake.SetDefaultNode(serverId)

	// diagnose
	app.Register(cherryGops.New())
	// 注册调度组件
	app.Register(cherryCron.New())
	// 注册redis组件
	app.Register(cherryredis.New())
	// 注册数据配置组件
	app.Register(data.New())
	// 注册检测中心节点组件，确认中心节点启动后，再启动当前节点
	app.Register(checkCenter.New())
	// 注册db组件
	app.Register(db.New())
	// 注册初始化game节点缓存组件
	app.Register(serverinit.New())

	app.AddActors(
		&job.ActorJob{},
		&player.ActorPlayers{},
		//&item.ActorItems{},
	)

	go scanner()
	app.Startup()

}

func scanner() {
	// 从标准输入流中接收输入数据
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		split := strings.Split(line, " ")
		if split[0] == "hotfix" {
			if len(split) < 3 {
				fmt.Println("hotfix 脚本路径 函数名")
				fmt.Println("例如: hotfix gameserver.go.patch gameserver.GetPatch()")
				continue
			}
			filePath := split[1] // 补丁脚本的路径
			evalText := split[2] // 补丁脚本内执行的函数名
			clog.Info("hotfix file:", filePath, "eval:", evalText)
			// 加载补丁函数foo.GetPatch()
			_, err := hotfix.ApplyFunc(filePath, evalText, symbols.Symbols)
			if err != nil {
				clog.Error(err)
				continue
			}
			clog.Info("hotfix success")
		}
	}
}
