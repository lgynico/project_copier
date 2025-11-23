package cherryDiscovery

import (
	cfacade "github.com/lgynico/project-copier/cherry/facade"
	clog "github.com/lgynico/project-copier/cherry/logger"
	cherryProfile "github.com/lgynico/project-copier/cherry/profile"
)

const Name = "discovery_component"

type Component struct {
	cfacade.Component
	cfacade.IDiscovery
}

func New() *Component {
	return &Component{}
}

func (*Component) Name() string {
	return Name
}

func (p *Component) Init() {
	config := cherryProfile.GetConfig("cluster").GetConfig("discovery")
	if config.LastError() != nil {
		clog.Error("`cluster` property not found in profile file.")
		return
	}

	mode := config.GetString("mode")
	if mode == "" {
		clog.Error("`discovery.mode` property not found in profile file.s")
		return
	}

	discovery, ok := discoveryMap[mode]
	if discovery == nil || !ok {
		clog.Errorf("mode = %s not found in discovery config.", mode)
		return
	}

	clog.Infof("Select discovery [mode = %s]", mode)
	p.IDiscovery = discovery
	p.IDiscovery.Load(p.App())
}

func (p *Component) Stop() {
	p.IDiscovery.Stop()
}
