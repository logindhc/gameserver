package cherryActor

import (
	"time"

	ctimeWheel "gameserver/cherry/extend/time_wheel"
)

var (
	globalTimer = ctimeWheel.NewTimeWheel(10*time.Millisecond, 3600)
)

func init() {
	globalTimer.Start()
}
