package dropmgr

import (
	"errors"
	"gameserver/internal/data"
	"gameserver/nodes/game/res/restype"
	"math/rand"
	"time"
)

var (
	Instance = &dropMgr{}
)

type dropMgr struct {
}

func (d *dropMgr) Drop(playerId int64, dropGroupId int, count int) ([][]int, error) {
	// 使用新的随机种子生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return d.drop(playerId, dropGroupId, count, r)
}
func (d *dropMgr) DropRand(playerId int64, dropGroupId int, min, max int) ([][]int, error) {
	// 使用新的随机种子生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	count := r.Intn(max-min+1) + min
	return d.drop(playerId, dropGroupId, count, r)
}

func (d *dropMgr) drop(playerId int64, dropGroupId int, count int, r *rand.Rand) ([][]int, error) {
	row, b := data.DropConfig.GetByGroupId(dropGroupId)
	if !b {
		return nil, nil
	}
	totalWeight := 0
	for _, drop := range row {
		totalWeight += drop.Weight
	}
	if totalWeight <= 0 {
		return nil, errors.New("totalWeight <= 0")
	}
	res := make([][]int, 0, count)
	for range count {
		randomValue := r.Intn(totalWeight)
		currentWeight := 0
		for _, drop := range row {
			currentWeight += drop.Weight
			if randomValue < currentWeight {
				heroId := drop.HeroId
				if heroId > 0 {
					// 掉落英雄
					res = append(res, []int{restype.Hero, heroId, 1})
				}
				itemId := drop.ItemId
				if itemId > 0 {
					// 掉落物品
					res = append(res, []int{restype.Item, itemId, 1})
				}
				gold := drop.Coin
				if gold > 0 {
					// 掉落金币
					res = append(res, []int{restype.Currency, int(restype.Gold), gold})
				}
				break
			}
		}
	}
	return res, nil
}
