package hw04lrucache

import (
	"errors"
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
}

type cacheItem struct {
	key   Key
	value interface{}
}

// NewCache создает новый *lruCache в интерфейсе Cache.
func NewCache(capacity int) (Cache, error) {
	if capacity > 0 {
		return &lruCache{
			capacity: capacity,
			queue:    NewList(),
			items:    make(map[Key]*ListItem, capacity),
		}, nil
	}
	return nil, ErrInvalidCapacity
}

var ErrInvalidCapacity = errors.New("invalid capacity: <= 0")

// Get получает элемент по ключу:
// - если элемент присутствует в словаре, то переместить элемент в начало очереди и вернуть его значение и true;
// - если элемента нет в словаре, то вернуть nil и false.
func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	var ok bool
	var listItem *ListItem
	if listItem, ok = l.items[key]; ok {
		l.queue.MoveToFront(listItem)
	} else {
		listItem = &ListItem{Value: &cacheItem{value: nil}}
	}
	return listItem.Value.(*cacheItem).value, ok
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
	var ok bool
	var listItem *ListItem
	if listItem, ok = l.items[key]; ok {
		listItem.Value.(*cacheItem).value = value
		l.queue.MoveToFront(listItem)
	} else {
		if l.capacity == l.queue.Len() { // буфер заполнен, нужно удаление одного cacheItem
			delete(l.items, l.queue.Back().Value.(*cacheItem).key)
			l.queue.Remove(l.queue.Back())
		}
		l.items[key] = l.queue.PushFront(&cacheItem{key: key, value: value})
	}
	return ok
}

// Clear очищает lruCache.
func (l *lruCache) Clear() {
	*l = lruCache{}
}
