package buffer

import (
	"fmt"
	cstring "gameserver/cherry/extend/string"
	clog "gameserver/cherry/logger"
	"gameserver/internal/utils"
	queue "github.com/duke-git/lancet/v2/datastructure/queue"
	"gorm.io/gorm"
	"math/rand"
	"reflect"
	"strings"
	"sync"
	"time"
)

type LoggerBuffer[K string | int64, T any] struct {
	queue         *queue.ArrayQueue[T]
	db            *gorm.DB
	prefix        string
	monthSharding bool
	locker        sync.Locker
	batchSize     int
}

func NewLoggerBuffer[K string | int64, T any](db *gorm.DB, prefix string, monthSharding bool) IBuffer[K, T] {
	batchSize := 200
	var buffer = &LoggerBuffer[K, T]{
		queue:         queue.NewArrayQueue[T](batchSize),
		db:            db,
		prefix:        prefix,
		monthSharding: monthSharding,
		locker:        &sync.Mutex{},
		batchSize:     batchSize,
	}
	go buffer.flushLoop() // 启动后台任务处理更新与删除
	return buffer
}

// flushLoop 是一个后台循环，用于定期批量入库
func (d *LoggerBuffer[K, T]) flushLoop() {
	interval := time.Duration(flushIntervals+rand.Intn(flushIntervals)) * time.Minute
	clog.Debugf("%s# start flushLoop task interval %d", d.prefix, interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			d.flush(false)
		}

	}
}

// Add 方法实现
func (d *LoggerBuffer[K, T]) Add(entity *T) *T {
	d.locker.Lock()
	d.queue.Enqueue(*entity)
	d.locker.Unlock()
	if d.queue.Size() >= d.batchSize {
		d.flush(false)
	}
	//go d.add(entity)
	return entity
}

func (d *LoggerBuffer[K, T]) Update(entity *T) {
}
func (d *LoggerBuffer[K, T]) Remove(id K) {
}
func (d *LoggerBuffer[K, T]) RemoveAll() {
}

const (
	DefaultOnupdate = "default" // 第一次插入才赋值，后面都不会赋值
	Repeat          = "repeat"
	Total           = "total"
)

type SqlFieldStruct struct {
	fieldName     string
	sqlName       string
	fieldVal      interface{}
	isNull        bool
	isPrimary     bool
	onupdate      string
	typeIsStr     bool
	isMonthShared bool
}

// Flush 方法实现
func (d *LoggerBuffer[K, T]) Flush() {
	d.flush(true)
}
func (d *LoggerBuffer[K, T]) flush(immediately bool) {
	if d.queue.IsEmpty() {
		return
	}
	size := d.queue.Size()
	clog.Debugf("%s# batch add num %d \n", d.prefix, size)
	d.locker.Lock()
	var entities = make([]*T, 0, size)
	for i := 0; i < size; i++ {
		entity, ok := d.queue.Dequeue()
		if !ok {
			break
		}
		entities = append(entities, &entity)
	}
	d.locker.Unlock()
	var sqlBuilder strings.Builder
	for _, entity := range entities {
		sql := d.generateSql(entity)
		sqlBuilder.WriteString(sql)
	}
	if immediately {
		tx := d.db.Exec(sqlBuilder.String())
		if tx.Error != nil {
			clog.Errorf("%s# sync batch add error %s", d.prefix, tx.Error.Error())
			return
		}
		clog.Infof("%s# batch sync add num %d success", d.prefix, size)
	} else {
		go func() {
			tx := d.db.Exec(sqlBuilder.String())
			if tx.Error != nil {
				clog.Errorf("%s# async batch add error %s", d.prefix, tx.Error.Error())
				return
			}
			clog.Infof("%s# async batch add num %d success", d.prefix, size)
		}()
	}
}
func (d *LoggerBuffer[K, T]) add(entity *T) {
	sql := d.generateSql(entity)
	tx := d.db.Exec(sql)
	if tx.Error != nil {
		clog.Errorf("%s# add error %s", d.prefix, tx.Error.Error())
		return
	}
	clog.Infof("%s# add success", d.prefix)
}

func (d *LoggerBuffer[K, T]) generateSql(entity *T) string {
	var sqlBuilder strings.Builder
	//先反射获取对应标记生成的sql
	entityType := reflect.TypeOf(*entity)
	entityValue := reflect.ValueOf(*entity)
	var temp []*SqlFieldStruct
	for i := 0; i < entityType.NumField(); i++ {
		temp = append(temp, d.processField(entityType, entityValue, i))
	}
	var left, values, updates strings.Builder
	monthShardingVal := 0
	for i := 0; i < len(temp); i++ {
		field := temp[i]
		if field.isNull {
			continue
		}
		if field.isMonthShared {
			monthShardingVal = field.fieldVal.(int)
		}
		left.WriteString(fmt.Sprintf("`%s`,", field.sqlName))
		fv := fmt.Sprintf("%v", field.fieldVal)
		if field.typeIsStr {
			fv = fmt.Sprintf("`%v`", field.fieldVal)
		}
		values.WriteString(fmt.Sprintf("%s,", fv))
		if field.onupdate == Repeat {
			updates.WriteString(fmt.Sprintf("`%s`=%s,", field.sqlName, fv))
		} else if field.onupdate == Total {
			updates.WriteString(fmt.Sprintf("`%s`=`%s`+%s,", field.sqlName, field.sqlName, fv))
		} else {
			updates.WriteString(fmt.Sprintf("`%v`=`%v`,", field.sqlName, field.sqlName))
		}
	}
	tbName := d.prefix
	if d.monthSharding {
		tbName = utils.GetMonthTbName(tbName, monthShardingVal)
	}
	leftStr := strings.TrimRight(left.String(), ",")
	valuesStr := strings.TrimRight(values.String(), ",")
	updateStr := strings.TrimRight(updates.String(), ",")
	sqlBuilder.WriteString(fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s;", tbName, leftStr, valuesStr, updateStr))
	sql := sqlBuilder.String()
	return sql
}

func (d *LoggerBuffer[K, T]) processField(entityType reflect.Type, entityValue reflect.Value, i int) *SqlFieldStruct {
	field := entityType.Field(i)
	fieldValue := entityValue.Field(i)
	isPrimary := false
	gormTag := field.Tag.Get("gorm")
	if gormTag != "" {
		gormTag = strings.TrimSpace(gormTag)
		if strings.Contains(gormTag, "primaryKey") {
			isPrimary = true
		}
	}
	onupdate := field.Tag.Get("onupdate")
	if onupdate != Repeat && onupdate != Total {
		onupdate = DefaultOnupdate
	}
	isMonthShared := false
	if field.Tag.Get("monthSharding") == "true" {
		isMonthShared = true
	}
	isNull := false
	value := fieldValue.Interface()
	if fieldValue.Kind() == reflect.Ptr {
		// 如果字段是指针，获取指针指向的值
		if fieldValue.IsNil() {
			isNull = true
		} else {
			value = fieldValue.Elem().Interface()
		}
	}
	typeIsStr := false
	if field.Type.Kind() == reflect.String {
		typeIsStr = true
	}
	sqlField := &SqlFieldStruct{
		fieldName:     field.Name,
		sqlName:       cstring.CamelToSnake(field.Name),
		typeIsStr:     typeIsStr,
		isPrimary:     isPrimary,
		fieldVal:      value,
		isNull:        isNull,
		onupdate:      onupdate,
		isMonthShared: isMonthShared,
	}
	return sqlField
}
