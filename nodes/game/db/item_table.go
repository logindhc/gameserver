package db

import (
	"database/sql/driver"
	"errors"
	"gameserver/internal/persistence"
	"gameserver/internal/persistence/repository"
	jsoniter "github.com/json-iterator/go"
)

var ItemRepository repository.IRepository[int64, ItemTable]

// ItemTable 角色道具表
type (
	ItemTable struct {
		ID    int64   `gorm:"primaryKey;autoIncrement:false;column:id;comment:id" json:"玩家id"`
		Items ItemMap `gorm:"type:longtext;column:items;comment:道具集合" json:"items"`
	}
	ItemMap map[int]int
)

func (*ItemTable) TableName() string {
	return "item"
}

func (i *ItemTable) InitRepository() {
	ItemRepository = repository.NewDefaultRepository[int64, ItemTable](database.GetGameDB(), i.TableName())
	persistence.RegisterRepository(ItemRepository)
}

func (i *ItemTable) GetItems() ItemMap {
	if i.Items == nil {
		i.Items = ItemMap{}
	}
	return i.Items
}

// 这个不能带*ItemMap 只能是ItemMap
func (im ItemMap) Value() (driver.Value, error) {
	if im == nil {
		return nil, nil
	}
	data, err := jsoniter.Marshal(im)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func (im *ItemMap) Scan(value interface{}) error {
	if value == nil {
		*im = ItemMap{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	var m ItemMap
	if err := jsoniter.Unmarshal(b, &m); err != nil {
		return err
	}
	*im = m
	return nil
}
