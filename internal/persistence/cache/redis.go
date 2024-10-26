package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	clog "gameserver/cherry/logger"
	"github.com/go-redis/redis/v8"
	"time"
)

// RedisCache 是基于 Redis 实现的缓存
type RedisCache[K string | int64, T any] struct {
	client     *redis.Client
	prefix     string
	expiration time.Duration
}

// NewRedisCache 创建一个新的 RedisCache 实例
func NewRedisCache[K string | int64, T any](client *redis.Client, prefix string, expiration time.Duration) *RedisCache[K, T] {
	return &RedisCache[K, T]{client: client, prefix: prefix, expiration: expiration}
}

// Get 从缓存中获取指定 ID 的值
func (r *RedisCache[K, T]) Get(id K) *T {
	key := fmt.Sprintf("%s:%v", r.prefix, id)
	data, err := r.client.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		return nil // 没有找到对应的值
	} else if err != nil {
		clog.Errorf("%s#id:%v Get失败", r.prefix, id)
		return nil
	}

	var entity T
	err = json.Unmarshal([]byte(data), &entity)
	if err != nil {
		clog.Errorf("%s#id:%v Get反序列化失败", r.prefix, id)
		return nil
	}
	return &entity
}

// Put 将值放入缓存中
func (r *RedisCache[K, T]) Put(id K, entity *T) *T {
	key := fmt.Sprintf("%s:%v", r.prefix, id)
	data, err := json.Marshal(entity)
	if err != nil {
		clog.Errorf("%s#id:%v Put序列化失败", r.prefix, id)
		return nil
	}
	tx := r.client.Set(context.Background(), key, data, r.expiration)
	if tx.Err() != nil {
		clog.Errorf("%s#id:%v Put失败", r.prefix, id)
		return nil
	}
	return entity
}

// Remove 从缓存中移除指定 ID 的值
func (r *RedisCache[K, T]) Remove(id K) {
	key := fmt.Sprintf("%s:%v", r.prefix, id)
	err := r.client.Del(context.Background(), key).Err()
	if err != nil {
		clog.Errorf("%s#id:%v Remove失败", r.prefix, id)
		return
	}
}

// Clear 清空缓存，一般不建议使用
func (r *RedisCache[K, T]) Clear() {
	keys, err := r.client.Keys(context.Background(), fmt.Sprintf("%s:*", r.prefix)).Result()
	if err != nil {
		clog.Errorf("%s#Clear失败", r.prefix)
		return
	}
	if len(keys) <= 0 {
		return
	}
	err = r.client.Del(context.Background(), keys...).Err()
	if err != nil {
		clog.Errorf("%s#Clear失败", r.prefix)
	}
}
