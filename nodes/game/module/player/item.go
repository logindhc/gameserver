package player

import (
	clog "gameserver/cherry/logger"
	cproto "gameserver/cherry/net/proto"
	"gameserver/internal/pb"
	"gameserver/nodes/game/db"
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
	response := &pb.PlayerSelectResponse{}
	// 游戏设定单服单角色，协议设计成可返回多角色
	playerTable := db.PlayerRepository.Get(session.Uid)
	if playerTable != nil {
		playerInfo := buildPBPlayer(playerTable)
		response.List = append(response.List, &playerInfo)
	}
	i.Response(session, response)
	clog.Debugf("[Item] getInfo playerId = %d", playerId)
}
