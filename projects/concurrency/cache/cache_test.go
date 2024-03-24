package cache

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

/* TODO: How would you dynamically test caches of different types? */
func TestPutAndGet(t *testing.T) {
	cacheSize := 5
	cache := NewCache[int, int](cacheSize)
	rand := getRandGenerator()

	key, val := rand.Int(), rand.Int()
	cache.Put(key, val)
	cachedVal, isCached := cache.Get(key)
	require.Equal(t, true, isCached)
	require.Equal(t, val, cachedVal)
}

func TestPutAndGetLock(t *testing.T) {
	cacheSize := 5
	cache := NewCache[int, int](cacheSize)
	rand := getRandGenerator()

	t.Run("Get cache in parallel", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			cache.Get(rand.Int())
		}
	})

	t.Run("Put cache in parallel", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			cache.Put(rand.Int(), rand.Int())
		}
	})
}

func TestStatsUpdate(t *testing.T) {
	cacheSize, iterations := 500, 50
	rand := getRandGenerator()

	t.Run("Test cache hit count", func(t *testing.T) {
		t.Parallel()
		cache := NewCache[int, int](cacheSize)
		key, val := rand.Int(), rand.Int()
		cache.Put(key, val)
		for i := 0; i < iterations; i++ {
			cachedVal, exists := cache.Get(key)
			require.Equal(t, true, exists)
			require.Equal(t, val, cachedVal)
			require.Equal(t, i+1, cache.hits)
		}
	})

	t.Run("Test cache miss count", func(t *testing.T) {
		t.Parallel()
		cache := NewCache[int, int](cacheSize)
		key := rand.Int()
		for i := 0; i < iterations; i++ {
			_, exists := cache.Get(key)
			require.Equal(t, false, exists)
			require.Equal(t, i+1, cache.misses)
		}
	})

	t.Run("Test cache write count", func(t *testing.T) {
		t.Parallel()
		cache := NewCache[int, int](cacheSize)
		key, val := rand.Int(), rand.Int()
		for i := 0; i < iterations; i++ {
			cache.Put(key, val)
			require.Equal(t, i+1, cache.writes)
		}
	})
}

func TestAddNewNodeToHead(t *testing.T) {
	cacheSize := 5 /* Won't matter as we are not triggering cache.Put() */
	iterations := 50
	cache := NewCache[int, int](cacheSize)
	for i := 0; i < iterations; i++ {
		val := rand.Int()
		newNode := &CacheNode[int]{val: val}
		cache.AddNewNodeToHead(newNode)
		require.Same(t, newNode, cache.head)
	}
}

func TestAddExistingNodeToHead(t *testing.T) {
	type pair struct{ key, val int }
	tcs := []struct {
		name         string
		cacheElems   []pair
		addToHeadIdx int
		cacheSize    int
	}{
		/* Ensure cache elems and index are value WRT cacheSize */
		{
			name:         "single elem - add head to head",
			cacheElems:   []pair{{1, 2}},
			addToHeadIdx: 0,
			cacheSize:    5,
		},
		{
			name:         "two elems - add tail to head",
			cacheElems:   []pair{{1, 2}, {3, 4}},
			addToHeadIdx: 1,
			cacheSize:    5,
		},
		{
			name:         "multiple elems - add random to head",
			cacheElems:   []pair{{1, 2}, {3, 4}, {5, 6}, {7, 8}},
			addToHeadIdx: 2,
			cacheSize:    5,
		},
	}

	for _, tc := range tcs {
		cache := NewCache[int, int](tc.cacheSize)

		/* Add all elems to cache - put in reverse order => order of LRU will be same as the order of tc.cacheElems */
		for i := len(tc.cacheElems) - 1; i >= 0; i-- {
			elem := tc.cacheElems[i]
			cache.Put(elem.key, elem.val)
		}

		/* Grab existing node and capture snapshot of current params */
		existingNodeKey := tc.cacheElems[tc.addToHeadIdx].key
		existingNode, exists := cache.cacheMap[existingNodeKey]
		require.Equal(t, true, exists)
		prevBeforeUpdate, nextBeforeUpdate, tailBeforeUpdate, headBeforeUpdate := existingNode.prev, existingNode.next, cache.tail, cache.head

		/* Add existing node to head and compare with expected values */
		err := cache.AddExistingNodeToHead(existingNode)
		require.NoError(t, err)
		require.Equal(t, cache.head, existingNode)
		if prevBeforeUpdate != nil {
			require.Equal(t, nextBeforeUpdate, prevBeforeUpdate.next)
		}
		if nextBeforeUpdate != nil {
			require.Equal(t, prevBeforeUpdate, nextBeforeUpdate.prev)
		}
		if len(tc.cacheElems) > 1 {
			require.Equal(t, existingNode.next, headBeforeUpdate)
		}
		if tailBeforeUpdate == existingNode && len(tc.cacheElems) > 1 {
			require.Equal(t, cache.tail, prevBeforeUpdate)
		}
	}
}

func TestLRU(t *testing.T) {
	cacheSize := 5
	rand := getRandGenerator()

	t.Run("No removals until cache limit", func(t *testing.T) {
		cache := NewCache[int, int](cacheSize)
		for i := 0; i < cacheSize; i++ {
			lruNodeBeforePut, sizeBeforePut := cache.tail, cache.size

			key, val := generateIntNotInCache(cache), rand.Int()
			cache.Put(key, val)

			if lruNodeBeforePut == nil {
				continue
			}
			lruNodeAfterPut := cache.tail
			require.Equal(t, lruNodeAfterPut, lruNodeBeforePut)
			require.Equal(t, sizeBeforePut+1, cache.size)
		}
	})

	t.Run("LRU removed after cache limit", func(t *testing.T) {
		cache := NewCache[int, int](cacheSize)
		iterations := cacheSize * 10
		rand := getRandGenerator()

		/* Fill cache to the limit */
		for i := 0; i < cacheSize; i++ {
			key, val := generateIntNotInCache(cache), rand.Int()
			cache.Put(key, val)
		}
		/* LRU must be eliminated after each put */
		for i := 0; i < iterations; i++ {
			lruNodeBeforePut := cache.tail
			key, val := generateIntNotInCache(cache), rand.Int()
			cache.Put(key, val)
			if lruNodeBeforePut == nil {
				continue
			}
			lruNodeAfterPut := cache.tail
			require.NotSame(t, lruNodeAfterPut, lruNodeBeforePut)
		}
	})

}

/* Pseudo-random number with same seed */
func getRandGenerator() *rand.Rand {
	seed := 500
	src := rand.NewSource(int64(seed))
	return rand.New(src)
}

func generateIntNotInCache(cache *Cache[int, int]) int {
	key := rand.Int()
	_, exists := cache.Get(key)
	for exists {
		key = rand.Int()
		_, exists = cache.Get(key)
	}
	return key
}
