package resmgr

import (
	clog "gameserver/cherry/logger"
	"gameserver/internal/data"
	"gameserver/nodes/game/db"
	dropmgr "gameserver/nodes/game/drop"
	"gameserver/nodes/game/res/restype"
)

type boxRes struct {
}

func (i *boxRes) BaseID() int {
	return restype.Box
}
func (i *boxRes) CheckConfigId(id int) bool {
	_, b := data.ShopBoxConfig.Get(id)
	return b
}
func (i *boxRes) Get(playerId int64, id int) int {
	return 0
}

func (i *boxRes) Enough(playerId int64, id, count int) bool {
	return i.Get(playerId, id) >= count
}

func (i *boxRes) Change(playerId int64, id int, count int) (int, [][]int) {
	if count <= 0 {
		//获取宝箱直接调用掉落逻辑
		return 0, nil
	}
	boxRow, b := data.ShopBoxConfig.Get(id)
	if !b {
		clog.Warnf("box config not exist:%v", id)
		return 0, nil
	}
	drops := boxRow.Drops
	ress := make([][]int, 0)
	//随机掉落奖励
	for _, drop := range drops {
		if len(drop) != 3 {
			continue
		}
		dropId := drop[0]
		minCount := drop[1]
		maxCount := drop[2]
		randRes, err := dropmgr.Instance.DropRand(playerId, dropId, minCount, maxCount)
		if err != nil {
			return 0, nil
		}
		ress = append(ress, randRes...)
	}
	// 获得宝箱固定获得金币
	if boxRow.RewardGold > 0 {
		ress = append(ress, []int{restype.Currency, int(restype.Gold), boxRow.RewardGold})
	}
	// 获得宝箱固定获得宝箱经验，会触发宝箱升级
	if boxRow.RewardExp > 0 {
		shopTable := db.ShopRepository.GetOrCreate(playerId)
		boxLv := shopTable.BoxLevel
		shopTable.BoxExp += int32(boxRow.RewardExp)
		if shopBoxLvRow, ok := data.ShopBoxLvConfig.Get(int(boxLv)); ok {
			if diff := shopTable.BoxExp - int32(shopBoxLvRow.Exp); diff >= 0 { //宝箱等级升级了
				if nextConfig, upLvOk := data.ShopBoxLvConfig.Get(int(shopTable.BoxLevel + 1)); upLvOk {
					shopTable.BoxLevel = int32(nextConfig.Id)
					shopTable.BoxExp = diff
					clog.Infof("%d box lv up:%v exp:%v", playerId, shopTable.BoxLevel, shopTable.BoxExp)
				}
			}
		}
		db.ShopRepository.Update(shopTable)
	}
	clog.Infof("open box:%v,%v", id, count)
	return 0, ress
}
