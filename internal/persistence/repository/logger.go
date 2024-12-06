package repository

import (
	"gameserver/internal/persistence/buffer"
	"gorm.io/gorm"
)

type LoggerRepository[K string | int64, T any] struct {
	db            *gorm.DB
	buffer        buffer.IBuffer[K, T]
	prefix        string
	monthSharding bool
}

func NewLoggerRepository[K string | int64, T any](db *gorm.DB, prefix string, monthSharding bool) *LoggerRepository[K, T] {
	r := &LoggerRepository[K, T]{
		db:            db,
		buffer:        buffer.NewLoggerBuffer[K, T](db, prefix, monthSharding),
		prefix:        prefix,
		monthSharding: monthSharding,
	}
	return r
}

func (r *LoggerRepository[K, T]) Add(entity *T) *T {
	if entity == nil {
		return nil
	}
	r.buffer.Add(entity)
	return entity
}
func (r *LoggerRepository[K, T]) Flush() {
	r.buffer.Flush()
}

func (r *LoggerRepository[K, T]) Get(id K) *T {
	return nil
}
func (r *LoggerRepository[K, T]) GetAll() []*T {
	return nil
}
func (r *LoggerRepository[K, T]) GetOrCreate(id K) *T {
	return nil
}
func (r *LoggerRepository[K, T]) Remove(id K) {
}
func (r *LoggerRepository[K, T]) Update(entity *T, immediately ...bool) {
}
func (r *LoggerRepository[K, T]) Where(query interface{}, args ...interface{}) (tx *gorm.DB) {
	return nil
}
