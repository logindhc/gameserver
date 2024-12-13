package buffer

import (
	clog "gameserver/cherry/logger"
	"gameserver/internal/concurrent"
	"gameserver/internal/persistence/cache"
	"gorm.io/gorm"
	"math/rand"
	"sync/atomic"
	"time"
)

const batchSize = 1000

type DelayedBuffer[K string | int64, T any] struct {
	cache      cache.ICache[K, T]
	db         *gorm.DB
	prefix     string
	updates    *atomic.Value
	deletes    *concurrent.ConcurrentSet[K]
	bufferSize int
}

func NewDelayedBuffer[K string | int64, T any](db *gorm.DB, cache cache.ICache[K, T], prefix string) *DelayedBuffer[K, T] {
	bufferSize := 5000
	buffer := &DelayedBuffer[K, T]{
		cache:      cache,
		db:         db,
		prefix:     prefix,
		updates:    &atomic.Value{},
		deletes:    concurrent.NewConcurrentSet[K](),
		bufferSize: bufferSize,
	}
	buffer.updates.Store(concurrent.NewConcurrentSet[K]())
	go buffer.flushLoop() // 启动后台任务处理更新与删除
	return buffer
}

// flushLoop 是一个后台循环，用于定期将缓存中的更改同步到数据库
func (d *DelayedBuffer[K, T]) flushLoop() {
	interval := time.Duration(flushIntervals+rand.Intn(flushIntervals)) * time.Minute
	clog.Debugf("%s# start flushLoop task interval %d", d.prefix, interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			d.flush()
		}
	}
}

// Add 方法实现
func (d *DelayedBuffer[K, T]) Add(entity *T) *T {
	k := getKey(entity)
	d.deletes.Remove(k.(K))
	//go func() {
	tx := d.db.Create(entity)
	if tx.Error == nil {
		clog.Debugf("%s#id:%v 添加成功", d.prefix, k)
	} else {
		clog.Errorf("%s#id:%v 添加失败 %v", d.prefix, k, entity)
	}
	//}()
	return entity
}

// Update 方法实现
func (d *DelayedBuffer[K, T]) Update(entity *T) {
	id := getKey(entity)
	if d.deletes.Has(id.(K)) {
		clog.Errorf("%s#id:%v 更新时已经被删除", d.prefix, id)
		return
	}
	pending := d.updates.Load().(*concurrent.ConcurrentSet[K])
	pending.Add(id.(K))
	size := pending.Size()
	//clog.Debugf("updates %p %v \n", pending, size)
	if size >= d.bufferSize {
		d.flush()
	}
}

// Remove 方法实现
func (d *DelayedBuffer[K, T]) Remove(id K) {
	d.deletes.Add(id)
	d.updates.Load().(*concurrent.ConcurrentSet[K]).Remove(id)
	var entity T
	d.db.Model(entity).Where("id = ?", id).Delete(nil)
	d.deletes.Remove(id)
	clog.Debugf("%s#id:%v 删除成功", d.prefix, id)
}

// RemoveAll 方法实现
func (d *DelayedBuffer[K, T]) RemoveAll() {
	// 清空缓存并触发刷新
	d.cache.Clear()
	d.deletes.Clear()
	d.updates.Load().(*concurrent.ConcurrentSet[K]).Clear()
	var entity = new(*T)
	d.db.Delete(entity)
}

// Flush 方法实现
func (d *DelayedBuffer[K, T]) Flush() {
	d.flush()
}

func (d *DelayedBuffer[K, T]) flush() {
	// 处理更新
	flushes := d.updates.Swap(concurrent.NewConcurrentSet[K]())
	f := flushes.(*concurrent.ConcurrentSet[K])
	size := f.Size()
	if size <= 0 {
		return
	}
	clog.Debugf("%s# [%p] delayerd flush num %d", d.prefix, flushes, size)
	all := f.All()
	count := 0
	start := time.Now()
	var entities = make([]*T, 0, size)
	for _, id := range all {
		entity := d.cache.Get(id)
		if entity == nil {
			clog.Errorf("%s#id:%v 更新失败，缓存中不存在", d.prefix, id)
			continue
		}
		count++
		entities = append(entities, entity)
		f.Remove(id)
	}

	for i := 0; i < len(entities); i += batchSize {
		end := i + batchSize
		if end > len(entities) {
			end = len(entities)
		}
		batch := entities[i:end]
		if err := d.db.Save(&batch).Error; err != nil {
			clog.Errorf("%s# Batch save failed: %v", d.prefix, err)
		}
	}
	//
	//tx := d.db.Save(entities)
	//if tx.Error != nil {
	//	clog.Errorf("%s# 批量更新失败 %v", d.prefix, tx.Error)
	//}
	clog.Debugf("%s# [%p] delayerd sync flush num %d success %d, cos %v", d.prefix, flushes, size, count, time.Since(start))
}
