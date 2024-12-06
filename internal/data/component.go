package data

import (
	cherryDataConfig "gameserver/cherry/components/data-config"
	cherryMapStructure "gameserver/cherry/extend/mapstructure"
	"gameserver/internal/types"
)

var (
	SdkConfig  = &sdkConfig{}
	CodeConfig = &codeConfig{}
)

func New() *cherryDataConfig.Component {
	dataConfig := cherryDataConfig.New()
	dataConfig.Register(
		SdkConfig,
		CodeConfig,
	)
	return dataConfig
}

func DecodeData(input interface{}, output interface{}) error {
	return cherryMapStructure.HookDecode(
		input,
		output,
		"json",
		types.GetDecodeHooks(),
	)
}
