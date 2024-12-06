package sdk

import (
	cherryGin "gameserver/cherry/components/gin"
	cherryError "gameserver/cherry/error"
	cherryString "gameserver/cherry/extend/string"
	cfacade "gameserver/cherry/facade"
	"gameserver/internal/code"
	"gameserver/internal/data"
	rpcCenter "gameserver/internal/rpc/center"
	sessionKey "gameserver/internal/session_key"
)

type devSdk struct {
	app cfacade.IApplication
}

func (devSdk) SdkId() int32 {
	return DevMode
}

func (p devSdk) Login(_ *data.SdkRow, params Params, callback Callback) {
	pcode, _ := params.GetString("code")
	channel := params.GetInt("channel", 0)
	platform := params.GetInt("platform", 0)

	if pcode == "" || channel == 0 || platform == 0 {
		err := cherryError.Errorf("account or password params is empty.")
		callback(code.LoginError, nil, err)
		return
	}

	info, ok := rpcCenter.GetAccountInfo(p.app, int32(channel), int32(platform), pcode)
	if ok != code.OK {
		callback(code.LoginError, nil)
		return
	}

	callback(code.OK, map[string]string{
		sessionKey.OpenID:   cherryString.ToString(pcode),
		sessionKey.PlayerID: cherryString.ToString(info.Uid),
		sessionKey.ServerID: cherryString.ToString(info.ServerId),
	})
}

func (devSdk) PayCallback(_ *data.SdkRow, _ *cherryGin.Context) {
}
