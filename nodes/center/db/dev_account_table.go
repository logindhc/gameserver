package db

import (
	cherryError "gameserver/cherry/error"
	cherrySnowflake "gameserver/cherry/extend/snowflake"
	cherryTime "gameserver/cherry/extend/time"
	cherryLogger "gameserver/cherry/logger"
	"gameserver/internal/code"
	"gameserver/internal/persistence"
	"gameserver/internal/persistence/repository"
)

// DevAccountTable 开发模式的帐号信息表(platform.TypeDevMode)
type DevAccountTable struct {
	ID          int64  `gorm:"column:id;primary_key;autoIncrement:false;comment:帐号id" json:"id"`
	AccountName string `gorm:"column:account_name;comment:帐号名" json:"accountName"`
	Password    string `gorm:"column:password;comment:密码" json:"-"`
	CreateIP    string `gorm:"column:create_ip;comment:创建ip" json:"createIP"`
	CreateTime  int64  `gorm:"column:create_time;comment:创建时间" json:"createTime"`
}

func (*DevAccountTable) TableName() string {
	return "dev_account"
}

var AccountRepository *repository.LruRepository[int64, DevAccountTable]

func (p *DevAccountTable) InitRepository() {
	AccountRepository = repository.NewLruRepository[int64, DevAccountTable](database.GetGameDB(), p.TableName())
	persistence.RegisterRepository(AccountRepository)
}

func DevAccountRegister(accountName, password, ip string) int32 {
	devAccount, _ := DevAccountWithName(accountName)
	if devAccount != nil {
		return code.AccountNameIsExist
	}

	devAccountTable := &DevAccountTable{
		ID:          cherrySnowflake.NextId(),
		AccountName: accountName,
		Password:    password,
		CreateIP:    ip,
		CreateTime:  cherryTime.Now().Unix(),
	}
	add := AccountRepository.Add(devAccountTable)
	if add == nil {
		return code.AccountRegisterError
	}
	devAccountCache.Put(accountName, devAccountTable)
	return code.OK
}

func DevAccountWithName(accountName string) (*DevAccountTable, error) {
	val, found := devAccountCache.GetIfPresent(accountName)
	if found == false {
		val = new(DevAccountTable)
		tx := AccountRepository.Where("account_name", accountName).Find(&val)
		if tx.RowsAffected == 0 {
			return nil, cherryError.Error("account not found")
		}
		devAccountCache.Put(accountName, val)
	}

	return val.(*DevAccountTable), nil
}

// loadDevAccount 节点启动时预加载帐号数据
func loadDevAccount() {
	list := AccountRepository.GetAll()
	for _, account := range list {
		devAccountCache.Put(account.AccountName, account)
	}
	cherryLogger.Info("preload DevAccountTable")
}
