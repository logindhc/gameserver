package data

import (
	cherryError "gameserver/cherry/error"
	cherryLogger "gameserver/cherry/logger"
)

var ShopBoxLvConfig = &shopBoxLvConfig{}

type (
	ShopBoxLvRow struct {
		Id       int `json:"id"`       //
		Exp      int `json:"exp"`      // 经验
		MinBoxId int `json:"minBoxId"` // 等级奖励(道具id）
		MaxBoxId int `json:"maxBoxId"` // 等级奖励
	}

	shopBoxLvConfig struct {
		maps map[int]*ShopBoxLvRow
	}
)

func (c *shopBoxLvConfig) Init() {
	c.maps = make(map[int]*ShopBoxLvRow)
}

func (c *shopBoxLvConfig) OnLoad(maps interface{}, _ bool) (int, error) {
	list, ok := maps.([]interface{}) // map结构：maps.(map[string]interface{})
	if !ok {
		return 0, cherryError.Error("maps convert to map[string]interface{} error.")
	}

	loadMaps := make(map[int]*ShopBoxLvRow)
	for index, data := range list {
		loadConfig := &ShopBoxLvRow{}
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

func (c *shopBoxLvConfig) Name() string {
	return "ShopBoxLvConfig"
}

func (c *shopBoxLvConfig) OnAfterLoad(_ bool) {}

func (c *shopBoxLvConfig) Get(key int) (*ShopBoxLvRow, bool) {
	row, ok := c.maps[key]
	return row, ok
}

func (c *shopBoxLvConfig) List() []*ShopBoxLvRow {
	var list []*ShopBoxLvRow
	for _, row := range c.maps {
		list = append(list, row)
	}
	return list
}
