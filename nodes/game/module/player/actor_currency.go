package player

import (
	clog "gameserver/cherry/logger"
	cproto "gameserver/cherry/net/proto"
	"gameserver/internal/pb"
	"gameserver/nodes/game/db"
)

type (
	ActorCurrency struct {
		*ActorPlayer
	}
)

func NewActorCurrency(player *ActorPlayer) *ActorCurrency {
	return &ActorCurrency{player}
}

func (i *ActorCurrency) pushInfo(session *cproto.Session) {
	i.Push(session, "currencyInfo", i.buildPbInfo(session.Uid))
}

func (i *ActorCurrency) getInfo(session *cproto.Session, _ *pb.None) {
	i.Response(session, i.buildPbInfo(session.Uid))
}

func (i *ActorCurrency) buildPbInfo(playerId int64) *pb.S2CCurrencyInfo {
	tb := db.CurrencyRepository.GetOrCreate(playerId)
	tbMap := tb.GetMaps()
	info := &pb.S2CCurrencyInfo{
		Currencys: make(map[int32]int64),
	}
	for id, count := range tbMap {
		info.Currencys[int32(id)] = int64(count)
	}
	clog.Infof("currency info:%v", info)
	return info
}
