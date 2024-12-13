package player

import (
	clog "gameserver/cherry/logger"
	cproto "gameserver/cherry/net/proto"
	"gameserver/internal/pb"
)

type (
	ActorItem struct {
		*ActorPlayer
	}
)

func NewActorItem(player *ActorPlayer) *ActorItem {
	return &ActorItem{player}
}

func (p *ActorItem) getInfo(session *cproto.Session, _ *pb.None) {
	response := &pb.S2CItemInfo{
		Items: map[int32]int64{
			1: 1,
		},
	}
	p.Response(session, response)
}

// use 玩家使用道具
func (p *ActorItem) use(session *cproto.Session, req *pb.C2SItemUse) {
	itemId := req.ItemId
	count := req.Count
	clog.Debugf("[%d] [useItem] itemId = %d, count = %d", p.Id, itemId, count)
	response := &pb.S2CItemUse{}
	p.Response(session, response)
}
