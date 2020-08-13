package lru

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
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		_, inCache := c.Set("aaa", 100)
		require.False(t, inCache)

		_, inCache = c.Set("bbb", 200)
		require.False(t, inCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		_, inCache = c.Set("aaa", 300)
		require.True(t, inCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("capacity", func(t *testing.T) {
		c := NewCache(2)
		_, _ = c.Set("a", 100)
		_, _ = c.Set("b", 200)
		deleted, _ := c.Set("c", 300)
		require.Equal(t, 100, deleted)

		v, ok := c.Get("a")
		require.False(t, ok)
		require.Nil(t, v)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(1)
		_, _ = c.Set("aaa", 100)
		c.Clear()

		_, ok := c.Get("aaa")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			c.Set(strconv.Itoa(i), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			c.Get(strconv.Itoa(rand.Intn(10000)))
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			c.Len()
		}
	}()

	wg.Wait()
}
