package hw04_lru_cache //nolint:golint,stylecheck
import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

type lruCache struct {
	cap   int
	queue List
	items map[Key]*ListItem
	m     sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	lru := lruCache{cap: capacity}
	lru.Clear()
	return &lru
}

func (lru *lruCache) Set(itemKey Key, value interface{}) bool {
	lru.m.Lock()
	defer lru.m.Unlock()
	listItem, ok := lru.items[itemKey]
	if ok {
		cacheItem := listItem.Value.(*cacheItem)
		cacheItem.value = value
		lru.queue.MoveToFront(listItem)
	} else {
		if lru.cap == lru.queue.Len() { // delete last item
			lastListItem := lru.queue.Back()
			cacheItem := lastListItem.Value.(*cacheItem)
			delete(lru.items, cacheItem.key)
			lru.queue.Remove(lastListItem)
		}
		newItem := &cacheItem{key: itemKey, value: value}
		newListItem := lru.queue.PushFront(newItem)
		lru.items[itemKey] = newListItem
	}

	return ok
}

func (lru *lruCache) Get(itemKey Key) (value interface{}, isExist bool) {
	lru.m.Lock()
	defer lru.m.Unlock()
	listItem, ok := lru.items[itemKey]
	if ok {
		lru.queue.MoveToFront(listItem)
		value = listItem.Value.(*cacheItem).value
		isExist = true
	}
	return
}

func (lru *lruCache) Clear() {
	lru.m.Lock()
	defer lru.m.Unlock()
	lru.items = make(map[Key]*ListItem)
	lru.queue = NewList()
}
