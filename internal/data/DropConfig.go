package data

import (
	cherryError "gameserver/cherry/error"
	cherryLogger "gameserver/cherry/logger"
)

var DropConfig = &dropConfig{}

type (
	DropRow struct {
		Id        int `json:"id"`        // 唯一ID
		DropGroup int `json:"dropGroup"` // 掉落组
		HeroId    int `json:"heroId"`    // 英雄
		ItemId    int `json:"itemId"`    // 道具
		Coin      int `json:"coin"`      // 金币
		Weight    int `json:"weight"`    // 权重
	}

	dropConfig struct {
		maps   map[int]*DropRow
		groups map[int][]*DropRow
	}
)

func (c *dropConfig) Init() {
	c.maps = make(map[int]*DropRow)
}

func (c *dropConfig) OnLoad(maps interface{}, _ bool) (int, error) {
	list, ok := maps.([]interface{}) // map结构：maps.(map[string]interface{})
	if !ok {
		return 0, cherryError.Error("maps convert to map[string]interface{} error.")
	}

	loadMaps := make(map[int]*DropRow)
	groups := make(map[int][]*DropRow)
	for index, data := range list {
		loadConfig := &DropRow{}
		err := DecodeData(data, loadConfig)
		if err != nil {
			cherryLogger.Warnf("decode error. [id = %v, %v], err = %s", index, loadConfig, err)
			continue
		}
		loadMaps[loadConfig.Id] = loadConfig
		groups[loadConfig.DropGroup] = append(groups[loadConfig.DropGroup], loadConfig)
	}
	c.maps = loadMaps
	c.groups = groups

	return len(list), nil
}

func (c *dropConfig) Name() string {
	return "DropConfig"
}

func (c *dropConfig) OnAfterLoad(_ bool) {}

func (c *dropConfig) Get(key int) (*DropRow, bool) {
	row, ok := c.maps[key]
	return row, ok
}

func (c *dropConfig) GetByGroupId(key int) ([]*DropRow, bool) {
	row, ok := c.groups[key]
	return row, ok
}

func (c *dropConfig) List() []*DropRow {
	var list []*DropRow
	for _, row := range c.maps {
		list = append(list, row)
	}
	return list
}
