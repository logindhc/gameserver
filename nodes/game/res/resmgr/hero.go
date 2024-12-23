package resmgr

import (
	clog "gameserver/cherry/logger"
	"gameserver/internal/data"
	"gameserver/nodes/game/db"
	"gameserver/nodes/game/res/restype"
)

type heroRes struct {
}

func (i *heroRes) BaseID() int {
	return restype.Hero
}

func (i *heroRes) CheckConfigId(id int) bool {
	_, b := data.HeroConfig.Get(id)
	return b
}
func (i *heroRes) Get(playerId int64, id int) int {
	Hero := db.HeroRepository.GetOrCreate(playerId)
	heroRow, _ := data.HeroConfig.Get(id)
	heroGroupId := heroRow.HeroGroup
	if _, ok := Hero.GetHeros()[heroGroupId]; !ok {
		return 0
	}
	return 1
}

func (i *heroRes) Enough(playerId int64, id, count int) bool {
	return i.Get(playerId, id) >= count
}

func (i *heroRes) Change(playerId int64, id int, count int) (int, [][]int) {
	tb := db.HeroRepository.GetOrCreate(playerId)
	heroRow, _ := data.HeroConfig.Get(id)
	heroGroupId := heroRow.HeroGroup
	isNew := 0
	convertorCount := 0
	if count > 0 { // 增加
		if _, ok := tb.GetHeros()[heroGroupId]; !ok {
			// 创建英雄，根据获取id来控制
			tb.GetHeros()[heroGroupId] = heroRow.HeroLevel
			isNew = 1
			count--
		}
		if count > 0 {
			//转换成碎片
			convertorCount = heroRow.ConveterCount * count
		}
	} else { // 减少
		//TODO 英雄一般不操作删除，只有GM才删除
		clog.Infof("%d hero delete %d", playerId, id)
	}
	if isNew == 1 {
		db.HeroRepository.Update(tb)
	}
	var convertor [][]int = nil
	if convertorCount > 0 {
		convertor = [][]int{{restype.Item, heroRow.PieceId, convertorCount}}
	}
	clog.Infof("%d change hero:%v,%v,%v", playerId, id, count, convertorCount)
	return isNew, convertor
}
