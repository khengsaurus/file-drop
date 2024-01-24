package utils

import (
	"container/list"
	"sync"
	"time"
)

type node struct {
	data      string
	keyPtr    *list.Element
	createdAt time.Time
}

type LruCache struct {
	capacity int
	items    map[string]*node
	lock     *sync.Mutex
	queue    *list.List
}

func NewLruCache(
	capacity int,
	ttl time.Duration,
	clearInterval time.Duration,
) LruCache {
	lruCache := LruCache{
		capacity: capacity,
		items:    make(map[string]*node),
		lock:     &sync.Mutex{},
		queue:    list.New(),
	}
	go lruCache.SetClearInterval(ttl, clearInterval)
	return lruCache
}

func (c *LruCache) Put(key string, value string) *list.Element {
	c.lock.Lock()
	defer c.lock.Unlock()

	var toDelete *list.Element
	if item, exists := c.items[key]; exists {
		item.data = value
		c.items[key] = item
		c.queue.MoveToFront(item.keyPtr)
	} else {
		if c.capacity == len(c.items) {
			toDelete := c.queue.Back()
			c.queue.Remove(toDelete)
			delete(c.items, toDelete.Value.(string))
		}
		c.items[key] = &node{
			data:      value,
			keyPtr:    c.queue.PushFront(key),
			createdAt: time.Now(),
		}
	}
	return toDelete
}

func (c *LruCache) Get(key string) string {
	c.lock.Lock()
	defer c.lock.Unlock()

	item := c.items[key]
	if item != nil {
		c.queue.MoveToFront(item.keyPtr)
		return item.data
	}
	return ""
}

func (c *LruCache) SetClearInterval(ttl, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		func() {
			c.lock.Lock()
			defer c.lock.Unlock()

			var next *list.Element
			for curr := c.queue.Front(); curr != nil; curr = next {
				next = curr.Next()
				key := curr.Value.(string)
				if item, exists := c.items[key]; exists {
					if time.Since(item.createdAt) > ttl {
						c.queue.Remove(curr)
						delete(c.items, key)
					}
				}
			}
		}()
	}
}
