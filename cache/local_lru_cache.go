package cache

import (
	"github.com/hashicorp/golang-lru"
	"log"
)

type LRUCache struct {
	cache *lru.Cache
}

func NewLRUCache(cacheSize int) *LRUCache {
	c, err := lru.New(cacheSize)

	if err != nil {
		log.Fatal("failed to initialize cache", err)
	}

	return &LRUCache{
		cache: c,
	}
}

func (c *LRUCache) Set(key interface{}, data interface{}) {
	c.cache.Add(key, data)
}

func (c *LRUCache) Get(key interface{}) (interface{}, bool) {
	return c.cache.Get(key)
}