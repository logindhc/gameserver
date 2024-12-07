package db

import (
	"gameserver/internal/persistence"
	"gameserver/internal/persistence/repository"
	"gameserver/internal/utils"
)

var DotLoginRepository *repository.LoggerRepository[int64, DotLogin]

func (log *DotLogin) InitRepository() {
	DotLoginRepository = repository.NewLoggerRepository[int64, DotLogin](database.GetLogDB(), "dot_login", true)
	persistence.RegisterRepository(DotLoginRepository)
}

type DotLogin struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement:false" `
	DayIndex   int    `gorm:"column:day_index;primaryKey;autoIncrement:false" monthSharding:"true" partition:"day_index"`
	FirstTime  *int64 `gorm:"column:first_time"`
	LastTime   *int64 `gorm:"column:last_time" onupdate:"repeat"`
	TotalCount *int   `gorm:"column:total_count" onupdate:"total"`
}

func (log *DotLogin) TableName() string {
	//DayIndex格式为yyyyMMdd
	return utils.GetMonthTbName("dot_login", log.DayIndex) //用这个DayIndex 做表名分月
}
