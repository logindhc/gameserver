package account

import (
	"fmt"
	cstring "gameserver/cherry/extend/string"
	cactor "gameserver/cherry/net/actor"
	"gameserver/internal/code"
	"gameserver/internal/constant"
	"gameserver/internal/pb"
	"gameserver/nodes/center/db"
	"math/rand"
)

type (
	ActorAccount struct {
		cactor.Base
	}
)

// OnInit center为后端节点，不直接与客户端通信，所以了一些remote函数，供RPC调用
func (p *ActorAccount) OnInit() {
	p.Remote().Register("getAccountInfo", p.getAccountInfo)
}

// getAccountInfo 获取uid
func (p *ActorAccount) getAccountInfo(req *pb.AccountInfo) (*pb.AccountInfo, int32) {
	id := fmt.Sprintf("%d_%s", req.Channel, req.OpenId)
	account := db.AccountRepository.Get(id)
	if account != nil {
		req.Uid = account.UID
		req.ServerId = account.ServerId
		return req, code.OK
	}
	members := p.App().Discovery().ListByType(constant.GameNodeType)
	if len(members) == 0 {
		return nil, code.ServerError
	}
	//	TODO 随机选一个游戏服，需要优化负载均衡
	sid := members[rand.Intn(len(members))].GetNodeId()
	serverId := cstring.ToInt32D(sid)
	account = db.CreateAccount(req.Channel, req.OpenId, req.Platform, serverId)

	if account.UID == 0 || serverId == 0 {
		return nil, code.AccountAuthFail
	}
	req.Uid = account.UID
	req.ServerId = serverId
	return req, code.OK
}
