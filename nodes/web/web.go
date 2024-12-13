package web

import (
	"gameserver/cherry"
	cherryCron "gameserver/cherry/components/cron"
	cherryGin "gameserver/cherry/components/gin"
	cherryredis "gameserver/cherry/components/redis"
	cherryFile "gameserver/cherry/extend/file"
	checkCenter "gameserver/internal/component/check_center"
	serverinit "gameserver/internal/component/init_server"
	"gameserver/internal/data"
	"gameserver/nodes/web/controller"
	"gameserver/nodes/web/sdk"
	"github.com/gin-gonic/gin"
)

func Run(profileFilePath, nodeId string) {
	// 配置cherry引擎,加载profile配置文件
	app := cherry.Configure(profileFilePath, nodeId, false, cherry.Cluster)

	// 注册调度组件
	app.Register(cherryCron.New())

	// 注册检查中心服是否启动组件
	app.Register(checkCenter.New())

	// 注册数据配表组件
	app.Register(data.New())

	// 注册redis组件
	app.Register(cherryredis.New())

	// 注册初始化game节点缓存组件
	app.Register(serverinit.New())

	// 加载http server组件
	httpServerComponent(app)

	// 加载sdk逻辑
	sdk.Init(app)

	// 启动cherry引擎
	app.Startup()
}

func httpServerComponent(app *cherry.AppBuilder) {
	gin.SetMode(gin.DebugMode)

	// new http server
	httpServer := cherryGin.NewHttp("http_server", app.Address())
	httpServer.Use(cherryGin.Cors())

	httpServer.Use(cherryGin.Md5Filter(app.Settings().GetString("md5_key", "")))

	// http server使用gin组件搭建，这里增加一个RecoveryWithZap中间件
	httpServer.Use(cherryGin.RecoveryWithZap(true))

	// 映射h5客户端静态文件到static目录，例如：http://127.0.0.1/static/protocol.js
	httpServer.Static("/static", "./static/")

	// 加载./view目录的html模板文件(详情查看gin文档)
	viewFiles := cherryFile.WalkFiles("./view/", ".html")
	if len(viewFiles) < 1 {
		panic("view files not found.")
	}
	httpServer.LoadHTMLFiles(viewFiles...)

	//注册 controller
	httpServer.Register(new(controller.Controller),
		new(controller.GMController),
		new(controller.MailController))
	// 注册 http server
	app.Register(httpServer)
}
