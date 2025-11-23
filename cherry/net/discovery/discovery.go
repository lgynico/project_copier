package cherryDiscovery

import (
	cfacade "github.com/lgynico/project-copier/cherry/facade"
	clog "github.com/lgynico/project-copier/cherry/logger"
)

var discoveryMap = make(map[string]cfacade.IDiscovery)

func init() {

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
