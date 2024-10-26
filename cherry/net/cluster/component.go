package cherryCluster

import (
	cfacade "gameserver/cherry/facade"
	cherryNatsCluster "gameserver/cherry/net/cluster/nats_cluster"
)

const (
	Name = "cluster_component"
)

type Component struct {
	cfacade.Component
	cfacade.ICluster
}

func New() *Component {
	return &Component{}
}

func (c *Component) Name() string {
	return Name
}

func (c *Component) Init() {
	c.ICluster = c.loadCluster()
	c.ICluster.Init()
}

func (c *Component) OnStop() {
	c.ICluster.Stop()
}

func (c *Component) loadCluster() cfacade.ICluster {
	return cherryNatsCluster.New(c.App())
}
