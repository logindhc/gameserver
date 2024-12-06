package db

import (
	"context"
	"fmt"
	cherryredis "gameserver/cherry/components/redis"
	cherryString "gameserver/cherry/extend/string"
	"time"
)

var (
// uid缓存 key:uid, value:playerId
//uidCache = cache.New(
//	cache.WithMaximumSize(-1),
//	cache.WithExpireAfterAccess(60*time.Minute),
//)
//
//// 玩家表缓存 key:playerId, value:*PlayerTable
//playerTableCache = cache.New(
//	cache.WithMaximumSize(-1),
//	cache.WithExpireAfterAccess(60*time.Minute),
//)
//
//// 玩家昵称缓存 key:playerName, value:playerId
//playerNameCache = cache.New(
//	cache.WithMaximumSize(-1),
//	cache.WithExpireAfterAccess(60*time.Minute),
//)

// 英雄表缓存 key:playerId, value:*HeroTable
// 道具表缓存 key:playerId, value:*ItemTable

)
var (
	userIdRedisKey = "openId:"
	expTime        = time.Hour * 48
)

func GetLoginToken(openId string) int64 {
	key := fmt.Sprintf("%s%s", userIdRedisKey, openId)
	userStr, err := cherryredis.GetRds().Get(context.Background(), key).Result()
	if err != nil {
		return 0
	}
	return cherryString.ToInt64D(userStr, 0)
}

func SetLoginToken(openId string, userId int64) {
	key := fmt.Sprintf("%s%s", userIdRedisKey, openId)
	cherryredis.GetRds().Set(context.Background(), key, userId, expTime)
}
