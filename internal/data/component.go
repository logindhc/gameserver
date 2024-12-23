package data

import (
	cherryDataConfig "gameserver/cherry/components/data-config"
	cherryMapStructure "gameserver/cherry/extend/mapstructure"
	"gameserver/internal/types"
)

func New() *cherryDataConfig.Component {
	dataConfig := cherryDataConfig.New()
	dataConfig.Register(
		CodeConfig,
		SdkConfig,
		ShopBoxConfig,
		ShopBoxLvConfig,
		DropConfig,
		ItemConfig,
		LevelRewardConfig,
		GameConfig,
		HeroConfig,
		ResTypeConfig,
		CurrencyConfig,
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
