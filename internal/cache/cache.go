package cache

import (
	"sync"
	"time"
)

type Item struct {
	URL  string
	Name string
	at   time.Time
}

type Cache struct {
	mu    sync.RWMutex
	ttl   time.Duration
	items map[string]Item
}

func New(ttl time.Duration) *Cache {
	return &Cache{
		ttl:   ttl,
		items: make(map[string]Item),
	}
}

func (c *Cache) Get(key string) (Item, bool) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		return Item{}, false
	}

	if time.Since(item.at) > c.ttl {
		c.Delete(key)
		return Item{}, false
	}

	return item, true
}

func (c *Cache) Set(key, faviconUrl, name string) {
	c.mu.Lock()
	c.items[key] = Item{URL: faviconUrl, Name: name, at: time.Now()}
	c.mu.Unlock()
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}
