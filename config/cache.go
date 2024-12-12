package config

import (
	"sync"
	"time"
)

// cacheEntry 캐시 항목 구조체
type cacheEntry struct {
	value      string
	expiration int64
}

// Cache 캐시 구조체
type Cache struct {
	items map[string]cacheEntry
	mutex sync.RWMutex
}

var (
	cacheInstance *Cache
	once          sync.Once
)

// GetCache 캐시 인스턴스 반환
func GetCache() *Cache {
	once.Do(func() {
		cacheInstance = &Cache{
			items: make(map[string]cacheEntry),
		}
		go cacheInstance.cleanupExpiredItems(1 * time.Minute)
	})
	return cacheInstance
}

// Set 캐시에 항목 추가
func (c *Cache) Set(key, value string, duration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	} else {
		expiration = 0 // 0이면 만료되지 않습니다.
	}
	c.items[key] = cacheEntry{
		value:      value,
		expiration: expiration,
	}
}

// Get 캐시에서 항목 가져오기
func (c *Cache) Get(key string) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	entry, found := c.items[key]
	if !found {
		return ""
	}
	if entry.expiration > 0 && time.Now().UnixNano() > entry.expiration {
		return ""
	}
	return entry.value
}

// Delete 캐시에서 항목 삭제
func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.items, key)
}

// cleanupExpiredItems 만료된 항목 정리
func (c *Cache) cleanupExpiredItems(interval time.Duration) {
	for {
		time.Sleep(interval)
		c.mutex.Lock()
		for key, entry := range c.items {
			if entry.expiration > 0 && time.Now().UnixNano() > entry.expiration {
				delete(c.items, key)
			}
		}
		c.mutex.Unlock()
	}
}
