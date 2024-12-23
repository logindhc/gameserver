package db

import (
	"database/sql/driver"
	"errors"
	clog "gameserver/cherry/logger"
	"gameserver/internal/persistence"
	"gameserver/internal/persistence/repository"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
)

var ShopRepository repository.IRepository[int64, ShopTable]

// ShopTable 角色道具表
type (
	ShopTable struct {
		ID              int64   `gorm:"primaryKey;autoIncrement:false;column:id;comment:id" json:"玩家id"`
		BoxLevel        int32   `gorm:"column:box_level;comment:宝箱等级" json:"box_level"`
		BoxExp          int32   `gorm:"column:box_exp;comment:宝箱经验" json:"box_exp"`
		Shops           ShopMap `gorm:"type:longtext;column:shops;comment:商品集合" json:"shops"`
		BuyCount        ShopMap `gorm:"type:longtext;column:bug_count;comment:购买次数" json:"bug_count"`
		RefreshCount    int32   `gorm:"column:refresh_count;comment:当天已刷新次数" json:"refresh_count"`
		LastRefreshTime int64   `gorm:"column:last_refresh_time;comment:上次刷新日期，第二天登录或者在线零点重置" json:"last_refresh_time"`
	}
	ShopMap map[int]int
)

func (*ShopTable) TableName() string {
	return "shop"
}

func (i *ShopTable) InitRepository() {
	ShopRepository = repository.NewDefaultRepository[int64, ShopTable](database.GetGameDB(), i.TableName())
	persistence.RegisterRepository(ShopRepository)
}
func (i *ShopTable) BeforeCreate(_ *gorm.DB) (err error) {
	if i.BoxLevel == 0 {
		i.BoxLevel = 1
	}
	clog.Debugf("%s# before create", i.TableName())
	return
}

func (i *ShopTable) GetShops() ShopMap {
	if i.Shops == nil {
		i.Shops = ShopMap{}
	}
	return i.Shops
}

// 这个不能带*ShopMap 只能是ShopMap
func (im ShopMap) Value() (driver.Value, error) {
	if im == nil {
		return nil, nil
	}
	data, err := jsoniter.Marshal(im)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func (im *ShopMap) Scan(value interface{}) error {
	if value == nil {
		*im = ShopMap{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	var m ShopMap
	if err := jsoniter.Unmarshal(b, &m); err != nil {
		return err
	}
	*im = m
	return nil
}
