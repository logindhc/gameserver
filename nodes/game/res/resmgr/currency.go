package resmgr

import (
	clog "gameserver/cherry/logger"
	"gameserver/internal/data"
	"gameserver/nodes/game/db"
	"gameserver/nodes/game/res/restype"
	"math"
)

type currencyRes struct {
}

func (i *currencyRes) BaseID() int {
	return restype.Currency
}

func (i *currencyRes) CheckConfigId(id int) bool {
	_, b := data.CurrencyConfig.Get(id)
	return b
}
func (i *currencyRes) Get(playerId int64, id int) int {
	res := db.CurrencyRepository.GetOrCreate(playerId)
	resMap := res.GetMaps()
	if resMap == nil {
		return 0
	}
	count, ok := resMap[id]
	if !ok {
		return 0
	}
	return count
}

func (i *currencyRes) Enough(playerId int64, id, count int) bool {
	return i.Get(playerId, id) >= count
}

func (i *currencyRes) Change(playerId int64, id int, count int) (int, [][]int) {
	now := i.Get(playerId, id) + count
	if now < 0 {
		now = 0
	}
	if now > math.MaxInt {
		now = math.MaxInt
		clog.Errorf("currency count overflow:%v,%v,%v", playerId, id, count)
	}
	tb := db.CurrencyRepository.GetOrCreate(playerId)
	tb.GetMaps()[id] = now
	db.CurrencyRepository.Update(tb)
	clog.Infof("change currency:%v,%v,%v,%v", playerId, id, count, now)
	return now, nil
}
