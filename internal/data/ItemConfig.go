package data

import (
	cherryError "gameserver/cherry/error"
	cherryLogger "gameserver/cherry/logger"
)

var ItemConfig = &itemConfig{}

type (
	ItemRow struct {
		Id          int     `json:"id"`          // 唯一ID
		ItemType    int     `json:"itemType"`    // 类型（1碎片、2宝箱）
		ItemQuality int     `json:"itemQuality"` // 品质
		ItemUse     bool    `json:"itemUse"`     // 是否可使用
		ItemBundle  bool    `json:"itemBundle"`  // 是否宝箱
		Weight      [][]int `json:"weight"`      // 权重[权重,最小数量，最大数量]
		UseParm     [][]int `json:"useParm"`     // 使用效果
	}

	itemConfig struct {
		maps map[int]*ItemRow
	}
)

func (c *itemConfig) Init() {
	c.maps = make(map[int]*ItemRow)
}

func (c *itemConfig) OnLoad(maps interface{}, _ bool) (int, error) {
	list, ok := maps.([]interface{}) // map结构：maps.(map[string]interface{})
	if !ok {
		return 0, cherryError.Error("maps convert to map[string]interface{} error.")
	}

	loadMaps := make(map[int]*ItemRow)
	for index, data := range list {
		loadConfig := &ItemRow{}
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

func (c *itemConfig) Name() string {
	return "ItemConfig"
}

func (c *itemConfig) OnAfterLoad(_ bool) {}

func (c *itemConfig) Get(key int) (*ItemRow, bool) {
	row, ok := c.maps[key]
	return row, ok
}

func (c *itemConfig) List() []*ItemRow {
	var list []*ItemRow
	for _, row := range c.maps {
		list = append(list, row)
	}
	return list
}
