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
	ActorItem struct {
		*ActorPlayer
	}
)

func NewActorItem(player *ActorPlayer) *ActorItem {
	return &ActorItem{player}
}
func (i *ActorItem) pushInfo(session *cproto.Session) {
	i.Push(session, "itemInfo", i.buildPbInfo(session.Uid))
}

func (i *ActorItem) getInfo(session *cproto.Session, _ *pb.None) {
	i.Response(session, i.buildPbInfo(session.Uid))
}

// use 玩家使用道具
func (i *ActorItem) use(session *cproto.Session, req *pb.C2SItemUse) {
	itemId := req.ItemId
	itemRow, b := data.ItemConfig.Get(int(itemId))
	if !b {
		i.ResponseCode(session, code.ItemNotEnough)
		return
	}
	if !itemRow.ItemUse {
		i.ResponseCode(session, code.ItemNotAvailable)
		return
	}
	count := req.Count
	enough := resmgr.Instance.EnoughOne(session.Uid, restype.Item, int(itemId), int(count))
	if !enough {
		i.ResponseCode(session, code.ItemNotEnough)
		return
	}
	ress, err := resmgr.Instance.DelOne(session.Uid, restype.Item, int(itemId), int(count), opcode.ItemUse, "item_use")
	if err != nil {
		i.ResponseCode(session, code.Error)
		return
	}
	clog.Debugf("[%d] [useItem] itemId = %d, count = %d", session.Uid, itemId, count)
	response := &pb.S2CItemUse{}
	i.Response(session, response)
	i.PushResUpdateInfo(session, ress)
}

func (i *ActorItem) buildPbInfo(playerId int64) *pb.S2CItemInfo {
	tb := db.ItemRepository.GetOrCreate(playerId)
	tbMap := tb.GetItems()
	info := &pb.S2CItemInfo{
		Items: make(map[int32]int64),
	}
	for id, count := range tbMap {
		info.Items[int32(id)] = int64(count)
	}
	return info
}
