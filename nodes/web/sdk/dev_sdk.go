package sdk

import (
	"gameserver/internal/code"
	"gameserver/internal/data"
	rpcCenter "gameserver/internal/rpc/center"
	cherryGin "github.com/cherry-game/cherry/components/gin"
	cherryError "github.com/cherry-game/cherry/error"
	cherryString "github.com/cherry-game/cherry/extend/string"
	cfacade "github.com/cherry-game/cherry/facade"
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
