package job

import (
	ctimeWheel "gameserver/cherry/extend/time_wheel"
	"time"
)

var (
	GlobalTimer = ctimeWheel.NewTimeWheel(10*time.Millisecond, 3600)
)

func init() {
	GlobalTimer.Start()
}
