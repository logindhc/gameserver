package db

import (
	cherrySnowflake "gameserver/cherry/extend/snowflake"
	"gameserver/internal/persistence"
	"gameserver/internal/persistence/repository"
	"time"
)

var AccountRepository repository.IRepository[string, AccountTable]

func (p *AccountTable) InitRepository() {
	AccountRepository = repository.NewDefaultRepository[string, AccountTable](database.GetGameDB(), p.TableName())
	persistence.RegisterRepository(AccountRepository)
}

// AccountTable 帐号信息表
type AccountTable struct {
	ID           string `gorm:"column:id;primary_key;autoIncrement:false;comment:帐号ID=渠道_openId" json:"id"`
	UID          int64  `gorm:"column:uid;comment:玩家ID" json:"uid"`
	OtherOpenId  string `gorm:"column:other_open_id;comment:第三方平台open_id" json:"other_open_id"`
	Platform     int32  `gorm:"column:platform;comment:平台ID" json:"platform"`
	Channel      int32  `gorm:"column:channel;comment:渠道ID" json:"channel"`
	RegisterTime int64  `gorm:"column:register_time;comment:注册时间" json:"register_time"`
	ServerId     int32  `gorm:"column:server_id;comment:游戏服ID" json:"server_id"`
}

func (*AccountTable) TableName() string {
	return "account"
}

func CreateAccount(id string, channel int32, openId string, platform int32, serverId int32) *AccountTable {
	account := AccountRepository.Get(id)
	if account != nil {
		return account
	}
	//新注册
	account = &AccountTable{
		ID:           id,
		OtherOpenId:  openId,
		Platform:     platform,
		Channel:      channel,
		ServerId:     serverId,
		RegisterTime: time.Now().Unix(),
		UID:          cherrySnowflake.NextId(),
	}
	return AccountRepository.Add(account)
}
