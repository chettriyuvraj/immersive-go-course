package cache

type Cache[K comparable, V any] struct {
	cacheMap map[K]V
	size     int
}

/* Always initialize new cache using this function */
func NewCache[K comparable, V any](size int) *Cache[K, V] {
	return &Cache[K, V]{
		cacheMap: make(map[K]V),
		size:     size,
	}
}

func (c *Cache[K, V]) Get(key K) (*V, bool) {
	val, exists := c.cacheMap[key]
	if !exists {
		return nil, false
	}

	return &val, true
}

func (c *Cache[K, V]) Put(key K, val V) {
	c.cacheMap[key] = val
}
