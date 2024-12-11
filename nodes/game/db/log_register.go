package db

import (
	"gameserver/internal/persistence"
	"gameserver/internal/persistence/repository"
)

var logRegisterPrefix = "log_register"

var LogRegisterRepository repository.IRepository[string, LogRegister]

func (d *LogRegister) InitRepository() {
	LogRegisterRepository = repository.NewLoggerRepository[string, LogRegister](database.GetLogDB(), logRegisterPrefix)
	persistence.RegisterRepository(LogRegisterRepository)
}

type LogRegister struct {
	Device   string `gorm:"column:device;primaryKey;autoIncrement:false;comment:设备ID or OpenId"`
	Channel  int32  `gorm:"column:channel;primaryKey;autoIncrement:false;comment:渠道"`
	Platform int32  `gorm:"column:platform;primaryKey;autoIncrement:false;comment:平台"`
	Time     int64  `gorm:"column:time;comment:激活时间戳"`
}

func (d *LogRegister) TableName() string {
	return logRegisterPrefix
}
