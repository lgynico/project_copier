package cherryCluster

import (
	cfacade "github.com/lgynico/project-copier/cherry/facade"
	cnatsCluster "github.com/lgynico/project-copier/cherry/net/cluster/nats"
)

type Component struct {
	cfacade.Component
	cfacade.ICluster
}

func New() *Component {
	return &Component{}
}

func (p *Component) Name() string {
	return "cluster_component"
}

func (p *Component) Init() {
	p.ICluster = p.loadCluster()
	p.ICluster.Init()
}

func (p *Component) OnStop() {
	p.ICluster.Stop()
}

func (p *Component) loadCluster() cfacade.ICluster {
	// FIXME 这里依赖了具体实现
	return cnatsCluster.New(p.App())
}
