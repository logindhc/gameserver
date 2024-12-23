package db

import (
	"database/sql/driver"
	"errors"
	"gameserver/internal/persistence"
	"gameserver/internal/persistence/repository"
	jsoniter "github.com/json-iterator/go"
)

var HeroRepository repository.IRepository[int64, HeroTable]

// HeroTable 角色英雄表
type (
	HeroTable struct {
		ID    int64   `gorm:"primaryKey;autoIncrement:false;column:id;comment:id" json:"玩家id"`
		Heros HeroMap `gorm:"type:longtext;column:heros;comment:英雄集合" json:"heros"`
	}
	HeroMap map[int]int
)

func (*HeroTable) TableName() string {
	return "hero"
}

func (i *HeroTable) InitRepository() {
	HeroRepository = repository.NewDefaultRepository[int64, HeroTable](database.GetGameDB(), i.TableName())
	persistence.RegisterRepository(HeroRepository)
}

func (i *HeroTable) GetHeros() HeroMap {
	if i.Heros == nil {
		i.Heros = HeroMap{}
	}
	return i.Heros
}

// 这个不能带*ItemMap 只能是ItemMap
func (im HeroMap) Value() (driver.Value, error) {
	if im == nil {
		return nil, nil
	}
	data, err := jsoniter.Marshal(im)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func (im *HeroMap) Scan(value interface{}) error {
	if value == nil {
		*im = HeroMap{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	var m HeroMap
	if err := jsoniter.Unmarshal(b, &m); err != nil {
		return err
	}
	*im = m
	return nil
}
