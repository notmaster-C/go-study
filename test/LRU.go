package test

import (
	"container/list"
	"fmt"
	"sync"
)

// LRUCache 定义LRU缓存结构
type LRUCache struct {
	mu       sync.Mutex                    // 互斥锁，确保线程安全
	capacity int                           // 缓存容量
	items    map[interface{}]*list.Element // 存储缓存项的map
	l        *list.List                    // 双向链表，存储缓存项
}

// CacheItem 缓存项结构
type CacheItem struct {
	key   interface{}
	value interface{}
}

// NewLRUCache 创建一个新的LRU缓存
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[interface{}]*list.Element),
		l:        list.New(),
	}
}

// Get 获取缓存项
func (c *LRUCache) Get(key interface{}) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查key是否存在于map中
	elem, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// 如果存在，移动到链表尾部，表示最近使用
	c.l.MoveToBack(elem)
	return elem.Value.(*CacheItem).value, true
}

// Put 添加或更新缓存项
func (c *LRUCache) Put(key, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查key是否存在
	elem, exists := c.items[key]
	if exists {
		// 如果存在，更新值并移动到链表尾部
		c.l.MoveToBack(elem)
		elem.Value.(*CacheItem).value = value
		return
	}

	// 如果不存在，创建新的缓存项
	if c.l.Len() == c.capacity {
		// 如果已达到容量，删除最老的缓存项（链表头部）
		oldest := c.l.Front()
		c.l.Remove(oldest)
		delete(c.items, oldest.Value.(*CacheItem).key)
	}

	// 添加新的缓存项到链表尾部和map中
	elem = c.l.PushBack(&CacheItem{key: key, value: value})
	c.items[key] = elem
}

func main() {
	lru := NewLRUCache(2)

	lru.Put("a", 1)
	lru.Put("b", 2)
	fmt.Println(lru.Get("a")) // 输出 1

	lru.Put("c", 3)           // 淘汰 "b"
	fmt.Println(lru.Get("b")) // 输出 false

	lru.Put("a", 4)           // "a" 被更新，不是被淘汰
	fmt.Println(lru.Get("a")) // 输出 4
	fmt.Println(lru.Get("c")) // 输出 3
}
