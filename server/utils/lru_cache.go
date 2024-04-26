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

// Returns the newly added or updated *list.Element
func (c *LruCache) Put(key string, value string) *list.Element {
	c.lock.Lock()
	defer c.lock.Unlock()

	if item, exists := c.items[key]; exists {
		item.data = value
		c.queue.MoveToFront(item.keyPtr)
		return item.keyPtr
	} else {
		if len(c.items) == c.capacity {
			toDelete := c.queue.Back()
			c.queue.Remove(toDelete)
			delete(c.items, toDelete.Value.(string))
		}
		ele := &node{
			data:      value,
			keyPtr:    c.queue.PushFront(key),
			createdAt: time.Now(),
		}
		c.items[key] = ele
		return ele.keyPtr
	}
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
