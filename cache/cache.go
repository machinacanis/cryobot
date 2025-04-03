package cache

import (
	"github.com/coocood/freecache"
	"strconv"
)

type CryoCache struct {
	cacheName string           // 缓存名称
	expire    int              // 缓存过期时间（秒）
	cache     *freecache.Cache // Freecache 实例
}

// NewCryoCache 创建一个新的 CryoCache 缓存管理器实例
func NewCryoCache(cacheName string, cacheSize int, expire ...int) *CryoCache {
	if expire == nil {
		expire = append(expire, 60) // 默认过期时间为 60 秒
	}
	if cacheSize <= 0 {
		cacheSize = 128 // 默认缓存大小为 128MB
	}
	return &CryoCache{
		cacheName: cacheName,
		cache:     freecache.NewCache(cacheSize * 1024 * 1024),
		expire:    expire[0],
	}
}

// Set 设置缓存值
func (c *CryoCache) Set(key string, value []byte, expire ...int) error {
	if expire == nil {
		expire = append(expire, c.expire) // 使用默认过期时间
	}
	err := c.cache.Set([]byte(key), value, expire[0])
	return err
}

// Get 获取缓存值
func (c *CryoCache) Get(key string) ([]byte, error) {
	value, err := c.cache.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	return value, nil
}

// Del 删除缓存值
func (c *CryoCache) Del(key string) bool {
	ok := c.cache.Del([]byte(key))
	if !ok {
		return false
	}
	return true
}

// Clear 清除缓存
func (c *CryoCache) Clear() {
	c.cache.Clear()
}

// GetCacheName 获取缓存名称
func (c *CryoCache) GetCacheName() string {
	return c.cacheName
}

// GetExpire 获取缓存过期时间
func (c *CryoCache) GetExpire() int {
	return c.expire
}

// SetExpire 设置缓存过期时间
func (c *CryoCache) SetExpire(expire int) {
	c.expire = expire
}

// GetCache 获取 Freecache 实例
func (c *CryoCache) GetCache() *freecache.Cache {
	return c.cache
}

// SetString 设置字符串缓存值
func (c *CryoCache) SetString(key string, value string, expire ...int) error {
	err := c.Set(key, []byte(value), expire...)
	if err != nil {
		return err
	}
	return nil
}

// GetString 获取字符串缓存值
func (c *CryoCache) GetString(key string) (string, error) {
	value, err := c.Get(key)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

// SetInt 设置整数缓存值
func (c *CryoCache) SetInt(key string, value int, expire ...int) error {
	err := c.Set(key, []byte(string(value)), expire...)
	if err != nil {
		return err
	}
	return nil
}

// GetInt 获取整数缓存值
func (c *CryoCache) GetInt(key string) (int, error) {
	value, err := c.Get(key)
	if err != nil {
		return 0, err
	}
	intValue, err := strconv.Atoi(string(value))
	if err != nil {
		return 0, err
	}
	return intValue, nil
}

// SetBool 设置布尔值缓存
func (c *CryoCache) SetBool(key string, value bool, expire ...int) error {
	var boolValue int
	if value {
		boolValue = 1
	} else {
		boolValue = 0
	}
	err := c.Set(key, []byte(strconv.Itoa(boolValue)), expire...)
	if err != nil {
		return err
	}
	return nil
}

// GetBool 获取布尔值缓存
func (c *CryoCache) GetBool(key string) (bool, error) {
	value, err := c.Get(key)
	if err != nil {
		return false, err
	}
	boolValue, err := strconv.Atoi(string(value))
	if err != nil {
		return false, err
	}
	return boolValue == 1, nil
}

// SetFloat 设置浮点数缓存值
func (c *CryoCache) SetFloat(key string, value float64, expire ...int) error {
	err := c.Set(key, []byte(strconv.FormatFloat(value, 'f', -1, 64)), expire...)
	if err != nil {
		return err
	}
	return nil
}

// GetFloat 获取浮点数缓存值
func (c *CryoCache) GetFloat(key string) (float64, error) {
	value, err := c.Get(key)
	if err != nil {
		return 0, err
	}
	floatValue, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return 0, err
	}
	return floatValue, nil
}
