package cherryDiscovery

import (
	cfacade "gameserver/cherry/facade"
	clog "gameserver/cherry/logger"
)

var (
	discoveryMap = make(map[string]cfacade.IDiscovery)
)

func init() {
	Register(&DiscoveryDefault{})
	Register(&DiscoveryNATS{})
	//RegisterDiscovery(&DiscoveryETCD{})
}

func Register(discovery cfacade.IDiscovery) {
	if discovery == nil {
		clog.Fatal("Discovery instance is nil")
		return
	}

	if discovery.Name() == "" {
		clog.Fatalf("Discovery name is empty. %T", discovery)
		return
	}
	discoveryMap[discovery.Name()] = discovery
}
