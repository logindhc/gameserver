package player

import (
	clog "gameserver/cherry/logger"
	cproto "gameserver/cherry/net/proto"
	"gameserver/internal/pb"
)

type (
	Item struct {
		*ActorPlayer
	}
)

func (i *Item) OnInit() {
	clog.Debugf("[Item] path = %s init!", i.PathString())
	i.Local().Register("getItemInfo", i.getItemInfo) // 注册 查看角色
}

func (i *Item) getItemInfo(session *cproto.Session, _ *pb.None) {
	playerId := session.Uid

	clog.Debugf("[Item] getInfo playerId = %d", playerId)
}
