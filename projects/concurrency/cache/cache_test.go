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
	cachedVal, isCached := cache.Get(key) /* Note that Get returns a pointer */
	require.Equal(t, true, isCached)
	require.Equal(t, val, *cachedVal)
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
