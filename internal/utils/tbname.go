package utils

import (
	"fmt"
	ctime "gameserver/cherry/extend/time"
)

func GetMonthTbName(tableNamePrefix string, fieldVal int32) string {
	if fieldVal == 0 {
		return fmt.Sprintf("%s_%s", tableNamePrefix, ctime.Now().ToShortMonthFormat())
	}
	yyyy := fieldVal / 10000
	mm := fieldVal % 10000 / 100
	return fmt.Sprintf("%s_%04d%02d", tableNamePrefix, yyyy, mm)
}
