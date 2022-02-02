package lru

import "container/list"

type Cache struct {
	maxBytes  int64 // 最大内存
	usedBytes int64 // 已用内存
	ll        *list.List
	cache     map[string]*list.Element
	// optional
	OnEvicted func(key string, value Value) // 记录被移除时的回调函数
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int // 计算所需的bytes
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (Value, bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele) // 设为最新
		kv := ele.Value.(*entry)
		return kv.value, ok
	}
	return nil, false
}

func (c *Cache) Set(key string, value Value) {
	if e, ok := c.cache[key]; ok {
		// update
		c.ll.MoveToFront(e)
		kv := e.Value.(*entry)
		c.usedBytes += int64(value.Len() - kv.value.Len())
		kv.value = value
	} else {
		// insert
		newEle := c.ll.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.cache[key] = newEle
		c.usedBytes += int64(len(key) + value.Len())
	}
	// 检测是否内存超了
	for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) RemoveOldest() {
	last := c.ll.Back() // 取最旧的元素
	if last != nil {
		c.ll.Remove(last)
		kv := last.Value.(*entry)
		delete(c.cache, kv.key)
		c.usedBytes -= int64(len(kv.key) + kv.value.Len())
		// 回调
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
