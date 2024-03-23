package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/* TODO: How would you dynamically test caches of different types? */
func TestCachePutAndGet(t *testing.T) {
	cache := NewCache[string, string](5)
	key, val := "testKey", "testVal"
	cache.Put(key, val)
	cachedVal, isCached := cache.Get(key)
	require.Equal(t, true, isCached)
	require.Equal(t, val, cachedVal)
}

func TestCachePutAndGetParallel(t *testing.T) {
	cache := NewCache[string, string](5)

	t.Run("Get cache in parallel", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			cache.Get("k1")
		}
	})

	t.Run("Put cache in parallel", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			cache.Put("k1", "v1")
		}
	})
}

func TestCacheStatsUpdate(t *testing.T) {
	cache := NewCache[string, string](500)
	t.Run("Test cache hit count", func(t *testing.T) {
		t.Parallel()
		iterations := 50
		key, val := "k1", "v1"
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
		iterations := 50
		key := "k2"
		cache := NewCache[string, string](5)
		for i := 0; i < iterations; i++ {
			_, exists := cache.Get(key)
			require.Equal(t, false, exists)
			require.Equal(t, i+1, cache.misses)
		}
	})

	t.Run("Test cache write count", func(t *testing.T) {
		t.Parallel()
		iterations := 50
		key, val := "k3", "v3"
		cache := NewCache[string, string](5)
		for i := 0; i < iterations; i++ {
			cache.Put(key, val)
			require.Equal(t, i+1, cache.writes)
		}

	})

}
