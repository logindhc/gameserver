package rpcCenter

import (
	"fmt"
	cfacade "gameserver/cherry/facade"
	clog "gameserver/cherry/logger"
	"gameserver/internal/code"
	"gameserver/internal/constant"
	"gameserver/internal/pb"
	"gameserver/internal/rpc"
)

// route = 节点类型.节点handler.remote函数

const (
	opsActor     = ".ops"
	accountActor = ".account.%d_%s"
)

const (
	ping           = "ping"
	getAccountInfo = "getAccountInfo"
)

const (
	sourcePath = ".system"
)

// Ping 访问center节点，确认center已启动
func Ping(app cfacade.IApplication) bool {
	nodeId := GetCenterNodeID(app)
	if nodeId == "" {
		return false
	}

	rsp := &pb.Bool{}
	targetPath := nodeId + opsActor
	errCode := app.ActorSystem().CallWait(sourcePath, targetPath, ping, nil, rsp)
	if code.IsFail(errCode) {
		return false
	}

	return rsp.Value
}

// GetAccountInfo 获取帐号UID和区服信息
func GetAccountInfo(app cfacade.IApplication, channel, platform int32, openId string) (*rpc.AccountInfo, int32) {
	req := &rpc.AccountInfo{
		Channel:  channel,
		Platform: platform,
		OpenId:   openId,
	}

	targetPath := GetTargetPath(app, fmt.Sprintf(accountActor, channel, openId))
	rsp := &rpc.AccountInfo{}
	errCode := app.ActorSystem().CallWait(sourcePath, targetPath, getAccountInfo, req, rsp)
	if code.IsFail(errCode) {
		clog.Warnf("[GetAccountInfo] errCode = %v", errCode)
		return nil, errCode
	}
	return rsp, code.OK
}

func GetCenterNodeID(app cfacade.IApplication) string {
	list := app.Discovery().ListByType(constant.CenterType)
	if len(list) > 0 {
		return list[0].GetNodeId()
	}
	return ""
}

func GetTargetPath(app cfacade.IApplication, actorID string) string {
	nodeId := GetCenterNodeID(app)
	return nodeId + actorID
}
