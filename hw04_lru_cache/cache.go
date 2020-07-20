package hw04_lru_cache //nolint:golint,stylecheck
import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	m        sync.Mutex
	cache    map[Key]*Item
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.m.Lock()
	defer l.m.Unlock()

	item, ok := l.cache[key]

	if ok {
		l.queue.MoveToFront(item)
		item.Value.(*cacheItem).value = value
	} else {
		if l.queue.Len() == l.capacity {
			lastQueueEl := l.queue.Back()
			delete(l.cache, lastQueueEl.Value.(*cacheItem).key)
			l.queue.Remove(lastQueueEl)
		}

		newCacheItem := &cacheItem{
			key:   key,
			value: value,
		}

		l.cache[key] = l.queue.PushFront(newCacheItem)
	}

	return ok
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.m.Lock()
	defer l.m.Unlock()

	item, ok := l.cache[key]

	if ok {
		l.queue.MoveToFront(item)
		return item.Value.(*cacheItem).value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.m.Lock()
	defer l.m.Unlock()

	for key, item := range l.cache {
		l.queue.Remove(item)
		delete(l.cache, key)
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		cache:    make(map[Key]*Item),
	}
}
