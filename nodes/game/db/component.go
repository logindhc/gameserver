package db

import (
	cherryGORM "gameserver/cherry/components/gorm"
	cherryTime "gameserver/cherry/extend/time"
	cherryUtils "gameserver/cherry/extend/utils"
	cherryLogger "gameserver/cherry/logger"
	"gameserver/internal/job"
	"gameserver/internal/persistence"
	"time"
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
	// 记录最近触发的月份
	lastJobTime = 0
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
	c.AutoMigrate(defaultModels, logModels, false)
	persistence.Start(defaultModels)
	persistence.Start(logModels)

	for _, fn := range onLoadFuncList {
		cherryUtils.Try(fn, func(errString string) {
			cherryLogger.Warnf(errString)
		})
	}
	CheckLogJob(c)
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

// 每小时检查一次，如果当月检查一次就不会重复检查，每月日志表生成
func CheckLogJob(c *Component) {
	//开服就检查一次
	c.AutoMigrate(nil, logModels, true)

	job.GlobalTimer.BuildEveryFunc(time.Hour, func() {
		intMonthTime := cherryTime.Now().ToShortIntMonthFormat()
		if lastJobTime < intMonthTime {
			c.AutoMigrate(nil, logModels, true)
			lastJobTime = intMonthTime
		}
	})
}
