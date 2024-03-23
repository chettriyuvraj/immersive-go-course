package cache

import "sync"

type CacheStats struct {
	hitRate                 float64
	untouchedElems          int
	readAvg                 int
	totalReads, totalWrites int
}

type CacheNode[V any] struct {
	val   V
	reads int
}

func (ce *CacheNode[V]) isUntouched() bool {
	return ce.reads == 0
}

type Cache[K comparable, V any] struct {
	cacheMap             map[K]CacheNode[V]
	size                 int
	mu                   sync.Mutex
	hits, misses, writes int
	// evictedWithoutTouch  int
}

/* Always initialize new cache using this function */
func NewCache[K comparable, V any](size int) *Cache[K, V] {
	return &Cache[K, V]{
		cacheMap: make(map[K]CacheNode[V]),
		size:     size,
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, exists := c.cacheMap[key]
	if !exists {
		c.misses++
		return getZero[V](), false
	}

	c.hits++
	node.reads++
	return node.val, true
}

func (c *Cache[K, V]) Put(key K, val V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheMap[key] = CacheNode[V]{val: val}
	c.size++
	c.writes++
}

func getZero[V any]() V {
	var result V
	return result
}
