package resmgr

import (
	clog "gameserver/cherry/logger"
	"gameserver/internal/data"
	"gameserver/nodes/game/db"
	"gameserver/nodes/game/res/restype"
	"math"
)

type itemRes struct {
}

func (i *itemRes) BaseID() int {
	return restype.Item
}
func (i *itemRes) CheckConfigId(id int) bool {
	_, b := data.ItemConfig.Get(id)
	return b
}
func (i *itemRes) Get(playerId int64, id int) int {
	item := db.ItemRepository.GetOrCreate(playerId)
	items := item.GetItems()
	if items == nil {
		return 0
	}
	count, ok := items[id]
	if !ok {
		return 0
	}
	return count
}

func (i *itemRes) Enough(playerId int64, id, count int) bool {
	return i.Get(playerId, id) >= count
}

func (i *itemRes) Change(playerId int64, id int, count int) (int, [][]int) {
	now := i.Get(playerId, id) + count
	if now < 0 {
		now = 0
	}
	if now > math.MaxInt {
		now = math.MaxInt
		clog.Errorf("item count overflow:%v,%v,%v", playerId, id, count)
	}
	item := db.ItemRepository.GetOrCreate(playerId)
	item.GetItems()[id] = now
	db.ItemRepository.Update(item)
	clog.Infof("change item:%v,%v,%v", id, count, now)
	return now, nil
}
