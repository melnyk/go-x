package hotcold

import (
	"testing"
	"time"

	conceptcache "go.melnyk.org/concept/cache"
)

type cacheinmap[K comparable, V any] map[K]V

func (c cacheinmap[K, V]) Set(key K, value V, _ time.Duration) {
	c[key] = value
}
func (c cacheinmap[K, V]) Get(key K) (V, bool) {
	v, ok := c[key]
	return v, ok
}
func (c cacheinmap[K, V]) Delete(key ...K) {
	for _, k := range key {
		delete(c, k)
	}
}
func (c cacheinmap[K, V]) Reset() error {
	for key := range c {
		delete(c, key)
	}
	return nil
}

func TestNewCache(t *testing.T) {
	hot := &conceptcache.Empty[int, int]{}
	cold := &conceptcache.Empty[int, int]{}

	cache := NewCache[int, int](hot, cold, time.Minute).(*hotcoldcache[int, int])

	if cache.hot != hot {
		t.Errorf("Expected hot cache to be %v, got %v", hot, cache.hot)
	}
	if cache.cold != cold {
		t.Errorf("Expected cold cache to be %v, got %v", cold, cache.cold)
	}
	if cache.hotperiod != time.Minute {
		t.Errorf("Expected ttl to be %v, got %v", time.Minute, cache.hotperiod)
	}
}

func TestCacheSet(t *testing.T) {
	var hot cacheinmap[string, int] = map[string]int{
		"test": 1,
	}
	var cold cacheinmap[string, int] = map[string]int{
		"test":  100,
		"test2": 2,
		"test3": 3,
	}

	cache := NewCache[string, int](hot, cold, time.Minute)

	cache.Set("test", 5, time.Minute)
	if hot["test"] != 5 {
		t.Errorf("Expected hot cache to be %v, got %v", 5, hot["test"])
	}
	if cold["test"] != 5 {
		t.Errorf("Expected cold cache to be %v, got %v", 5, cold["test"])
	}

	cache.Set("testN", 1111, time.Minute)
	if hot["testN"] != 1111 {
		t.Errorf("Expected hot cache to be %v, got %v", 1111, hot["testN"])
	}
	if cold["testN"] != 1111 {
		t.Errorf("Expected cold cache to be %v, got %v", 1111, cold["testN"])
	}
}

func TestCacheGet(t *testing.T) {
	var hot cacheinmap[string, int] = map[string]int{
		"test": 1,
	}
	var cold cacheinmap[string, int] = map[string]int{
		"test":  100,
		"test2": 2,
	}
	cache := NewCache[string, int](hot, cold, time.Minute)

	if v, ok := cache.Get("test"); !ok || v != 1 {
		t.Errorf("Expected cache to contain test=1, got %v", v)
	}

	if v, ok := cache.Get("test2"); !ok || v != 2 {
		t.Errorf("Expected cache to contain test2=2, got %v", v)
	}

	if hot["test2"] != 2 {
		t.Errorf("Expected hot cache to contain test2=2, got %v", hot["test2"])
	}

	if v, ok := cache.Get("test3"); ok {
		t.Errorf("Expected cache to not contain test3, got %v", v)
	}
}

func TestCacheDelete(t *testing.T) {
	var hot cacheinmap[string, int] = map[string]int{
		"test": 1,
	}
	var cold cacheinmap[string, int] = map[string]int{
		"test":  100,
		"test2": 2,
	}
	cache := NewCache[string, int](hot, cold, time.Minute)

	cache.Delete("test")
	if _, ok := hot["test"]; ok {
		t.Errorf("Expected hot cache to not contain test, got %v", hot["test"])
	}

	if _, ok := cold["test"]; ok {
		t.Errorf("Expected cold cache to not contain test, got %v", cold["test"])
	}

	cache.Delete("test2")
	if _, ok := hot["test2"]; ok {
		t.Errorf("Expected hot cache to not contain test2, got %v", hot["test2"])
	}
	if _, ok := cold["test2"]; ok {
		t.Errorf("Expected cold cache to not contain test2, got %v", cold["test2"])
	}
}

func TestCacheReset(t *testing.T) {
	var hot cacheinmap[string, int] = map[string]int{
		"test": 1,
	}
	var cold cacheinmap[string, int] = map[string]int{
		"test":  100,
		"test2": 2,
	}

	cache := NewCache[string, int](hot, cold, time.Minute)

	cache.Reset()

	if len(hot) != 0 {
		t.Errorf("Expected hot cache to be empty, got %v", hot)
	}
	if len(cold) != 0 {
		t.Errorf("Expected cold cache to be empty, got %v", cold)
	}

	a := &conceptcache.Empty[string, int]{}
	cache = NewCache[string, int](a, cold, time.Minute)

	err := cache.Reset()

	if err == nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
