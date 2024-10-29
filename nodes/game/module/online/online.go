package online

import (
	cherryFacade "gameserver/cherry/facade"
	clog "gameserver/cherry/logger"
	"sync"
)

var (
	currentServerId int32 = 0 // 当前游戏节点的serverId
)

var (
	lock        = &sync.RWMutex{}
	playerIdMap = make(map[int64]string) // key:playerId, value:agentActorPath
)

func BindPlayer(playerId int64, agentActorPath string) {
	if playerId < 1 || agentActorPath == "" {
		return
	}

	lock.Lock()
	defer lock.Unlock()

	playerIdMap[playerId] = agentActorPath
}

func UnBindPlayer(playerId cherryFacade.UID) int64 {
	if playerId < 1 {
		return 0
	}

	lock.Lock()
	defer lock.Unlock()

	delete(playerIdMap, playerId)

	playerIdCount := len(playerIdMap)

	if playerIdCount == 0 {
		clog.Infof("Unbind player uid = %d, playerIdCount = %d", playerId, playerIdCount)
	}

	return playerId
}

func Count() int {
	lock.Lock()
	defer lock.Unlock()

	count := len(playerIdMap)
	return count
}
