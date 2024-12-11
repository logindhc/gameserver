package utils

import (
	"context"
	cherryredis "gameserver/cherry/components/redis"
	clog "gameserver/cherry/logger"
	"gameserver/internal/constant"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func UpdateGameNodeOnline(nodeId string, online float64) {
	err := cherryredis.GetRds().ZAdd(ctx, constant.GameOnlineKey, &redis.Z{Score: online, Member: nodeId}).Err()
	if err != nil {
		clog.Errorf("UpdateOnline [nodeId = %s] redis error: %v", nodeId, err)
	}
}

func DelGameNodeOnline(nodeId string) {
	err := cherryredis.GetRds().ZRem(ctx, constant.GameOnlineKey, nodeId).Err()
	if err != nil {
		clog.Errorf("DelOnline [nodeId = %s] redis error: %v", nodeId, err)
	}
}

func GetAllGameNodeIdByRank() ([]string, error) {
	// 按分数从小到到排序获取所有元素
	result, err := cherryredis.GetRds().ZRangeWithScores(ctx, constant.GameOnlineKey, 0, -1).Result()
	if err != nil {
		clog.Error("GetAllOnlineByRank redis error: ", err)
		return nil, err
	}
	var nodeIds = make([]string, 0, len(result))
	for _, z := range result {
		nodeIds = append(nodeIds, z.Member.(string))
	}
	return nodeIds, nil
}
