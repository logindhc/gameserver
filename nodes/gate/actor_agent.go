package gate

import (
	cstring "gameserver/cherry/extend/string"
	cfacade "gameserver/cherry/facade"
	clog "gameserver/cherry/logger"
	cactor "gameserver/cherry/net/actor"
	"gameserver/cherry/net/parser/pomelo"
	cproto "gameserver/cherry/net/proto"
	"gameserver/internal/code"
	"gameserver/internal/data"
	"gameserver/internal/pb"
	rpcCenter "gameserver/internal/rpc/center"
	sessionKey "gameserver/internal/session_key"
	"gameserver/internal/token"
)

var (
	duplicateLoginCode []byte
)

type (
	// ActorAgent 每个网络连接对应一个ActorAgent
	ActorAgent struct {
		cactor.Base
	}
)

func (p *ActorAgent) OnInit() {
	duplicateLoginCode, _ = p.App().Serializer().Marshal(&cproto.I32{
		Value: code.PlayerDuplicateLogin,
	})

	p.Local().Register("login", p.login)
	p.Remote().Register("setSession", p.setSession)
}

func (p *ActorAgent) setSession(req *pb.StringKeyValue) {
	if req.Key == "" {
		return
	}

	if agent, ok := pomelo.GetAgent(p.ActorID()); ok {
		agent.Session().Set(req.Key, req.Value)
	}
}

// login 用户登录验证 (*pb.LoginResponse, int32)
func (p *ActorAgent) login(session *cproto.Session, req *pb.C2SLogin) {
	agent, found := pomelo.GetAgent(p.ActorID())
	if !found {
		return
	}

	// 验证token
	userToken, errCode := p.validateToken(req.Token)
	if code.IsFail(errCode) {
		agent.Response(session, errCode)
		return
	}

	// 验证channelId是否配置
	sdkRow := data.SdkConfig.Get(userToken.Channel)
	if sdkRow == nil {
		agent.ResponseCode(session, code.ChannelIDError, true)
		return
	}

	//// 根据token带来的sdk参数，从中心节点获取uid
	info, errCode := rpcCenter.GetAccountInfo(p.App(), userToken.Channel, userToken.Platform, userToken.OpenId)
	if info == nil || code.IsFail(errCode) {
		agent.ResponseCode(session, errCode, true)
		return
	}

	//玩家id
	uid := info.Uid
	p.checkGateSession(uid)

	if err := agent.Bind(uid); err != nil {
		clog.Warn(err)
		agent.ResponseCode(session, code.AccountBindFail, true)
		return
	}
	agent.Session().Set(sessionKey.AccountID, cstring.ToString(info.AccountId))
	agent.Session().Set(sessionKey.ServerID, cstring.ToString(info.ServerId))
	agent.Session().Set(sessionKey.ChannelID, cstring.ToString(userToken.Channel))
	agent.Session().Set(sessionKey.PlatformID, cstring.ToString(userToken.Platform))
	agent.Session().Set(sessionKey.OpenID, cstring.ToString(userToken.OpenId))

	response := &pb.S2CLogin{
		Uid:    uid,
		Params: nil,
	}
	clog.Infof("login uid = %d", uid)
	agent.Response(session, response)
}

func (p *ActorAgent) validateToken(base64Token string) (*token.Token, int32) {
	userToken, ok := token.DecodeToken(base64Token)
	if ok == false {
		return nil, code.AccountTokenValidateFail
	}

	sdkRow := data.SdkConfig.Get(userToken.Channel)
	if sdkRow == nil {
		return nil, code.ChannelIDError
	}

	statusCode, ok := token.Validate(userToken, sdkRow.Salt)
	if ok == false {
		return nil, statusCode
	}

	return userToken, code.OK
}

func (p *ActorAgent) checkGateSession(uid cfacade.UID) {
	if agent, found := pomelo.GetAgentWithUID(uid); found {
		agent.Kick(duplicateLoginCode, true)
	}

	rsp := &cproto.PomeloKick{
		Uid:    uid,
		Reason: duplicateLoginCode,
	}

	// 遍历所有网关节点，踢除旧的session
	members := p.App().Discovery().ListByType(p.App().NodeType(), p.App().NodeId())
	for _, member := range members {
		// user是gate.go里自定义的agentActorID
		actorPath := cfacade.NewPath(member.GetNodeId(), "user")
		p.Call(actorPath, pomelo.KickFuncName, rsp)
	}
}

// onSessionClose  当agent断开时，关闭对应的ActorAgent
func (p *ActorAgent) onSessionClose(agent *pomelo.Agent) {
	session := agent.Session()
	serverId := session.GetString(sessionKey.ServerID)
	if serverId == "" {
		return
	}

	// 通知game节点关闭session
	childId := cstring.ToString(session.Uid)
	if childId != "" {
		targetPath := cfacade.NewChildPath(serverId, "player", childId)
		p.Call(targetPath, "sessionClose", nil)
	}

	// 自己退出
	p.Exit()
	clog.Infof("sessionClose path = %s", p.Path())
}
