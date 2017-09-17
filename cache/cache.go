package cache

import (
	"encoding/binary"
	"sync"
	"sync/atomic"
)

//freecache
//会根据cache size 进行预分配
type Cache struct {
	locks     [256]sync.Mutex
	segments  [256]segment
	hitCount  int64
	missCount int64
}

func hashFunc(data []byte) uint64 {
	return Sum64(data)
}

// cache 最小对象数量为512KB
// 假如 cache size 设置相当大的情况下，
// 建议设置`debug.SetGCPercent()` 来限制内存的消耗和GC的停顿时间
func NewCache(size int) (cache *Cache) {
	if size < 512*1024 {
		size = 512 * 1024
	}
	cache = new(Cache)
	for i := 0; i < 256; i++ {
		cache.segments[i] = newSegment(size/256, i)
	}
	return
}

// 如果key的数量大于65535 或者 value的数量大于总cache的1/1024
// 缓存entry不会被写入成功。
// expireSeconds当设置为0时候，表示永不过期， 但当cache满后，旧的entry也会被淘汰
func (cache *Cache) Set(key, value []byte, expireSeconds int) (err error) {
	hashVal := hashFunc(key)
	segId := hashVal & 255
	cache.locks[segId].Lock()
	err = cache.segments[segId].set(key, value, hashVal, expireSeconds)
	cache.locks[segId].Unlock()
	return
}

func (cache *Cache) Get(key []byte) (value []byte, err error) {
	hashVal := hashFunc(key)
	segId := hashVal & 255
	cache.locks[segId].Lock()
	value, err = cache.segments[segId].get(key, hashVal)
	cache.locks[segId].Unlock()
	if err == nil {
		atomic.AddInt64(&cache.hitCount, 1)
	} else {
		atomic.AddInt64(&cache.missCount, 1)
	}
	return
}

func (cache *Cache) TTL(key []byte) (timeLeft uint32, err error) {
	hashVal := hashFunc(key)
	segId := hashVal & 255
	timeLeft, err = cache.segments[segId].ttl(key, hashVal)
	return
}

func (cache *Cache) Del(key []byte) (affected bool) {
	hashVal := hashFunc(key)
	segId := hashVal & 255
	cache.locks[segId].Lock()
	affected = cache.segments[segId].del(key, hashVal)
	cache.locks[segId].Unlock()
	return
}

func (cache *Cache) SetInt(key int64, value []byte, expireSeconds int) (err error) {
	var bKey [8]byte
	binary.LittleEndian.PutUint64(bKey[:], uint64(key))
	return cache.Set(bKey[:], value, expireSeconds)
}

func (cache *Cache) GetInt(key int64) (value []byte, err error) {
	var bKey [8]byte
	binary.LittleEndian.PutUint64(bKey[:], uint64(key))
	return cache.Get(bKey[:])
}

func (cache *Cache) DelInt(key int64) (affected bool) {
	var bKey [8]byte
	binary.LittleEndian.PutUint64(bKey[:], uint64(key))
	return cache.Del(bKey[:])
}

func (cache *Cache) EvacuateCount() (count int64) {
	for i := 0; i < 256; i++ {
		count += atomic.LoadInt64(&cache.segments[i].totalEvacuate)
	}
	return
}

func (cache *Cache) ExpiredCount() (count int64) {
	for i := 0; i < 256; i++ {
		count += atomic.LoadInt64(&cache.segments[i].totalExpired)
	}
	return
}

func (cache *Cache) EntryCount() (entryCount int64) {
	for i := 0; i < 256; i++ {
		entryCount += atomic.LoadInt64(&cache.segments[i].entryCount)
	}
	return
}

func (cache *Cache) AverageAccessTime() int64 {
	var entryCount, totalTime int64
	for i := 0; i < 256; i++ {
		totalTime += atomic.LoadInt64(&cache.segments[i].totalTime)
		entryCount += atomic.LoadInt64(&cache.segments[i].totalCount)
	}
	if entryCount == 0 {
		return 0
	} else {
		return totalTime / entryCount
	}
}

func (cache *Cache) HitCount() int64 {
	return atomic.LoadInt64(&cache.hitCount)
}

func (cache *Cache) LookupCount() int64 {
	return atomic.LoadInt64(&cache.hitCount) + atomic.LoadInt64(&cache.missCount)
}

func (cache *Cache) HitRate() float64 {
	lookupCount := cache.LookupCount()
	if lookupCount == 0 {
		return 0
	} else {
		return float64(cache.HitCount()) / float64(lookupCount)
	}
}

func (cache *Cache) OverwriteCount() (overwriteCount int64) {
	for i := 0; i < 256; i++ {
		overwriteCount += atomic.LoadInt64(&cache.segments[i].overwrites)
	}
	return
}

func (cache *Cache) Clear() {
	for i := 0; i < 256; i++ {
		cache.locks[i].Lock()
		newSeg := newSegment(len(cache.segments[i].rb.data), i)
		cache.segments[i] = newSeg
		cache.locks[i].Unlock()
	}
	atomic.StoreInt64(&cache.hitCount, 0)
	atomic.StoreInt64(&cache.missCount, 0)
}

func (cache *Cache) ResetStatistics() {
	atomic.StoreInt64(&cache.hitCount, 0)
	atomic.StoreInt64(&cache.missCount, 0)
	for i := 0; i < 256; i++ {
		cache.locks[i].Lock()
		cache.segments[i].resetStatistics()
		cache.locks[i].Unlock()
	}
}
