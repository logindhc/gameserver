package types

import (
	"encoding/json"
	cherryMapStructure "gameserver/cherry/extend/mapstructure"
	"github.com/spf13/cast"
	"reflect"
)

type IntMap map[int]int

func NewIntMap() IntMap {
	return make(map[int]int)
}

func (IntMap) Type() reflect.Type {
	return reflect.TypeOf(I32I64Map{})
}

func (p *IntMap) Hook() cherryMapStructure.DecodeHookFuncType {
	return func(_ reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if t == p.Type() {
			return p.ToMap(data), nil
		}
		return data, nil
	}
}

func (p *IntMap) ToMap(values interface{}) map[int]int {
	var maps = make(map[int]int)
	if values == nil {
		return maps
	}

	valueSlice := cast.ToSlice(values)
	if valueSlice == nil {
		return maps
	}

	if len(valueSlice) == 2 {
		k, kErr := cast.ToIntE(valueSlice[0])
		v, vErr := cast.ToIntE(valueSlice[1])
		if kErr == nil && vErr == nil {
			maps[k] = v
			return maps
		}
	}

	for _, value := range valueSlice {
		result, found := value.([]interface{})
		if found == false {
			break
		}

		if len(result) >= 2 {
			k := cast.ToInt(result[0])
			v := cast.ToInt(result[1])
			maps[k] = v
		}
	}

	return maps
}

func (p IntMap) ReadString(data string) {
	var jsonObject interface{}
	err := json.Unmarshal([]byte(data), &jsonObject)
	if err != nil {
		return
	}

	resultMap := p.ToMap(jsonObject)
	for k, v := range resultMap {
		p[k] = v
	}
}

func (p IntMap) Decrease(key int, decreaseValue int) (int, bool) {
	if decreaseValue < 1 {
		return 0, false
	}

	value, _ := p[key]
	if value < decreaseValue {
		return 0, false
	}

	p[key] = value - decreaseValue

	return p[key], true
}

func (p IntMap) Add(key int, addValue int) (int, bool) {
	if addValue < 1 {
		return 0, false
	}

	value, _ := p[key]
	p[key] = value + addValue

	return p[key], true
}

func (p IntMap) Get(key int) (int, bool) {
	value, found := p[key]
	if found {
		return value, true
	}
	return 0, false
}

func (p IntMap) Set(key int, value int) {
	p[key] = value
}

func (p IntMap) ContainKey(key int) bool {
	_, found := p[key]
	return found
}
