package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mutex    sync.Mutex
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if item, exist := lru.items[key]; exist {
		item.Value = &cacheItem{
			key:   key,
			value: value,
		}

		lru.queue.MoveToFront(item)

		return true
	}

	if lru.queue.Len()+1 > lru.capacity {
		lastItem := lru.queue.Back()

		delete(lru.items, lastItem.Value.(*cacheItem).key)

		lru.queue.Remove(lastItem)
	}

	newItem := lru.queue.PushFront(&cacheItem{
		key:   key,
		value: value,
	})

	lru.items[key] = newItem

	return false
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if item, exist := lru.items[key]; exist {
		lru.queue.MoveToFront(item)

		return item.Value.(*cacheItem).value, true
	}

	return nil, false
}

func (lru *lruCache) Clear() {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	lru.queue = NewList()
	lru.items = make(map[Key]*ListItem, lru.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mutex:    sync.Mutex{},
	}
}
