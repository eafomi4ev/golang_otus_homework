package hw04_lru_cache //nolint:golint,stylecheck

type Key string

type Cache interface {
	Set(key string, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key string) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                                 // Очистить кэш

}

type lruCache struct {
	capacity int
	queue    list
	items    map[string]*cacheItem
}

func (l *lruCache) Set(key string, value interface{}) bool {
	item, ok := l.items[key]

	if ok {
		l.queue.MoveToFront(item.value)
		item.value = l.queue.Front()
		item.value.Value = value
	} else {
		if l.queue.Len() == l.capacity {
			for cacheItemKey, cacheItemValue := range l.items {
				if cacheItemValue.value.Value == l.queue.Back().Value {
					delete(l.items, cacheItemKey)
					l.queue.Remove(l.queue.Back())
					break
				}
			}
		}

		newValue := l.queue.PushFront(value)
		l.items[key] = &cacheItem{
			key:   key,
			value: newValue,
		}
	}

	return ok
}

func (l *lruCache) Get(key string) (interface{}, bool) {
	item, ok := l.items[key]

	if ok {
		l.queue.MoveToFront(item.value)
	}

	var value interface{}
	if item != nil {
		value = item.value.Value
	}

	return value, ok
}

func (l *lruCache) Clear() {
	for key, item := range l.items {
		l.queue.Remove(item.value)
		delete(l.items, key)
	}
}

type cacheItem struct {
	key   string
	value *Item
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		items:    make(map[string]*cacheItem),
	}
}
