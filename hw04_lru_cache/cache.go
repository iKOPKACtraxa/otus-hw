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

type lruCache struct {
	// mu охраняет всю структуру
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
	itemsRev map[*ListItem]Key
}

type cacheItem struct {
	key   Key
	value interface{}
}

// NewCache создает новый *lruCache в интерфейсе Cache.
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		itemsRev: make(map[*ListItem]Key, capacity),
	}
}

// Get получает элемент по ключу:
// - если элемент присутствует в словаре, то переместить элемент в начало очереди и вернуть его значение и true;
// - если элемента нет в словаре, то вернуть nil и false.
func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	cacheItem := cacheItem{
		key:   key,
		value: nil,
	}
	var ok bool
	var listItem *ListItem
	if listItem, ok = l.items[cacheItem.key]; ok {
		l.queue.MoveToFront(listItem)
	} else {
		listItem = &ListItem{Value: nil}
	}
	return listItem.Value, ok
}

// Set добавляет элемент по ключу и значению:
// - если элемент присутствует в словаре, то обновить его значение и переместить элемент в начало очереди;
// - если элемента нет в словаре, то добавить в словарь и в начало очереди
//  (при этом, если размер очереди больше ёмкости кэша,
//  то необходимо удалить последний элемент из очереди и его значение из словаря);
// - возвращаемое значение - флаг, присутствовал ли элемент в кэше.
func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	cacheItem := cacheItem{
		key:   key,
		value: value,
	}
	var ok bool
	if _, ok = l.items[cacheItem.key]; ok {
		l.items[cacheItem.key].Value = cacheItem.value
		l.queue.MoveToFront(l.items[cacheItem.key])
	} else {
		if l.capacity == l.queue.Len() {
			// удаление из карт старого элемента
			keyToDel := l.itemsRev[l.queue.Back()]
			delete(l.itemsRev, l.queue.Back())
			delete(l.items, keyToDel)
			// удаление из очереди старого элемента
			l.queue.Remove(l.queue.Back())
		}
		newListItem := l.queue.PushFront(cacheItem.value)
		l.items[cacheItem.key] = newListItem
		l.itemsRev[newListItem] = cacheItem.key
	}
	return ok
}

// Clear очищает lruCache.
func (l *lruCache) Clear() {
	*l = lruCache{}
}
