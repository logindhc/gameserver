package db

import (
	cherrySnowflake "gameserver/cherry/extend/snowflake"
	clog "gameserver/cherry/logger"
	idgener "gameserver/internal/component/id"
	"gameserver/internal/persistence"
	"gameserver/internal/persistence/repository"
	"time"
)

var AccountRepository repository.IRepository[int64, AccountTable]

func (p *AccountTable) InitRepository() {
	AccountRepository = repository.NewSynchRepository[int64, AccountTable](database.GetGameDB())
	persistence.RegisterRepository(AccountRepository)
}

// AccountTable 帐号信息表
type AccountTable struct {
	ID           int64  `gorm:"column:id;primary_key;autoIncrement:false;comment:帐号ID" json:"id"`
	OpenId       string `gorm:"column:open_id;index;comment:第三方open_id" json:"open_id"`
	Channel      int32  `gorm:"column:channel;comment:渠道ID" json:"channel"`
	Platform     int32  `gorm:"column:platform;comment:平台ID" json:"platform"`
	RegisterTime int64  `gorm:"column:register_time;comment:注册时间" json:"register_time"`
	ServerId     int32  `gorm:"column:server_id;comment:游戏服ID" json:"server_id"`
	UID          int64  `gorm:"column:uid;comment:玩家ID" json:"uid"`
}

func (*AccountTable) TableName() string {
	return "account"
}

func GetAccount(openId string) *AccountTable {
	if account, ok := openId2AccountCache.GetIfPresent(openId); ok {
		return account.(*AccountTable)
	}
	account := &AccountTable{}
	tx := AccountRepository.Where("open_id = ?", openId).Find(account)
	if tx.RowsAffected == 0 {
		//没有账号
		return nil
	}
	openId2AccountCache.Put(openId, account)
	return account
}

func CreateAccount(channel int32, openId string, platform int32, serverId int32) *AccountTable {
	//新注册
	accountId := cherrySnowflake.NextId()
	account := &AccountTable{
		ID:           accountId,
		OpenId:       openId,
		Platform:     platform,
		Channel:      channel,
		ServerId:     serverId,
		RegisterTime: time.Now().Unix(),
		UID:          idgener.PlayerIdGenerator.NextID(serverId),
	}
	//先入库
	AccountRepository.Add(account)
	//再放到缓存中
	openId2AccountCache.Put(openId, account)
	return account
}

func loadDevAccount() {

}

func loadMaxID() {
	var results []struct {
		ServerID int32
		MaxID    int64
	}
	err := database.GetGameDB().Model(&AccountTable{}).Select("server_id,MAX(uid) as max_id").Group("server_id").Scan(&results).Error
	if err != nil {
		clog.Panicf("loadMaxID err: %v", err)
		return
	}
	for _, ret := range results {
		serverId, increment := idgener.ParseID(ret.MaxID)
		clog.Infof("loadMaxID serverID: %d, maxID: %d parse serverId: %d increment: %d", ret.ServerID, ret.MaxID, serverId, increment)
		idgener.PlayerIdGenerator.InitializeIncrement(serverId, increment)
	}
}
