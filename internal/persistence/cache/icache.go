package cache

type ICache[K string | int64, V any] interface {
	Get(id K) *V
	Put(id K, v *V) *V
	Remove(id K)
	Clear()
}
