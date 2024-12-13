package cache

import (
	clog "gameserver/cherry/logger"
	"sync"
	"sync/atomic"
)

// 配置高位和低位的占用位数
const (
	ServerIDBits  = 10
	IncrementBits = 27
	MaxServerID   = (1 << ServerIDBits) - 1
	MaxIncrement  = (1 << IncrementBits) - 1
)

type PlayerIDGenerator interface {
	InitializeIncrement(serverID int32, initialIncrement int64)
	NextID(serverID int32) int64
}

// DefIDGenerator 支持多个区服的 ID 生成器
type DefIDGenerator struct {
	serverIncrements sync.Map // 使用 sync.Map 保存每个区服的自增值
}

// NewIDGenerator 创建一个新的多区服 玩家ID 生成器
func NewIDGenerator() PlayerIDGenerator {
	return &DefIDGenerator{}
}

// InitializeIncrement 初始化某个区服的自增值（线程安全）
func (gen *DefIDGenerator) InitializeIncrement(serverID int32, initialIncrement int64) {
	if serverID > MaxServerID {
		clog.Panicf("ServerID exceeds the maximum value %d", MaxServerID)
	}
	if initialIncrement > MaxIncrement {
		clog.Panicf("Initial increment exceeds the maximum value %d", MaxIncrement)
	}
	// 使用 sync.Map 存储或更新自增值
	gen.serverIncrements.Store(serverID, &initialIncrement)
}

// NextID 生成某个区服的下一个 ID（线程安全）
func (gen *DefIDGenerator) NextID(serverID int32) int64 {
	if serverID > MaxServerID {
		clog.Panicf("ServerID exceeds the maximum value %d", MaxServerID)
	}
	// 获取或初始化区服的 increment
	incrementPtr, _ := gen.serverIncrements.LoadOrStore(serverID, new(int64))
	// 原子递增并返回
	inc := atomic.AddInt64(incrementPtr.(*int64), 1)
	if inc > MaxIncrement {
		clog.Panicf("Increment value for server %d exceeds the maximum limit", serverID)
	}
	// 生成 ID
	return (inc << ServerIDBits) | int64(serverID)
}
