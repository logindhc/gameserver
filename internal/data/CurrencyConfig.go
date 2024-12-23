package data

import (
	cherryError "gameserver/cherry/error"
	cherryLogger "gameserver/cherry/logger"
)

var CurrencyConfig = &currencyConfig{}

type (
	CurrencyRow struct {
		Id int `json:"id"` // 货币类型
	}

	currencyConfig struct {
		maps map[int]*CurrencyRow
	}
)

func (c *currencyConfig) Init() {
	c.maps = make(map[int]*CurrencyRow)
}

func (c *currencyConfig) OnLoad(maps interface{}, _ bool) (int, error) {
	list, ok := maps.([]interface{}) // map结构：maps.(map[string]interface{})
	if !ok {
		return 0, cherryError.Error("maps convert to map[string]interface{} error.")
	}

	loadMaps := make(map[int]*CurrencyRow)
	for index, data := range list {
		loadConfig := &CurrencyRow{}
		err := DecodeData(data, loadConfig)
		if err != nil {
			cherryLogger.Warnf("decode error. [id = %v, %v], err = %s", index, loadConfig, err)
			continue
		}
		loadMaps[loadConfig.Id] = loadConfig
	}
	c.maps = loadMaps

	return len(list), nil
}

func (c *currencyConfig) Name() string {
	return "CurrencyConfig"
}

func (c *currencyConfig) OnAfterLoad(_ bool) {}

func (c *currencyConfig) Get(key int) (*CurrencyRow, bool) {
	row, ok := c.maps[key]
	return row, ok
}

func (c *currencyConfig) List() []*CurrencyRow {
	var list []*CurrencyRow
	for _, row := range c.maps {
		list = append(list, row)
	}
	return list
}
