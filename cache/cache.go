package cache

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
)

type GfCache struct {
	CachePrefix string // 缓存前缀
	cache       *gcache.Cache
	tagSetMux   sync.Mutex
}

// New 使用内存缓存
func New(cachePrefix string) *GfCache {
	cache := &GfCache{
		CachePrefix: cachePrefix,
		cache:       gcache.New(),
	}
	return cache
}

// NewRedis 使用redis缓存
func NewRedis(cachePrefix string) *GfCache {
	cache := New(cachePrefix)
	cache.cache.SetAdapter(gcache.NewAdapterRedis(g.Redis()))
	return cache
}

// 设置tag缓存的keys
func (c *GfCache) cacheTagKey(ctx context.Context, key interface{}, tag string) {
	tagKey := c.CachePrefix + c.setTagKey(tag)
	if tagKey != "" {
		tagValue := []interface{}{key}
		value, _ := c.cache.Get(ctx, tagKey)
		if value != nil {
			var keyValue []interface{}
			// 若是字符串
			if keyStr, ok := value.Val().(string); ok {
				js, err := gjson.DecodeToJson(keyStr)
				if err != nil {
					g.Log().Error(ctx, err)
					return
				}
				keyValue = gconv.SliceAny(js.Interface())
			} else {
				keyValue = gconv.SliceAny(value)
			}
			for _, v := range keyValue {
				if !reflect.DeepEqual(key, v) {
					tagValue = append(tagValue, v)
				}
			}
		}
		_ = c.cache.Set(ctx, tagKey, tagValue, 0)
	}
}

// 获取带标签的键名
func (c *GfCache) setTagKey(tag string) string {
	if tag != "" {
		tag = "tag_" + tag
	}
	return tag
}

// Set sets cache with <tagKey>-<value> pair, which is expired after <duration>.
// It does not expire if <duration> <= 0.
func (c *GfCache) Set(ctx context.Context, key string, value interface{}, duration time.Duration, tag ...string) (err error) {
	c.tagSetMux.Lock()
	if len(tag) > 0 {
		c.cacheTagKey(ctx, key, tag[0])
	}
	err = c.cache.Set(ctx, c.CachePrefix+key, value, duration)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	c.tagSetMux.Unlock()
	return
}

// SetIfNotExist sets cache with <tagKey>-<value> pair if <tagKey> does not exist in the cache,
// which is expired after <duration>. It does not expire if <duration> <= 0.
func (c *GfCache) SetIfNotExist(ctx context.Context, key string, value interface{}, duration time.Duration, tag string) (bool, error) {
	c.tagSetMux.Lock()
	defer c.tagSetMux.Unlock()
	c.cacheTagKey(ctx, key, tag)
	return c.cache.SetIfNotExist(ctx, c.CachePrefix+key, value, duration)
}

// Get returns the value of <tagKey>.
// It returns nil if it does not exist or its value is nil.
func (c *GfCache) Get(ctx context.Context, key string) (*gvar.Var, error) {
	return c.cache.Get(ctx, c.CachePrefix+key)
}

// GetOrSet returns the value of <tagKey>,
// or sets <tagKey>-<value> pair and returns <value> if <tagKey> does not exist in the cache.
// The tagKey-value pair expires after <duration>.
//
// It does not expire if <duration> <= 0.
func (c *GfCache) GetOrSet(ctx context.Context, key string, value interface{}, duration time.Duration, tag string) (*gvar.Var, error) {
	c.tagSetMux.Lock()
	defer c.tagSetMux.Unlock()
	c.cacheTagKey(ctx, key, tag)
	return c.cache.GetOrSet(ctx, c.CachePrefix+key, value, duration)
}

// GetOrSetFunc returns the value of <tagKey>, or sets <tagKey> with result of function <f>
// and returns its result if <tagKey> does not exist in the cache. The tagKey-value pair expires
// after <duration>. It does not expire if <duration> <= 0.
func (c *GfCache) GetOrSetFunc(ctx context.Context, key string, f gcache.Func, duration time.Duration, tag string) (*gvar.Var, error) {
	c.tagSetMux.Lock()
	defer c.tagSetMux.Unlock()
	c.cacheTagKey(ctx, key, tag)
	return c.cache.GetOrSetFunc(ctx, c.CachePrefix+key, f, duration)
}

// GetOrSetFuncLock returns the value of <tagKey>, or sets <tagKey> with result of function <f>
// and returns its result if <tagKey> does not exist in the cache. The tagKey-value pair expires
// after <duration>. It does not expire if <duration> <= 0.
//
// Note that the function <f> is executed within writing mutex lock.
func (c *GfCache) GetOrSetFuncLock(ctx context.Context, key string, f gcache.Func, duration time.Duration, tag string) (*gvar.Var, error) {
	c.tagSetMux.Lock()
	defer c.tagSetMux.Unlock()
	c.cacheTagKey(ctx, key, tag)
	return c.cache.GetOrSetFuncLock(ctx, c.CachePrefix+key, f, duration)
}

// Contains returns true if <tagKey> exists in the cache, or else returns false.
func (c *GfCache) Contains(ctx context.Context, key string) (bool, error) {
	return c.cache.Contains(ctx, c.CachePrefix+key)
}

// Remove deletes the <tagKey> in the cache, and returns its value.
func (c *GfCache) Remove(ctx context.Context, key string) (*gvar.Var, error) {
	return c.cache.Remove(ctx, c.CachePrefix+key)
}

// Removes deletes <keys> in the cache.
func (c *GfCache) Removes(ctx context.Context, keys []string) {
	keysWithPrefix := make([]interface{}, len(keys))
	for k, v := range keys {
		keysWithPrefix[k] = c.CachePrefix + v
	}
	_, _ = c.cache.Remove(ctx, keysWithPrefix...)
}

// RemoveByTag deletes the <tag> in the cache, and returns its value.
func (c *GfCache) RemoveByTag(ctx context.Context, tag string) {
	c.tagSetMux.Lock()
	tagKey := c.setTagKey(tag)
	// 删除tagKey 对应的 key和值
	keys, _ := c.Get(ctx, tagKey)
	if !keys.IsNil() {
		// 如果是字符串
		if keyStr, ok := keys.Val().(string); ok {
			js, err := gjson.DecodeToJson(keyStr)
			if err != nil {
				g.Log().Error(ctx, err)
				return
			}
			ks := gconv.SliceStr(js.Interface())
			c.Removes(ctx, ks)
		} else {
			ks := gconv.SliceStr(keys.Val())
			c.Removes(ctx, ks)
		}
	}
	_, _ = c.Remove(ctx, tagKey)
	c.tagSetMux.Unlock()
}

// RemoveByTags deletes <tags> in the cache.
func (c *GfCache) RemoveByTags(ctx context.Context, tag []string) {
	for _, v := range tag {
		c.RemoveByTag(ctx, v)
	}
}

// Clear deletes all in the cache.
func (c *GfCache) Clear(ctx context.Context) error {
	return c.cache.Clear(ctx)
}

// Data returns a copy of all tagKey-value pairs in the cache as map type.
func (c *GfCache) Data(ctx context.Context) (map[interface{}]interface{}, error) {
	return c.cache.Data(ctx)
}

// Keys returns all keys in the cache as slice.
func (c *GfCache) Keys(ctx context.Context) ([]interface{}, error) {
	return c.cache.Keys(ctx)
}

// KeyStrings returns all keys in the cache as string slice.
func (c *GfCache) KeyStrings(ctx context.Context) ([]string, error) {
	return c.cache.KeyStrings(ctx)
}

// Values returns all values in the cache as slice.
func (c *GfCache) Values(ctx context.Context) ([]interface{}, error) {
	return c.cache.Values(ctx)
}

// Size returns the size of the cache.
func (c *GfCache) Size(ctx context.Context) (int, error) {
	return c.cache.Size(ctx)
}
