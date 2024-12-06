package cache

import (
	"context"
	"fmt"
	cherryredis "gameserver/cherry/components/redis"
	jsoniter "github.com/json-iterator/go"
)

var (
	RedisServerStatus = "SERVER:"
)

type Server struct {
	ServerId   int    `json:"server_id"`   // 服务器ID
	Status     int    `json:"status"`      // 服务器状态 0 维护状态 1正常状态
	ServerHost string `json:"server_host"` // 服务器地址
	ServerPort string `json:"server_port"` // 服务器端口
}

func GetServerStatus(serverId int) (*Server, error) {
	key := fmt.Sprintf("%s%v", RedisServerStatus, serverId)
	res, err := cherryredis.GetRds().Get(context.Background(), key).Result()
	server := new(Server)
	if err != nil {
		return server, err
	}
	err = jsoniter.Unmarshal([]byte(res), &server)
	if err != nil {
		return server, err
	}
	return server, nil
}

func UpdateServerStatus(server *Server) error {
	key := fmt.Sprintf("%s%v", RedisServerStatus, server.ServerId)
	data, err := jsoniter.Marshal(server)
	if err != nil {
		return err
	}
	tx := cherryredis.GetRds().Set(context.Background(), key, data, 0)
	if tx.Err() != nil {
		return err
	}
	return nil
}
