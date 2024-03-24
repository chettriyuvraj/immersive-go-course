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

// func TestLRU(t *Testing.T) {
// 	cacheSize := 5

// 	t.Run("Fill cache till limit") {

// 	}

// }

/* Pseudo-random number with same seed */
func getRandGenerator() *rand.Rand {
	seed := 500
	src := rand.NewSource(int64(seed))
	return rand.New(src)
}
