package lru

import "container/list"

type Cache struct {
	maxBytes int64 // cache最大容量
	nBytes   int64 // 当前cache大小
	ll       *list.List
	cache    map[string]*list.Element
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if val, ok := c.cache[key]; ok {
		c.ll.MoveToFront(val)
		return val.Value.(*entry).value, true
	}

	return
}

func (c *Cache) removeOldest() {
	back := c.ll.Back()
	c.ll.Remove(back)
	ele := back.Value.(*entry)
	c.nBytes = c.nBytes - int64(len(ele.key)) - int64(ele.value.Len())
	delete(c.cache, ele.key)
}

func (c *Cache) Set(key string, value Value) {
	if 0 == c.maxBytes {
		return
	}

	if val, ok := c.cache[key]; ok {
		c.ll.MoveToFront(val)
		c.nBytes = c.nBytes - int64(val.Value.(*entry).value.Len()) + int64(value.Len())
	} else {
		c.nBytes += int64(len(key)) + int64(value.Len())
		c.cache[key] = c.ll.PushFront(&entry{key: key, value: value})
	}

	for c.maxBytes > 0 && c.nBytes > c.maxBytes {
		c.removeOldest()
	}
}
