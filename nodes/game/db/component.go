package db

import (
	cherryGORM "gameserver/cherry/components/gorm"
	cherryUtils "gameserver/cherry/extend/utils"
	cherryLogger "gameserver/cherry/logger"
	"gameserver/internal/persistence"
)

var (
	onLoadFuncList []func() // db初始化时加载函数列表
	database       *Component
	// 表自动生成手动注册
	defaultModels = []interface{}{
		&PlayerTable{},
	}
	logModels = []interface{}{
		&DotLogin{},
	}
)

type Component struct {
	*cherryGORM.Component
}

func (c *Component) Name() string {
	return "db_game_component"
}

// Init 组件初始化函数
func (c *Component) Init() {
	c.Component.Init()
}

func (c *Component) OnAfterInit() {
	c.AutoMigrate(defaultModels, logModels)
	persistence.Start(defaultModels)
	persistence.Start(logModels)
	for _, fn := range onLoadFuncList {
		cherryUtils.Try(fn, func(errString string) {
			cherryLogger.Warnf(errString)
		})
	}
}

func (*Component) OnStop() {
	persistence.Stop()
}

func New() *Component {
	c := &Component{cherryGORM.NewComponent()}
	database = c
	return c
}

func addOnload(fn func()) {
	onLoadFuncList = append(onLoadFuncList, fn)
}
