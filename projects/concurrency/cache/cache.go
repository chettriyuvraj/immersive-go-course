package cache

import "sync"

type CacheStats struct {
	hitRate                 float64
	untouchedElems          int
	readAvg                 int
	totalReads, totalWrites int
}

type CacheNode[V any] struct {
	val        V
	reads      int
	prev, next *CacheNode[V]
}

func (ce *CacheNode[V]) isUntouched() bool {
	return ce.reads == 0
}

type Cache[K comparable, V any] struct {
	cacheMap             map[K]*CacheNode[V]
	head, tail           *CacheNode[V] /* Currently implementing LRU eviction policy */
	size, limit          int
	mu                   sync.Mutex
	hits, misses, writes int
	// evictedWithoutTouch  int
}

/* Always initialize new cache using this function */
func NewCache[K comparable, V any](limit int) *Cache[K, V] {
	return &Cache[K, V]{
		cacheMap: make(map[K]*CacheNode[V]),
		limit:    limit,
	}
}

func (cache *Cache[K, V]) Get(key K) (V, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	node, exists := cache.cacheMap[key]
	if !exists {
		cache.misses++
		return getZero[V](), false
	}

	cache.AddExistingNodeToHead(node)

	cache.hits++
	node.reads++
	return node.val, true
}

func (cache *Cache[K, V]) Put(key K, val V) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	newNode := &CacheNode[V]{val: val}
	if cache.size < cache.limit {
		cache.AddNewNodeToHead(newNode)
		cache.size++
	} else {
		cache.RemoveLRUNode() /* No size change - new node replaces LRU node */
	}

	cache.cacheMap[key] = newNode
	cache.writes++
}

/* Add a node that already exists in cache to the head of the list */
func (cache *Cache[K, V]) AddExistingNodeToHead(node *CacheNode[V]) error {
	if cache.head == node {
		return nil
	}

	if node.prev != nil {
		node.prev.next = node.next
		if cache.tail == node {
			cache.tail = node.prev
		}
	}
	if node.next != nil {
		node.next.prev = node.prev
	}

	node.next = cache.head
	node.prev = nil

	cache.head.prev = node
	cache.head = node
	return nil
}

/* Note: does not modify cache size */
func (cache *Cache[K, V]) AddNewNodeToHead(node *CacheNode[V]) {
	if cache.head == nil && cache.tail == nil {
		cache.head = node
		cache.tail = node
		return
	}

	node.next = cache.head
	node.prev = nil

	cache.head.prev = node
	cache.head = node
}

/* Note: does not modify cache size */
func (cache *Cache[K, V]) RemoveLRUNode() *CacheNode[V] {
	if cache.tail == nil {
		return nil
	}

	oldTail := cache.tail
	tailPrev := cache.tail.prev
	if tailPrev != nil {
		tailPrev.next = nil
		cache.tail = tailPrev
	}
	return oldTail
}

func (cache *Cache[K, V]) GetLRUNode() *CacheNode[V] {
	return cache.tail
}

func getZero[V any]() V {
	var result V
	return result
}
