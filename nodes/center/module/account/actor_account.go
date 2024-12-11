package account

import (
	"fmt"
	cstring "gameserver/cherry/extend/string"
	cherryLogger "gameserver/cherry/logger"
	cactor "gameserver/cherry/net/actor"
	"gameserver/internal/code"
	"gameserver/internal/constant"
	"gameserver/internal/pb"
	"gameserver/internal/utils"
	"gameserver/nodes/center/db"
)

type (
	ActorAccount struct {
		cactor.Base
	}
)

// OnInit center为后端节点，不直接与客户端通信，所以了一些remote函数，供RPC调用
func (a *ActorAccount) OnInit() {
	a.Remote().Register("getAccountInfo", a.getAccountInfo)
}

// getAccountInfo 获取uid
func (a *ActorAccount) getAccountInfo(req *pb.AccountInfo) (*pb.AccountInfo, int32) {
	id := fmt.Sprintf("%d_%d_%s", req.Channel, req.Platform, req.OpenId)
	account := db.AccountRepository.Get(id)
	if account != nil {
		req.Uid = account.UID
		req.ServerId = account.ServerId
		return req, code.OK
	}
	serverId := a.getServerId()
	account = db.CreateAccount(id, req.Channel, req.OpenId, req.Platform, serverId)

	if account.UID == 0 || serverId == 0 {
		return nil, code.AccountAuthFail
	}
	req.Uid = account.UID
	req.ServerId = serverId
	return req, code.OK
}

func (a *ActorAccount) getServerId() int32 {
	//根据最小负载的game节点
	serverId := int32(0)
	nodeIds, err := utils.GetAllGameNodeIdByRank()
	if err != nil {
		cherryLogger.Warnf("get game node id error. error=%s", err)
		return serverId
	}
	// 避免节点掉线，影响新用户
	members := a.App().Discovery().ListByType(constant.GameNodeType)
	for _, nodeId := range nodeIds {
		for _, member := range members {
			if member.GetNodeId() == nodeId {
				serverId = cstring.ToInt32D(nodeId)
				return serverId
			}
		}
	}
	return serverId
}
