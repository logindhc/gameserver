package persistence

import (
	clog "gameserver/cherry/logger"
	"reflect"
)

import (
	"container/list"
)

var (
	// 注册所有数据库实体
	repositories = list.New()
)

// 初始化所有实体的接口
type IModel interface {
	//反射调用InitRepository方法
	InitRepository()
}

func RegisterRepository(repo any) {
	repositories.PushBack(repo)
	clog.Debugf("Registered repository for %v", repo)
}

func Start(models []interface{}) {
	clog.Debug("Starting repositories...")
	for _, model := range models {
		repoVal := reflect.ValueOf(model)
		initMethod := repoVal.MethodByName("InitRepository")
		if initMethod.IsValid() && initMethod.Type().NumIn() == 0 { // 确保InitRepository方法存在且无}
			initMethod.Call(nil)
		}
	}
}

func Stop() {
	flushAllRepositories()
}

func flushAllRepositories() {
	clog.Debug("Flushing all repositories...")
	for i := repositories.Front(); i != nil; i = i.Next() {
		model := i.Value
		repoVal := reflect.ValueOf(model)
		flushMethod := repoVal.MethodByName("Flush")
		if flushMethod.IsValid() && flushMethod.Type().NumIn() == 0 { // 确保Flush方法存在且无参数
			flushMethod.Call(nil) // 调用Flush方法
		}
	}
}
