package data

import (
	cherryError "gameserver/cherry/error"
	cherryLogger "gameserver/cherry/logger"
)

var LevelRewardConfig = &levelRewardConfig{}

type (
	LevelRewardRow struct {
		Id          int            `json:"id"`          // 生存天数
		Reward      []int          `json:"reward"`      // 天数奖励
		LostReward  map[string]int `json:"lostReward"`  // 失败奖励(掉落)
		WinerReward map[string]int `json:"winerReward"` // 胜利奖励(掉落)
	}

	levelRewardConfig struct {
		maps map[int]*LevelRewardRow
	}
)

func (c *levelRewardConfig) Init() {
	c.maps = make(map[int]*LevelRewardRow)
}

func (c *levelRewardConfig) OnLoad(maps interface{}, _ bool) (int, error) {
	list, ok := maps.([]interface{}) // map结构：maps.(map[string]interface{})
	if !ok {
		return 0, cherryError.Error("maps convert to map[string]interface{} error.")
	}

	loadMaps := make(map[int]*LevelRewardRow)
	for index, data := range list {
		loadConfig := &LevelRewardRow{}
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

func (c *levelRewardConfig) Name() string {
	return "LevelRewardConfig"
}

func (c *levelRewardConfig) OnAfterLoad(_ bool) {}

func (c *levelRewardConfig) Get(key int) (*LevelRewardRow, bool) {
	row, ok := c.maps[key]
	return row, ok
}

func (c *levelRewardConfig) List() []*LevelRewardRow {
	var list []*LevelRewardRow
	for _, row := range c.maps {
		list = append(list, row)
	}
	return list
}
