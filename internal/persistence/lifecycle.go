package persistence

import (
	clog "gameserver/cherry/logger"
	"reflect"
)

var repositories = map[any]interface{}{}

func RegisterRepository(model any, repo interface{}) {
	repositories[model] = repo
}

func Start() {
	clog.Debug("Starting repositories...")
}

func Stop() {
	flushAllRepositories()
}

func flushAllRepositories() {
	clog.Info("Flushing all repositories...")
	for _, repo := range repositories {
		repoVal := reflect.ValueOf(repo)
		flushMethod := repoVal.MethodByName("Flush")
		if flushMethod.IsValid() && flushMethod.Type().NumIn() == 0 { // 确保Flush方法存在且无参数
			flushMethod.Call(nil) // 调用Flush方法
		}
	}
}
