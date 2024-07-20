package test

import (
	"fmt"
	"sort"
	"sync"
)

// 一. 请使用下面结构体实现SortMap, 请实现各个接口, 并添加测试用例
type SortMap[K comparable, V any] struct {
	data   map[K]V
	keys   []K
	rwlock sync.RWMutex
}

// NewSortMap 创建SortMap对象
func NewSortMap[K comparable, V any]() *SortMap[K, V] {
	return &SortMap[K, V]{data: make(map[K]V)}
}

// SetIfNotExist 设置key:val的映射, 如果key不存在才算设置成功
func (self *SortMap[K, V]) SetIfNotExist(key K, val V) bool {
	self.rwlock.Lock()
	defer self.rwlock.Unlock()
	_, exists := self.data[key]
	if !exists {
		self.data[key] = val
		return true
	}
	return false
}

// Search 搜索key
func (self *SortMap[K, V]) Search(key K) (V, bool) {
	self.rwlock.RLock()
	defer self.rwlock.RUnlock()
	val, exists := self.data[key]
	return val, exists
}

// Removes 删除多个key, 并返回删除的内容
func (self *SortMap[K, V]) Removes(keys []K) map[K]V {
	self.rwlock.Lock()
	defer self.rwlock.Unlock()
	removed := make(map[K]V)
	for _, key := range keys {
		if val, exists := self.data[key]; exists {
			removed[key] = val
			delete(self.data, key)
		}
	}
	return removed
}

// Clone 深度克隆一个SortMap对象
func (self *SortMap[K, V]) Clone() *SortMap[K, V] {
	self.rwlock.RLock()
	defer self.rwlock.RUnlock()
	clonedData := make(map[K]V, len(self.data))
	for k, v := range self.data {
		clonedData[k] = v
	}
	return &SortMap[K, V]{data: clonedData}
}

// RLockIterator 读锁定迭代整个对象
func (self *SortMap[K, V]) RLockIterator(f func(key K, val V) bool) {
	self.rwlock.RLock()
	defer self.rwlock.RUnlock()
	for k, v := range self.data {
		if !f(k, v) {
			break
		}
	}
}

// SortKey 对key进行排序, 返回排序后的val切片
func (self *SortMap[K, V]) SortKey(f func(one K, two K) bool) []V {
	keys := make([]K, 0, len(self.data))
	values := make([]V, 0, len(self.data))

	self.rwlock.RLock()
	defer self.rwlock.RUnlock()

	for k := range self.data {
		keys = append(keys, k)
	}

	// 使用提供的比较函数进行排序
	sort.SliceStable(keys, func(i, j int) bool {
		return f(keys[i], keys[j])
	})

	for _, k := range keys {
		values = append(values, self.data[k])
	}
	return values
}

func TestSortMap() {
	sm := NewSortMap[int, string]()

	// 测试用例
	sm.SetIfNotExist(1, "one")
	sm.SetIfNotExist(2, "two")
	sm.SetIfNotExist(1, "should not change") // 测试重复设置

	// fmt.Println("Search for key 1:", sm.Search(1)) // 应该返回 ("one", true)

	removed := sm.Removes([]int{2})
	fmt.Println("Removed:", removed) // 应该返回 map[2:two]

	clone := sm.Clone()
	fmt.Println("Cloned data:", clone.data) // 克隆的数据

	sm.RLockIterator(func(key int, val string) bool {
		fmt.Println("Iterator:", key, val)
		return true
	})

	// 测试SortKey
	sortedValues := sm.SortKey(func(one, two int) bool {
		return one < two
	})
	fmt.Println("Sorted values:", sortedValues) // 应该按key排序
}

// 二. 策划需要一套推广系统. 要求如下
/*
1. 每个用户一个整形的用户ID. 例如11223344. 要求全球唯一编号
2. 每个用户一个邀请码. 例如 f6Ds3f
3. 用户A可以将自己的邀请码发给自己的朋友B, B注册后生成自己的邀请码. B可以邀请C,D. D邀请F
4. 产品运营过程中, 需要在程序内存中实时查询任何一个用户的邀请关系网.
4.1  统计每个用户下级用户有多少人. 和用户ID列表. (分层统计人数和用户ID列表)
4.2  用户D使用B的邀请码注册成功后, 需要发送奖品给B和A.

请描述设计思路 (不使用数据库查询, 全部为内存管理和查询)
*/
