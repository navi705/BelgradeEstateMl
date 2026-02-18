package main

import (
	"container/list"
	"sync"
)

type CacheEntry struct {
	Key   string
	Value []byte
}

type LRUCache struct {
	capacity  int
	items     map[string]*list.Element
	evictList *list.List
	lock      sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity:  capacity,
		items:     make(map[string]*list.Element),
		evictList: list.New(),
	}
}

func (c *LRUCache) Get(key string) ([]byte, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		return ent.Value.(*CacheEntry).Value, true
	}
	return nil, false
}

func (c *LRUCache) Put(key string, value []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		ent.Value.(*CacheEntry).Value = value
		return
	}

	ent := &CacheEntry{Key: key, Value: value}
	entry := c.evictList.PushFront(ent)
	c.items[key] = entry

	if c.evictList.Len() > c.capacity {
		c.removeOldest()
	}
}

func (c *LRUCache) removeOldest() {
	ent := c.evictList.Back()
	if ent != nil {
		c.evictList.Remove(ent)
		kv := ent.Value.(*CacheEntry)
		delete(c.items, kv.Key)
	}
}
