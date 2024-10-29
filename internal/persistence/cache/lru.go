package cache

import (
	"github.com/goburrow/cache"
	"time"
)

type LRUCache[K string | int64, V any] struct {
	cache.Cache
}

func NewLRUCache[K string | int64, V any](expiration time.Duration) *LRUCache[K, V] {
	c := cache.New(
		cache.WithMaximumSize(-1),
		cache.WithExpireAfterAccess(expiration))
	return &LRUCache[K, V]{c}
}

func (l *LRUCache[K, V]) Get(id K) *V {
	v, ok := l.GetIfPresent(id)
	if !ok {
		return nil
	}
	return v.(*V)
}

func (l *LRUCache[K, V]) Put(id K, value *V) *V {
	l.Cache.Put(id, value)
	return value
}

func (l *LRUCache[K, V]) Remove(key K) {
	l.Invalidate(key)
}

func (l *LRUCache[K, V]) Clear() {
	l.InvalidateAll()
}
