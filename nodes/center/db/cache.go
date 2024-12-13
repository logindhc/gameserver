package db

import (
	"github.com/goburrow/cache"
	"time"
)

var (

	// accountId缓存 key:openId, value:account
	openId2AccountCache = cache.New(
		cache.WithMaximumSize(65535),
		cache.WithExpireAfterAccess(120*time.Minute),
	)
)
