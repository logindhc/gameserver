package sdk

import (
	cherryGin "gameserver/cherry/components/gin"
	cherryError "gameserver/cherry/error"
	cherryString "gameserver/cherry/extend/string"
	cfacade "gameserver/cherry/facade"
	"gameserver/internal/code"
	"gameserver/internal/data"
	rpcCenter "gameserver/internal/rpc/center"
)

type devSdk struct {
	app cfacade.IApplication
}

func (devSdk) SdkId() int32 {
	return DevMode
}

func (p devSdk) Login(_ *data.SdkRow, params Params, callback Callback) {
	accountName, _ := params.GetString("account")
	password, _ := params.GetString("password")

	if accountName == "" || password == "" {
		err := cherryError.Errorf("account or password params is empty.")
		callback(code.LoginError, nil, err)
		return
	}

	accountId := rpcCenter.GetDevAccount(p.app, accountName, password)
	if accountId < 1 {
		callback(code.LoginError, nil)
		return
	}

	callback(code.OK, map[string]string{
		"open_id": cherryString.ToString(accountId),
	})
}

func (devSdk) PayCallback(_ *data.SdkRow, _ *cherryGin.Context) {
}
