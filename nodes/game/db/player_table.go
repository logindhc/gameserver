package db

import (
	cherryTime "gameserver/cherry/extend/time"
	clog "gameserver/cherry/logger"
	cproto "gameserver/cherry/net/proto"
	"gameserver/internal/code"
	"gameserver/internal/persistence"
	"gameserver/internal/persistence/repository"
	sessionKey "gameserver/internal/session_key"
)

// PlayerTable 角色基础表
type PlayerTable struct {
	ID             int64  `gorm:"primaryKey;autoIncrement:false;column:id;comment:id" json:"id"`
	Channel        int32  `gorm:"index;column:channel;comment:渠道id" json:"channel"`
	Platform       int32  `gorm:"index;column:platform;comment:平台id" json:"platform"`
	AccountId      int64  `gorm:"index;column:account_id;comment:游戏账号" json:"accountId"`
	OpenId         string `gorm:"index;column:open_id;comment:平台open_id" json:"openId"`
	ServerId       int32  `gorm:"column:server_id;comment:创角时的游戏服id" json:"serverId"`
	Nickname       string `gorm:"column:nickname;comment:角色名称" json:"nickname"`
	Gender         int32  `gorm:"column:gender;comment:角色性别" json:"gender"`
	CreateTime     int64  `gorm:"column:create_time;comment:创建时间" json:"createTime"`
	LastLoginTime  int64  `gorm:"column:last_login_time;comment:最后登录时间" json:"lastLoginTime"`
	LastLogoutTime int64  `gorm:"column:last_logout_time;comment:最后登出时间" json:"LastLogoutTime"`
	White          int32  `gorm:"column:white;comment:白名单" json:"white"`
	GmLevel        int32  `gorm:"column:gm_level;comment:gm等级" json:"gmLevel"`
	Forbidden      int32  `gorm:"column:forbidden;comment:封禁" json:"forbidden"`
	Country        int32  `gorm:"column:country;comment:国家" json:"country"`
}

func (*PlayerTable) TableName() string {
	return "player"
}

var PlayerRepository repository.IRepository[int64, PlayerTable]

func (p *PlayerTable) InitRepository() {
	PlayerRepository = repository.NewDefaultRepository[int64, PlayerTable](database.GetGameDB(), p.TableName())
	persistence.RegisterRepository(PlayerRepository)
}

func CreatePlayer(session *cproto.Session) (*PlayerTable, int32) {
	channel := session.GetInt32(sessionKey.ChannelID)
	platform := session.GetInt32(sessionKey.PlatformID)
	openId := session.GetString(sessionKey.OpenID)
	serverId := session.GetInt32(sessionKey.ServerID)
	accountId := session.GetInt64(sessionKey.AccountID)

	if session.Uid < 1 || channel < 1 || openId == "" {
		clog.Warnf("create playerTable fail. pid or openId is error. [name = %v, pid = %v, openId = %v]",
			session.Uid,
			channel,
			openId,
		)
		return nil, code.PlayerCreateFail
	}
	playerTable := &PlayerTable{
		ID:         session.Uid,
		Channel:    channel,
		Platform:   platform,
		OpenId:     openId,
		AccountId:  accountId,
		ServerId:   serverId,
		Nickname:   "",
		Gender:     0,
		CreateTime: cherryTime.Now().ToSecond(),
		White:      0,
		GmLevel:    0,
		Forbidden:  0,
		Country:    61845,
	}

	add := PlayerRepository.Add(playerTable)
	if add == nil {
		return nil, code.PlayerCreateFail
	}
	return playerTable, code.OK
}
