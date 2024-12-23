package db

import (
	"database/sql/driver"
	"errors"
	"gameserver/internal/persistence"
	"gameserver/internal/persistence/repository"
	jsoniter "github.com/json-iterator/go"
)

var CurrencyRepository repository.IRepository[int64, CurrencyTable]

// CurrencyTable 角色货币表
type (
	CurrencyTable struct {
		ID     int64    `gorm:"primaryKey;autoIncrement:false;column:id;comment:id" json:"玩家id"`
		Moneys MoneyMap `gorm:"type:longtext;column:moneys;comment:货币集合" json:"moneys"`
	}
	MoneyMap map[int]int
)

func (*CurrencyTable) TableName() string {
	return "currency"
}

func (i *CurrencyTable) InitRepository() {
	CurrencyRepository = repository.NewDefaultRepository[int64, CurrencyTable](database.GetGameDB(), i.TableName())
	persistence.RegisterRepository(CurrencyRepository)
}

func (i *CurrencyTable) GetMaps() MoneyMap {
	if i.Moneys == nil {
		i.Moneys = MoneyMap{}
	}
	return i.Moneys
}

// 这个不能带*ItemMap 只能是ItemMap
func (im MoneyMap) Value() (driver.Value, error) {
	if im == nil {
		return nil, nil
	}
	data, err := jsoniter.Marshal(im)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func (im *MoneyMap) Scan(value interface{}) error {
	if value == nil {
		*im = MoneyMap{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	var m MoneyMap
	if err := jsoniter.Unmarshal(b, &m); err != nil {
		return err
	}
	*im = m
	return nil
}
