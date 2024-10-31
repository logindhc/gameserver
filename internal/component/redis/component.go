package redis

import (
	"context"
	cherryFacade "gameserver/cherry/facade"
	clog "gameserver/cherry/logger"
	cprofile "gameserver/cherry/profile"
	"github.com/go-redis/redis/v8"
)

var (
	rdb *Component
)

type (
	Component struct {
		redisConfig
		cherryFacade.Component
		rdb *redis.Client
	}

	redisConfig struct {
		Address  string `json:"address"`  // redis地址
		Password string `json:"password"` // 密码
		DB       int    `json:"db"`       // db index
	}
)

func (c *Component) Init() {
	dataConfig := cprofile.GetConfig("data_config").GetConfig(c.Name())
	if dataConfig.Unmarshal(&c.redisConfig) != nil {
		clog.Warnf("[data_config]->[%s] node in `%s` file not found.", c.Name(), cprofile.Name())
		return
	}
	c.newRedis()
}
func New() *Component {
	c := &Component{}
	rdb = c
	return c
}

func (c *Component) Name() string {
	return "redis"
}

func (r *Component) newRedis() {
	r.rdb = redis.NewClient(&redis.Options{
		Addr:     r.Address,
		Password: r.Password,
		DB:       r.DB,
	})
	_, err := r.rdb.Ping(context.Background()).Result()
	if err != nil {
		clog.Panic(err)
		return
	}
	rdb = r
	clog.Infof("%s connected.", r.rdb.String())
}

func (c *Component) OnAfterInit() {
}

func GetRds() *redis.Client {
	return rdb.rdb
}
