package hw04lrucache

import (
	"errors"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		a, err := NewCache(0)
		require.Nil(t, a)
		require.Truef(t, errors.Is(err, ErrInvalidCapacity), "actual error %q", err)

		c, err := NewCache(10)
		require.Truef(t, errors.Is(err, nil), "actual error %q", err)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c, _ := NewCache(5)

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

		c.Clear()
		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic1", func(t *testing.T) {
		c, _ := NewCache(3)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		wasInCache = c.Set("ccc", 300)
		require.False(t, wasInCache)

		wasInCache = c.Set("ddd", 400)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		val, ok = c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, 400, val)
	})
	t.Run("purge logic2", func(t *testing.T) {
		c, _ := NewCache(3)
		// добавляем 3 элемента
		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		wasInCache = c.Set("ccc", 300)
		require.False(t, wasInCache)
		// обновляем их Set'ами и Get'ами, при этом ааа будет последним обновленным
		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		wasInCache = c.Set("ccc", 333)
		require.True(t, wasInCache)

		wasInCache = c.Set("bbb", 222)
		require.True(t, wasInCache)
		// добавляем 4 элемент, он вытеснит "старый" ааа
		wasInCache = c.Set("ddd", 444)
		require.False(t, wasInCache)
		// проверяем, что старый ааа удален
		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
		// проверяем, что оставшиеся в порядке
		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 222, val)

		val, ok = c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, 333, val)

		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, 444, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	// t.Skip() // Remove me if task with asterisk completed.

	c, _ := NewCache(10)
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
