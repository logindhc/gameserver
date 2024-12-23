package player

import (
	clog "gameserver/cherry/logger"
	cproto "gameserver/cherry/net/proto"
	"gameserver/internal/code"
	"gameserver/internal/data"
	"gameserver/internal/pb"
	"gameserver/nodes/game/db"
	"gameserver/nodes/game/opcode"
	"gameserver/nodes/game/res/resmgr"
	"gameserver/nodes/game/res/restype"
)

type (
	ActorHero struct {
		*ActorPlayer
	}
)

func NewActorHero(player *ActorPlayer) *ActorHero {
	return &ActorHero{player}
}

func (i *ActorHero) pushInfo(session *cproto.Session) {
	i.Push(session, "heroInfo", i.buildPbInfo(session.Uid))
}

func (i *ActorHero) getInfo(session *cproto.Session, _ *pb.None) {
	i.Response(session, i.buildPbInfo(session.Uid))
}

// use 英雄升级
func (i *ActorHero) up(session *cproto.Session, req *pb.C2SHeroUp) {
	heroId := int(req.HeroId)
	if heroId <= 0 {
		i.ResponseCode(session, code.ParamError)
		return
	}
	heroRow, b := data.HeroConfig.Get(heroId)
	if !b {
		i.ResponseCode(session, code.ConfigError)
		return
	}

	tb := db.HeroRepository.GetOrCreate(session.Uid)
	//判断是否存在英雄,等级是否一致
	level, ok := tb.GetHeros()[heroRow.HeroGroup]
	if !ok {
		i.ResponseCode(session, code.HeroNotEnough)
		return
	}
	if level != heroRow.HeroLevel {
		i.ResponseCode(session, code.HeroLevelError)
		return
	}
	//下一级的配置
	nextLvHero, ok := data.HeroConfig.GetByGroupLevel(heroRow.HeroGroup, heroRow.HeroLevel+1)
	if !ok {
		i.ResponseCode(session, code.HeroMaxLevel)
		return
	}
	upCost := heroRow.UpCost
	enough := resmgr.Instance.GoldEnough(session.Uid, upCost)
	if !enough {
		i.ResponseCode(session, code.GoldNotEnough)
		return
	}

	itemId := heroRow.PieceId
	pieceCount := heroRow.PieceCount
	enough = resmgr.Instance.EnoughOne(session.Uid, restype.Item, itemId, pieceCount)
	if !enough {
		i.ResponseCode(session, code.ItemNotEnough)
		return
	}

	consumeRes := [][]int{
		{restype.Item, itemId, pieceCount},
		{restype.Currency, int(restype.Gold), upCost},
	}
	ress, err := resmgr.Instance.Del(session.Uid, consumeRes, opcode.HeroUp, "hero_up")
	if err != nil {
		i.ResponseCode(session, 1000)
		return
	}
	i.PushResUpdateInfo(session, ress)
	tb.GetHeros()[heroRow.HeroGroup] = nextLvHero.HeroLevel
	db.HeroRepository.Update(tb)

	clog.Debugf("%d hero up %d->%d", session.Uid, heroId, nextLvHero.Id)
	response := &pb.S2CHeroUp{}
	response.DelHeroId = req.HeroId
	response.HeroId = int32(nextLvHero.Id)
	i.Response(session, response)

}

func (i *ActorHero) buildPbInfo(playerId int64) *pb.S2CHeroInfo {
	tb := db.HeroRepository.GetOrCreate(playerId)
	tbMap := tb.GetHeros()
	heroInfo := &pb.S2CHeroInfo{
		Heros: make([]int32, 0, len(tbMap)),
	}
	for id, lv := range tbMap {
		heroRow, ok := data.HeroConfig.GetByGroupLevel(id, lv)
		if !ok {
			continue
		}
		heroInfo.Heros = append(heroInfo.Heros, int32(heroRow.Id))
	}
	return heroInfo
}
