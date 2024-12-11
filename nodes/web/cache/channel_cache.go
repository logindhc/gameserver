package cache

import (
	"context"
	"fmt"
	cherryredis "gameserver/cherry/components/redis"
	jsoniter "github.com/json-iterator/go"
)

var (
	RedisChannel = "admin:channel:"
)

type ChannelInfo struct {
	Id           int32           `json:"id"`           // 渠道ID
	Name         string          `json:"name"`         // 渠道名称
	ServerIds    map[string]bool `json:"server_ids"`   // 服务器ID
	Version      string          `json:"version"`      // 审核版本
	TipVersion   string          `json:"tipVersion"`   //提示更新最小版本
	ForceVersion string          `json:"forceVersion"` //强制更新最小版本
}

func GetChannelInfo(id int32) (*ChannelInfo, error) {
	key := fmt.Sprintf("%s%v", RedisChannel, id)
	res, err := cherryredis.GetRds().Get(context.Background(), key).Result()
	channel := new(ChannelInfo)
	if err != nil {
		return channel, err
	}
	err = jsoniter.Unmarshal([]byte(res), &channel)
	if err != nil {
		return channel, err
	}
	return channel, nil
}

func UpdateChannelInfo(channel *ChannelInfo) error {
	key := fmt.Sprintf("%s%v", RedisChannel, channel.Id)
	data, err := jsoniter.Marshal(channel)
	if err != nil {
		return err
	}
	tx := cherryredis.GetRds().Set(context.Background(), key, data, 0)
	if tx.Err() != nil {
		return err
	}
	return nil
}

func DelChannelInfo(id int32) error {
	key := fmt.Sprintf("%s%v", RedisChannel, id)
	tx := cherryredis.GetRds().Del(context.Background(), key)
	if tx.Err() != nil {
		return tx.Err()
	}
	return nil
}

func UpdateAppCheckInfo(channelId int32, version string, tipVersion string, forceVersion string) error {
	channel, err := GetChannelInfo(channelId)
	if err != nil {
		//不应该获取不到
		return err
	}
	channel.Version = version
	channel.TipVersion = tipVersion
	channel.ForceVersion = forceVersion
	key := fmt.Sprintf("%s%v", RedisChannel, channelId)
	data, err := jsoniter.Marshal(channel)
	if err != nil {
		return err
	}
	tx := cherryredis.GetRds().Set(context.Background(), key, data, 0)
	if tx.Err() != nil {
		return err
	}
	return nil
}

func DelAppCheckInfo(channelId int32) error {
	channel, err := GetChannelInfo(channelId)
	if err != nil {
		//不应该获取不到
		return err
	}
	channel.Version = ""
	channel.TipVersion = ""
	channel.ForceVersion = ""
	key := fmt.Sprintf("%s%v", RedisChannel, channelId)
	data, err := jsoniter.Marshal(channel)
	if err != nil {
		return err
	}
	tx := cherryredis.GetRds().Set(context.Background(), key, data, 0)
	if tx.Err() != nil {
		return err
	}
	return nil
}
