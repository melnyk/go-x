package hotcold

import (
	"time"

	"go.melnyk.org/concept"
)

type hotcoldcache[K comparable, V any] struct {
	hot       concept.Cache[K, V]
	cold      concept.Cache[K, V]
	hotperiod time.Duration
}

func (c *hotcoldcache[K, V]) Get(key K) (V, bool) {
	if v, ok := c.hot.Get(key); ok {
		return v, ok
	}
	if v, ok := c.cold.Get(key); ok {
		c.hot.Set(key, v, c.hotperiod)
		return v, ok
	}
	return *new(V), false
}

func (c *hotcoldcache[K, V]) Set(key K, value V, ttl time.Duration) {
	c.hot.Set(key, value, ttl)
	c.cold.Set(key, value, ttl)
}

func (c *hotcoldcache[K, V]) Delete(key ...K) {
	c.hot.Delete(key...)
	c.cold.Delete(key...)
}

func (c *hotcoldcache[K, V]) Reset() error {
	errH := c.hot.Reset()
	errC := c.cold.Reset()
	if errH != nil {
		return errH
	}
	return errC
}

func NewCache[K comparable, V any](
	hot concept.Cache[K, V],
	cold concept.Cache[K, V],
	hotperiod time.Duration,
) concept.Cache[K, V] {
	return &hotcoldcache[K, V]{
		hot:       hot,
		cold:      cold,
		hotperiod: hotperiod,
	}
}
