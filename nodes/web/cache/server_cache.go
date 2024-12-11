package cache

import (
	"context"
	"fmt"
	cherryredis "gameserver/cherry/components/redis"
	cherryLogger "gameserver/cherry/logger"
	jsoniter "github.com/json-iterator/go"
	"sync"
)

var (
	RedisServer = "admin:server:"
	serverMaps  = sync.Map{}
)

func InitServer() {
	keys, err := cherryredis.GetRds().Keys(context.Background(), RedisServer+"*").Result()
	if err != nil {
		return
	}
	for _, key := range keys {
		res, err := cherryredis.GetRds().Get(context.Background(), key).Result()
		server := new(Server)
		if err != nil {
			continue
		}
		err = jsoniter.Unmarshal([]byte(res), &server)
		if err != nil {
			continue
		}
		serverMaps.Store(server.ServerId, server)
		cherryLogger.Infof("init server info. server=%v", server)
	}
}

type Server struct {
	ServerId   int32  `json:"server_id"`   // 服务器ID
	Status     int32  `json:"status"`      // 服务器状态 0 维护状态 1正常状态
	ServerHost string `json:"server_host"` // 服务器地址
	ServerPort string `json:"server_port"` // 服务器端口
	CheckHost  string `json:"check_host"`  // 提审服地址
	CheckPort  string `json:"check_port"`  // 提审服端口
}

func GetServerInfo(serverId int32) (*Server, error) {
	value, ok := serverMaps.Load(serverId)
	if ok {
		return value.(*Server), nil
	}
	key := fmt.Sprintf("%s%v", RedisServer, serverId)
	res, err := cherryredis.GetRds().Get(context.Background(), key).Result()
	server := new(Server)
	if err != nil {
		cherryLogger.Warnf("get server info error. serverId=%v, error=%s", serverId, err)
		return server, err
	}
	err = jsoniter.Unmarshal([]byte(res), &server)
	if err != nil {
		cherryLogger.Warnf("unmarshal server info error. serverId=%v, error=%s", serverId, err)
		return server, err
	}
	return server, nil
}

func UpdateServerInfo(server *Server) error {
	key := fmt.Sprintf("%s%v", RedisServer, server.ServerId)
	data, err := jsoniter.Marshal(server)
	if err != nil {
		cherryLogger.Warnf("update server info error. server=%v, error=%s", server, err)
		return err
	}
	tx := cherryredis.GetRds().Set(context.Background(), key, data, 0)
	if tx.Err() != nil {
		cherryLogger.Warnf("update server info error. server=%v, error=%s", server, tx.Err())
		return err
	}
	serverMaps.Store(server.ServerId, server)
	cherryLogger.Infof("update server info. server=%v", server)
	return nil
}

func GetAllServer() ([]*Server, error) {
	var serverList []*Server
	serverMaps.Range(func(key, value interface{}) bool {
		server := value.(*Server)
		serverList = append(serverList, server)
		return true
	})
	return serverList, nil
}
