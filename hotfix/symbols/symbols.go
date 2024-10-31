package symbols

import (
	"reflect"

	"github.com/traefik/yaegi/stdlib"
)

var Symbols = map[string]map[string]reflect.Value{}

func init() {
	for k, v := range stdlib.Symbols {
		Symbols[k] = v
	}
}

// 点击生成符号表
//go:generate yaegi extract gameserver/hotfix
//go:generate yaegi extract gameserver/internal/code
//go:generate yaegi extract gameserver/internal/event
//go:generate yaegi extract gameserver/internal/event
//go:generate yaegi extract gameserver/internal/pb
//go:generate yaegi extract gameserver/internal/data
//go:generate yaegi extract gameserver/internal/session_key gameserver/internal/sessionKey
//go:generate yaegi extract gameserver/nodes/game/module/player
//go:generate yaegi extract gameserver/nodes/game/module/online
//go:generate yaegi extract gameserver/nodes/game/db

//go:generate yaegi extract gameserver/hotfix/test/model
//go:generate yaegi extract gameserver/cherry/components/cron gameserver/cherry/components/cherryCron
//go:generate yaegi extract gameserver/cherry/extend/string gameserver/cherry/extend/cherryString
//go:generate yaegi extract gameserver/cherry/net/proto gameserver/cherry/components/cherryProto
