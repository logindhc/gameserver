package db

import (
	cherryTime "gameserver/cherry/extend/time"
	clog "gameserver/cherry/logger"
	cproto "gameserver/cherry/net/proto"
	"gameserver/internal/code"
	"gameserver/internal/data"
	"gameserver/internal/persistence"
	"gameserver/internal/persistence/repository"
	sessionKey "gameserver/internal/session_key"
)

// PlayerTable 角色基础表
type PlayerTable struct {
	ID             int64  `gorm:"primaryKey;autoIncrement:false;column:id;comment:id" json:"id"`
	PID            int32  `gorm:"column:pid;comment:平台id" json:"pid"`
	OpenId         string `gorm:"column:open_id;comment:平台open_id" json:"openId"`
	ServerId       int32  `gorm:"column:server_id;comment:创角时的游戏服id" json:"serverId"`
	Name           string `gorm:"column:player_name;comment:角色名称" json:"name"`
	Gender         int32  `gorm:"column:gender;comment:角色性别" json:"gender"`
	Level          int32  `gorm:"column:level;comment:角色等级" json:"level"`
	Exp            int64  `gorm:"column:exp;comment:角色经验" json:"exp"`
	CreateTime     int64  `gorm:"column:create_time;comment:创建时间" json:"createTime"`
	LastLoginTime  int64  `gorm:"column:last_login_time;comment:最后登录时间" json:"lastLoginTime"`
	LastLogoutTime int64  `gorm:"column:last_logout_time;comment:最后登出时间" json:"LastLogoutTime"`
}

func (*PlayerTable) TableName() string {
	return "player"
}

var PlayerRepository *repository.LruRepository[int64, PlayerTable]

func (p *PlayerTable) InitRepository() {
	PlayerRepository = repository.NewLruRepository[int64, PlayerTable](database.GetGameDB(), p.TableName())
	persistence.RegisterRepository(PlayerRepository)
}

func CreatePlayer(session *cproto.Session, name string, serverId int32, playerInit *data.PlayerInitRow) (*PlayerTable, int32) {
	pid := session.GetInt32(sessionKey.PID)
	openId := session.GetString(sessionKey.OpenID)

	if session.Uid < 1 || pid < 1 || openId == "" {
		clog.Warnf("create playerTable fail. pid or openId is error. [name = %s, pid = %v, openId = %v]",
			name,
			pid,
			openId,
		)
		return nil, code.PlayerCreateFail
	}
	playerTable := &PlayerTable{
		ID:         session.Uid,
		PID:        pid,
		OpenId:     openId,
		ServerId:   serverId,
		Name:       name,
		Gender:     playerInit.Gender,
		Level:      playerInit.Level,
		Exp:        0,
		CreateTime: cherryTime.Now().ToMillisecond(),
	}

	add := PlayerRepository.Add(playerTable)
	if add == nil {
		return nil, code.PlayerCreateFail
	}
	// TODO 初始化角色相关的表
	// 道具表
	// 英雄表

	return playerTable, code.OK
}
