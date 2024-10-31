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

	defaultModels = []interface{}{
		&DevAccountTable{},
		&UserBindTable{},
	}
)

type Component struct {
	*cherryGORM.Component
}

func (c *Component) Name() string {
	return "db_center_component"
}

func (c *Component) Init() {
	c.Component.Init()
}

func (c *Component) OnAfterInit() {
	c.AutoMigrate(defaultModels, nil, false)
	persistence.Start(defaultModels)

	addOnload(loadDevAccount)
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
