package cache //nolint:golint,stylecheck

import (
	"container/list"
	"sync"
)

type Key string
type Value interface{}

type Cache interface {
	Set(key Key, value Value) bool
	Get(key Key) (Value, bool)
	Clear()
	String() string
}

type cacheItem struct {
	key Key
	val Value
}

type lruCache struct {
	sync.Mutex
	cap                 int
	list                *list.List
	items               map[interface{}]*list.Element
	hitCount, missCount int
}

func NewCache(cap int) *lruCache {
	return &lruCache{
		cap:cap,
		list:  list.New(),
		items: make(map[interface{}]*list.Element, cap),
	}
}

//Set just updates values or insert key/val to front of short-term FIFO.
func (c *lruCache) Set(key Key, value Value) bool {
	c.Lock()
	defer c.Unlock()

	//If item found in cache, just update value
	if item, ok := c.items[key]; ok {
		item.Value = cacheItem{key, value}
		c.list.MoveToFront(item)
		c.items[key]=item
		return true
	}

	//Remove back element of list, if list oversized
	if c.list.Len() >= c.cap {
		backItem := c.list.Back()
		c.list.Remove(backItem)
		delete(c.items, backItem.Value.(cacheItem).key)
	}

	c.items[key] = c.list.PushFront(cacheItem{key, value})
	return false
}

func (c *lruCache) Get(key Key) (Value, bool) {
	c.Lock()
	defer c.Unlock()
	item, ok := c.items[key]
	if ok {
		c.list.MoveToFront(item)
		return item.Value.(cacheItem).val, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.Lock()
	defer c.Unlock()

	c.list = list.New()
	c.items = make(map[interface{}]*list.Element, c.cap)
}

func (c *lruCache) List() []Value {
	c.Lock()
	items := make([]Value, c.cap)
	defer c.Unlock()
	for e := c.list.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value)
	}
	return items
}
