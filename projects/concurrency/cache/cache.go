package cache

import (
	"sync"
)

type CacheStats struct {
	hitRate                 float64
	untouchedElems          int
	avgReads                float64 /* Only for in-cache elems */
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
	cacheMap                    map[K]*CacheNode[V]
	head, tail                  *CacheNode[V] /* Currently implementing LRU eviction policy */
	size, limit                 int
	mu                          *sync.Mutex
	hits, misses, reads, writes int /* Every access to cacheMap regardless of it being Get/Put is considered for these fields */
	// evictedWithoutTouch  int
}

/* Always initialize new cache using this function */
func NewCache[K comparable, V any](limit int) *Cache[K, V] {
	return &Cache[K, V]{
		cacheMap: make(map[K]*CacheNode[V]),
		limit:    limit,
		mu:       &sync.Mutex{},
	}
}

func (cache *Cache[K, V]) Get(key K) (V, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	node, exists := cache.cacheMap[key]
	cache.reads++
	if !exists {
		cache.misses++
		return getZero[V](), false
	}

	node.reads++
	cache.AddExistingNodeToHead(node)

	cache.hits++
	return node.val, true
}

func (cache *Cache[K, V]) Put(key K, val V) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	/* If node already exists in cache */
	cache.reads++
	existingNode, exists := cache.cacheMap[key]
	if exists {
		existingNode.reads++
		err := cache.AddExistingNodeToHead(existingNode)
		if err != nil {
			return err
		}
		existingNode.val = val
		cache.writes++
		cache.hits++
		return nil
	}

	/* If node doesn't exist */
	newNode := &CacheNode[V]{val: val}
	if cache.size < cache.limit {
		cache.size++
	} else {
		cache.RemoveLRUNode() /* No size change - new node replaces LRU node */
	}
	cache.AddNewNodeToHead(newNode)
	cache.cacheMap[key] = newNode
	cache.misses++
	cache.writes++
	return nil
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

/*
- expects node to exist in cache, behaviour undefined if it does not
- does not modify cache.size and does not account for cache.size i.e. it will exceed cache.size if not called responsibly
*/
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

/*
- expects node to exist in cache, behaviour undefined if it does not
- does not modify cache.size and does not account for cache.size i.e. it will exceed cache.size if not called responsibly
*/
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

func (cache *Cache[K, V]) GetStats() CacheStats {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	stats := CacheStats{}

	/* Gather params to compute stats */
	totalCalls := cache.reads + cache.writes
	untouchedNodes, totalReadsForCachedNodes := 0, 0
	for _, node := range cache.cacheMap {
		if node.isUntouched() {
			untouchedNodes++
		}
		totalReadsForCachedNodes += node.reads
	}

	/* Compute stats */
	stats.hitRate = float64(cache.hits) / float64(totalCalls)
	stats.untouchedElems = untouchedNodes
	stats.avgReads = float64(totalReadsForCachedNodes) / float64(cache.reads)
	stats.totalReads = cache.reads
	stats.totalWrites = cache.writes

	return stats
}

func getZero[V any]() V {
	var result V
	return result
}
