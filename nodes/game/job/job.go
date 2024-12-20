package job

import (
	clog "gameserver/cherry/logger"
	cactor "gameserver/cherry/net/actor"
	"gameserver/internal/cache"
	"gameserver/internal/job"
	"gameserver/nodes/game/module/online"
	"time"
)

type (
	ActorJob struct {
		cactor.Base
	}
)

func (p *ActorJob) AliasID() string {
	return "job"
}

// OnInit 注册函数
func (p *ActorJob) OnInit() {
	p.job()
}
func (p *ActorJob) OnStop() {
	cache.DelGameNodeOnline(p.App().NodeId())
}

func (p *ActorJob) job() {
	//启动就上报一次
	cache.UpdateGameNodeOnline(p.App().NodeId(), float64(online.Count()))
	job.GlobalTimer.BuildEveryFunc(time.Minute, func() {
		nodeId := p.App().NodeId()
		count := online.Count()
		cache.UpdateGameNodeOnline(nodeId, float64(count))
		clog.Infof("[job] [nodeId = %s, onlineCount = %d]", nodeId, count)
	})
}
