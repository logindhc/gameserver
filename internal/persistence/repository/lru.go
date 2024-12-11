package repository

import (
	cherryredis "gameserver/cherry/components/redis"
	clog "gameserver/cherry/logger"
	"gameserver/internal/persistence/buffer"
	"gameserver/internal/persistence/cache"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type DefaultRepository[K string | int64, T any] struct {
	cache  cache.ICache[K, T]
	db     *gorm.DB
	buffer buffer.IBuffer[K, T]
	prefix string
}

// NewRedisRepository 使用redis作为缓存，但批量update要从redis获取数据会耗时
func NewRedisRepository[K string | int64, T any](db *gorm.DB, prefix string, params ...time.Duration) IRepository[K, T] {
	// lru 默认过期设置为2个小时
	expiration := 2 * time.Hour
	if len(params) > 0 {
		expiration = params[0]
	}
	redisCache := cache.NewRedisCache[K, T](cherryredis.GetRds(), prefix, expiration)
	r := &DefaultRepository[K, T]{
		db:     db,
		cache:  redisCache,
		buffer: buffer.NewDelayedBuffer[K, T](db, redisCache, prefix),
		prefix: prefix,
	}
	return r
}
func NewDefaultRepository[K string | int64, T any](db *gorm.DB, prefix string, params ...time.Duration) IRepository[K, T] {
	// lru 默认过期设置为2个小时
	expiration := 2 * time.Hour
	if len(params) > 0 {
		expiration = params[0]
	}
	lruCache := cache.NewLRUCache[K, T](expiration)
	r := &DefaultRepository[K, T]{
		db:     db,
		cache:  lruCache,
		buffer: buffer.NewDelayedBuffer[K, T](db, lruCache, prefix),
		prefix: prefix,
	}
	return r
}

func (r *DefaultRepository[K, T]) Get(id K) *T {
	// 从缓存中获取
	entity := r.cache.Get(id)
	if entity != nil {
		return entity
	}
	// 如果缓存中没有，则从数据库中获取
	tx := r.db.Where("id = ?", id).Find(&entity)
	if tx.RowsAffected == 0 {
		return nil
	}
	if entity != nil {
		e := r.cache.Put(id, entity)
		if e != nil {
			return e
		}
	}
	return entity
}

func (r *DefaultRepository[K, T]) GetAll() []*T {
	var entities []*T
	tx := r.db.Find(&entities)
	if tx.Error != nil {
		clog.Errorf("%s#all查询失败", r.prefix)
		return nil
	}
	return entities
}

func (r *DefaultRepository[K, T]) GetOrCreate(id K) *T {
	entity := r.Get(id)
	if entity == nil {
		entity = r.cache.Get(id)
		if entity == nil {
			entity = new(T)
			r.setId(entity, id)
			entity = r.Add(entity)
		}
	}
	return entity
}

func (r *DefaultRepository[K, T]) Add(entity *T) *T {
	if entity == nil {
		return nil
	}
	id := r.getId(entity)
	prev := r.cache.Get(id)
	if prev != nil {
		return prev
	}
	//先缓存再入库
	r.cache.Put(id, entity)
	return r.buffer.Add(entity)
}

func (r *DefaultRepository[K, T]) Remove(id K) {
	r.cache.Remove(id)
	r.buffer.Remove(id)
}

func (r *DefaultRepository[K, T]) Update(entity *T, immediately ...bool) {
	id := r.getId(entity)
	r.cache.Put(id, entity)
	if len(immediately) > 0 && immediately[0] {
		r.db.Save(entity)
		return
	}
	r.buffer.Update(entity)
}

func (r *DefaultRepository[K, T]) Flush() {
	r.buffer.Flush()
}

func (r *DefaultRepository[K, T]) Where(query interface{}, args ...interface{}) (tx *gorm.DB) {
	return r.db.Where(query, args...)
}

func (r *DefaultRepository[K, T]) setId(entity *T, id K) {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	idField := val.FieldByName("ID")
	if !idField.IsValid() {
		clog.Error("ID field not found")
	}
	idField.Set(reflect.ValueOf(id))
}

func (r *DefaultRepository[K, T]) getId(entity *T) K {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	idField := val.FieldByName("ID")
	if !idField.IsValid() {
		clog.Error("ID field not found")
	}
	id, ok := idField.Interface().(K)
	if !ok {
		clog.Error("ID Interface not found")
	}
	return id
}
