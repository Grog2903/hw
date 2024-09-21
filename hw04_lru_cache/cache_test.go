package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", 100)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		c.Clear()

		_, ok = c.Get("aaa")
		require.False(t, ok)
	})

	t.Run("remove by capacity", func(t *testing.T) {
		c := NewCache(3)

		c.Set("k1", 1)
		c.Set("k2", 2)
		c.Set("k3", 3)

		val, ok := c.Get("k1")
		require.True(t, ok)
		require.Equal(t, 1, val)

		val, ok = c.Get("k2")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = c.Get("k3")
		require.True(t, ok)
		require.Equal(t, 3, val)

		c.Set("k4", 4)

		_, ok = c.Get("k1")
		require.False(t, ok)

		val, ok = c.Get("k2")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = c.Get("k3")
		require.True(t, ok)
		require.Equal(t, 3, val)

		val, ok = c.Get("k4")
		require.True(t, ok)
		require.Equal(t, 4, val)
	})

	t.Run("remove by usage", func(t *testing.T) {
		c := NewCache(3)

		c.Set("k1", 1)
		c.Set("k2", 2)
		c.Set("k3", 3)

		val, ok := c.Get("k1")
		require.True(t, ok)
		require.Equal(t, 1, val)

		val, ok = c.Get("k2")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = c.Get("k3")
		require.True(t, ok)
		require.Equal(t, 3, val)

		c.Set("k1", 11)

		val, ok = c.Get("k1")
		require.True(t, ok)
		require.Equal(t, 11, val)

		c.Set("k4", 4)

		_, ok = c.Get("k2")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
