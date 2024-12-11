package db

import (
	"gameserver/internal/persistence"
	"gameserver/internal/persistence/repository"
	"gameserver/internal/utils"
)

var logLoginPrefix = "log_login"

var LogLoginRepository repository.IRepository[int64, LogLogin]

func (d *LogLogin) InitRepository() {
	LogLoginRepository = repository.NewLoggerRepository[int64, LogLogin](database.GetLogDB(), logLoginPrefix, true)
	persistence.RegisterRepository(LogLoginRepository)
}

type LogLogin struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement:false;comment:玩家ID"`
	DayIndex   int32  `gorm:"column:day_index;primaryKey;autoIncrement:false;comment:登录日期" monthSharding:"true" partition:"day_index"`
	FirstTime  *int64 `gorm:"column:first_time;comment:首次登录时间"`
	LastTime   *int64 `gorm:"column:last_time" onupdate:"repeat;comment:最近登录时间"`
	TotalCount int32  `gorm:"column:total_count" onupdate:"total;comment:累计次数"`
}

func (d *LogLogin) TableName() string {
	//DayIndex格式为yyyyMMdd
	return utils.GetMonthTbName(logLoginPrefix, d.DayIndex) //用这个DayIndex 做表名分月
}
