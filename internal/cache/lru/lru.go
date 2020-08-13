package lru

import (
	"container/list"
	"sync"
)

type cacheItem struct {
	key interface{}
	val interface{}
}

// LRU is simple LRU cache with doubly linked list under the hood.
type LRU struct {
	sync.RWMutex
	cap   int
	list  *list.List
	items map[interface{}]*list.Element
}

// New returns new LRU cache with given capacity.
func New(cap int) *LRU {
	return &LRU{
		cap:   cap,
		list:  list.New(),
		items: make(map[interface{}]*list.Element, cap),
	}
}

// Set just updates values or insert key/val.
func (c *LRU) Set(key interface{}, value interface{}) (deleted interface{}, updated bool) {
	c.Lock()
	defer c.Unlock()

	// If item found in cache, just update value
	if item, ok := c.items[key]; ok {
		item.Value = cacheItem{key, value}
		c.list.MoveToFront(item)
		c.items[key] = item
		return nil, true
	}

	// Remove back element of list, if list oversized
	var deletedValue interface{}
	if c.list.Len() >= c.cap {
		backItem := c.list.Back()
		c.list.Remove(backItem)
		delete(c.items, backItem.Value.(cacheItem).key)
		deletedValue = backItem.Value.(cacheItem).val
	}

	c.items[key] = c.list.PushFront(cacheItem{key, value})
	return deletedValue, false
}

// Get returns cached value and found flag by given key.
func (c *LRU) Get(key interface{}) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()
	item, ok := c.items[key]
	if ok {
		c.list.MoveToFront(item)
		return item.Value.(cacheItem).val, true
	}
	return nil, false
}

// Clear drops cache.
func (c *LRU) Clear() {
	c.Lock()
	defer c.Unlock()

	c.list = list.New()
	c.items = make(map[interface{}]*list.Element, c.cap)
}

// Len returns current cache size.
func (c *LRU) Len() int {
	c.RLock()
	defer c.RUnlock()
	return c.list.Len()
}
