package data

import (
	cherryError "gameserver/cherry/error"
	cherryLogger "gameserver/cherry/logger"
)

var GameConfig = &gameConfig{}

type (
	GameRow struct {
		Id      int            `json:"id"`      // ID（不可修改）
		Dec     string         `json:"dec"`     // 说明
		Val     int            `json:"val"`     // 参数
		ArrVal  []int          `json:"arrVal"`  // 数组
		Arr2Val [][]int        `json:"arr2Val"` // 二维数组
		MapVal  map[string]int `json:"mapVal"`  // Map结构kv
	}

	gameConfig struct {
		maps map[int]*GameRow
	}
)

func (c *gameConfig) Init() {
	c.maps = make(map[int]*GameRow)
}

func (c *gameConfig) OnLoad(maps interface{}, _ bool) (int, error) {
	list, ok := maps.([]interface{}) // map结构：maps.(map[string]interface{})
	if !ok {
		return 0, cherryError.Error("maps convert to map[string]interface{} error.")
	}

	loadMaps := make(map[int]*GameRow)
	for index, data := range list {
		loadConfig := &GameRow{}
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

func (c *gameConfig) Name() string {
	return "GameConfig"
}

func (c *gameConfig) OnAfterLoad(_ bool) {}

func (c *gameConfig) Get(key int) (*GameRow, bool) {
	row, ok := c.maps[key]
	return row, ok
}

func (c *gameConfig) List() []*GameRow {
	var list []*GameRow
	for _, row := range c.maps {
		list = append(list, row)
	}
	return list
}
