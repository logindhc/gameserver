package data

import (
	cherryError "gameserver/cherry/error"
	cherryLogger "gameserver/cherry/logger"
)

var ResTypeConfig = &resTypeConfig{}

type (
	ResTypeRow struct {
		Id  int    `json:"id"`  // 资源类型
		Des string `json:"des"` // 说明
	}

	resTypeConfig struct {
		maps map[int]*ResTypeRow
	}
)

func (c *resTypeConfig) Init() {
	c.maps = make(map[int]*ResTypeRow)
}

func (c *resTypeConfig) OnLoad(maps interface{}, _ bool) (int, error) {
	list, ok := maps.([]interface{}) // map结构：maps.(map[string]interface{})
	if !ok {
		return 0, cherryError.Error("maps convert to map[string]interface{} error.")
	}

	loadMaps := make(map[int]*ResTypeRow)
	for index, data := range list {
		loadConfig := &ResTypeRow{}
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

func (c *resTypeConfig) Name() string {
	return "ResTypeConfig"
}

func (c *resTypeConfig) OnAfterLoad(_ bool) {}

func (c *resTypeConfig) Get(key int) (*ResTypeRow, bool) {
	row, ok := c.maps[key]
	return row, ok
}

func (c *resTypeConfig) List() []*ResTypeRow {
	var list []*ResTypeRow
	for _, row := range c.maps {
		list = append(list, row)
	}
	return list
}
