package resmgr

import (
	"errors"
	"fmt"
	cstring "gameserver/cherry/extend/string"
	clog "gameserver/cherry/logger"
	"gameserver/internal/data"
	"gameserver/nodes/game/facade"
	"gameserver/nodes/game/opcode"
	"gameserver/nodes/game/res/restype"
)

// 产出类型
const (
	add = iota + 1 //--产出
	del            //消耗
)

type resManager struct {
	resMap map[int]facade.IRes
}

func (r *resManager) Register(res facade.IRes) {
	if _, ok := r.resMap[res.BaseID()]; ok {
		panic("重复注册资源")
	}
	r.resMap[res.BaseID()] = res
	clog.Infow("res register", "baseId", res.BaseID(), "res", res)
}
func (r *resManager) EnoughOne(playerId int64, baseId, id, count int, dup ...int) bool {
	ress, err := r.format(baseId, id, count)
	if err != nil {
		clog.Debugf("res enough err:%v", err)
		return false
	}
	return r.Enough(playerId, ress, dup...)
}
func (r *resManager) Enough(playerId int64, ress [][]int, dup ...int) bool {
	ress = r.merge(ress)
	for _, res := range ress {
		baseId, id, count, err := r.formatArgs(res, dup...)
		if err != nil {
			clog.Debugf("res enough err:%v", err)
			return false
		}
		if !r.resMap[baseId].Enough(playerId, id, count) {
			return false
		}
	}
	return true
}
func (r *resManager) AddOne(playerId int64, baseId, id, count int, opcode opcode.Type, remark string, dup ...int) ([][]int, error) {
	ress, err := r.format(baseId, id, count)
	if err != nil {
		return nil, err
	}
	return r.Add(playerId, ress, opcode, remark, dup...)
}
func (r *resManager) Add(playerId int64, ress [][]int, opcode opcode.Type, remark string, dup ...int) ([][]int, error) {
	return r.change(playerId, add, ress, opcode, remark, dup...)
}
func (r *resManager) DelOne(playerId int64, baseId, id, count int, opcode opcode.Type, remark string, dup ...int) ([][]int, error) {
	ress, err := r.format(baseId, id, count)
	if err != nil {
		clog.Debugf("del err:%v", err)
		return nil, err
	}
	return r.Del(playerId, ress, opcode, remark, dup...)
}
func (r *resManager) Del(playerId int64, ress [][]int, opcode opcode.Type, remark string, dup ...int) ([][]int, error) {
	return r.change(playerId, del, ress, opcode, remark, dup...)
}

func (r *resManager) GoldEnough(playerId int64, count int) bool {
	res, err := r.formatCurrency(restype.Gold, count)
	if err != nil {
		return false
	}
	return r.Enough(playerId, res)
}
func (r *resManager) GoldAdd(playerId int64, count int, opcode opcode.Type, remark string) ([][]int, error) {
	res, err := r.formatCurrency(restype.Gold, count)
	if err != nil {
		return nil, err
	}
	return r.Add(playerId, res, opcode, remark)
}
func (r *resManager) GoldDel(playerId int64, count int, opcode opcode.Type, remark string) ([][]int, error) {
	res, err := r.formatCurrency(restype.Gold, count)
	if err != nil {
		return nil, err
	}
	return r.Del(playerId, res, opcode, remark)
}

func (r *resManager) MoneyEnough(playerId int64, count int) bool {
	res, err := r.formatCurrency(restype.Money, count)
	if err != nil {
		return false
	}
	return r.Enough(playerId, res)
}
func (r *resManager) MoneyAdd(playerId int64, count int, opcode opcode.Type, remark string) ([][]int, error) {
	res, err := r.formatCurrency(restype.Money, count)
	if err != nil {
		return nil, err
	}
	return r.Add(playerId, res, opcode, remark)
}
func (r *resManager) MoneyDel(playerId int64, count int, opcode opcode.Type, remark string) ([][]int, error) {
	res, err := r.formatCurrency(restype.Money, count)
	if err != nil {
		return nil, err
	}
	return r.Del(playerId, res, opcode, remark)
}

func (r *resManager) DiamondEnough(playerId int64, count int) bool {
	res, err := r.formatCurrency(restype.Diamond, count)
	if err != nil {
		return false
	}
	return r.Enough(playerId, res)
}
func (r *resManager) DiamondAdd(playerId int64, count int, opcode opcode.Type, remark string) ([][]int, error) {
	res, err := r.formatCurrency(restype.Diamond, count)
	if err != nil {
		return nil, err
	}
	return r.Add(playerId, res, opcode, remark)
}
func (r *resManager) DiamondDel(playerId int64, count int, opcode opcode.Type, remark string) ([][]int, error) {
	res, err := r.formatCurrency(restype.Diamond, count)
	if err != nil {
		return nil, err
	}
	return r.Del(playerId, res, opcode, remark)
}

func (r *resManager) merge(ress [][]int) [][]int {
	//合并相同id的资源[[1,2,2],[1,2,3]] [[1,2,5]]
	mergedMap := make(map[[2]int]int)
	for _, arr := range ress {
		if len(arr) != 3 {
			continue // 忽略不符合要求的子数组
		}
		if r.resMap[arr[0]].BaseID() != arr[0] {
			clog.Errorf("资源ID错误:%d", arr[0])
			continue
		}
		key := [2]int{arr[0], arr[1]}
		mergedMap[key] += arr[2]
	}
	var result [][]int
	for key, value := range mergedMap {
		result = append(result, []int{key[0], key[1], value})
	}
	return result
}

func (r *resManager) change(playerId int64, optype int, ress [][]int, opcode opcode.Type, remark string, dup ...int) ([][]int, error) {
	if cstring.IsBlank(remark) || opcode <= 0 {
		return nil, errors.New("remark empty || opcode empty")
	}
	if len(ress) == 0 {
		return nil, errors.New("res empty")
	}
	retRes := make([][]int, 0, len(ress))
	for _, res := range ress {
		baseId, id, count, err := r.formatArgs(res, dup...)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			//跳过数量0的逻辑，调用的地方就不需要每次都判断数据为0的情况
			continue
		}
		if optype == del {
			count = -count
		}

		now, convert := r.resMap[baseId].Change(playerId, id, count)
		if (baseId == restype.Currency) || (baseId == restype.Item && now > 0) || (baseId == restype.Hero && now > 0) {
			//宝箱，英雄重复获得转换碎片了，就不用通知前端变更了
			retRes = append(retRes, []int{baseId, id, now})
		}
		//转换资源
		if count > 0 && len(convert) > 0 {
			for _, c := range convert {
				now, _ = r.resMap[c[0]].Change(playerId, c[1], c[2])
				//不使用递归，只转换一次,避免死循环
				retRes = append(retRes, []int{c[0], c[1], now})
			}
		}
	}
	clog.Infow("change res", "opcode", opcode, "remark", remark, "res", ress)
	return retRes, nil
}

func (r *resManager) formatArgs(res []int, dup ...int) (baseId, id, count int, err error) {
	var d = 0
	if len(dup) == 0 {
		d = 1
	}
	if len(res) != 3 {
		return 0, 0, 0, errors.New("res formatArgs len != 3")
	}

	baseId = res[0]
	if _, ok := data.ResTypeConfig.Get(baseId); !ok { //验证资源类型
		return 0, 0, 0, errors.New(fmt.Sprintf("res formatArgs baseId not exist [%d]", baseId))
	}
	id = res[1]
	if !r.resMap[baseId].CheckConfigId(id) { //验证资源id
		return 0, 0, 0, errors.New(fmt.Sprintf("res formatArgs id not exist [%d][%d]", baseId, id))
	}
	count = res[2] * d
	err = nil
	if count < 0 {
		err = errors.New("res formatArgs count < 0")
	}
	return
}

func (r *resManager) formatCurrency(id restype.CurrencyType, count int) (res [][]int, err error) {
	res = [][]int{{restype.Currency, int(id), count}}
	return res, nil
}

func (r *resManager) format(baseId, id, count int) (res [][]int, err error) {
	if _, ok := data.ResTypeConfig.Get(baseId); !ok {
		return nil, errors.New(fmt.Sprintf("res format baseId not exist [%d]", baseId))
	}
	if !r.resMap[baseId].CheckConfigId(id) {
		return nil, errors.New(fmt.Sprintf("res format id not exist [%d][%d]", baseId, id))
	}
	res = [][]int{{baseId, id, count}}
	return res, nil
}
