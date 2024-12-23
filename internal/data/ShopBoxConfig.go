package data

import (
	cherryError "gameserver/cherry/error"
	cherryLogger "gameserver/cherry/logger"
)

var ShopBoxConfig = &shopBoxConfig{}

type (
	ShopBoxRow struct {
		Id         int     `json:"id"`         // 唯一ID
		BoxName    string  `json:"boxName"`    // 宝箱名字
		Drops      [][]int `json:"drops"`      // [掉落组ID,掉落最小次数，掉落最大次数]
		RewardGold int     `json:"rewardGold"` // 获得金币
		RewardExp  int     `json:"rewardExp"`  // 获得宝箱经验
	}

	shopBoxConfig struct {
		maps map[int]*ShopBoxRow
	}
)

func (c *shopBoxConfig) Init() {
	c.maps = make(map[int]*ShopBoxRow)
}

func (c *shopBoxConfig) OnLoad(maps interface{}, _ bool) (int, error) {
	list, ok := maps.([]interface{}) // map结构：maps.(map[string]interface{})
	if !ok {
		return 0, cherryError.Error("maps convert to map[string]interface{} error.")
	}

	loadMaps := make(map[int]*ShopBoxRow)
	for index, data := range list {
		loadConfig := &ShopBoxRow{}
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

func (c *shopBoxConfig) Name() string {
	return "ShopBoxConfig"
}

func (c *shopBoxConfig) OnAfterLoad(_ bool) {}

func (c *shopBoxConfig) Get(key int) (*ShopBoxRow, bool) {
	row, ok := c.maps[key]
	return row, ok
}

func (c *shopBoxConfig) List() []*ShopBoxRow {
	var list []*ShopBoxRow
	for _, row := range c.maps {
		list = append(list, row)
	}
	return list
}
