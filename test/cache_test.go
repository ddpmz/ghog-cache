package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/ddpmz/ghog-cache/cache"
)

func TestBatch(t *testing.T) {
	// t.Run("testMemory", testMemory)
	t.Run("testRedis", testRedis)
	// t.Run("testMemoryClear", testMemoryClear)
	// t.Run("testRedisClear", testRedisClear)
}

// 缓存使用内存测试
func testMemory(t *testing.T) {
	c := cache.New("prefix")
	ctx := context.Background()
	// tag can batch Management Cache
	_ = c.Set(ctx, "person", g.Map{"name": "John", "age": 10}, 0)
	v, _ := c.Get(ctx, "person")
	fmt.Println(v)
	// 按键删除
	_, _ = c.Remove(ctx, "person")
}

// 缓存使用redis测试
func testRedis(t *testing.T) {
	config := gredis.Config{
		Address: "127.0.0.1:6379",
		Db:      1,
	}
	ctx := context.Background()
	gredis.SetConfig(&config)
	c := cache.NewRedis("prefix")
	// tag can batch Management Cache
	_ = c.Set(ctx, "person", g.Map{"name": "John", "age": 10}, 0)
	v, _ := c.Get(ctx, "person")
	fmt.Println(v)
	// 按键删除
	_, _ = c.Remove(ctx, "person")
}

// 缓存标签使用内存测试
func testMemoryWithTag(t *testing.T) {
	c := cache.New("prefix")
	ctx := context.Background()
	// tag can batch Management Cache

	_ = c.Set(ctx, "person01", g.Map{"name": "John", "age": 10}, 0, "tag_person")
	_ = c.Set(ctx, "family01", g.Map{"address": "Kan Yun street"}, 0, "tag_family")
	_ = c.Set(ctx, "work01", g.Map{"unit": "lf"}, 0, "tag_work")

	_ = c.Set(ctx, "person02", g.Map{"name": "John", "age": 10}, 0, "tag_person")
	_ = c.Set(ctx, "family02", g.Map{"address": "Kan Yun street"}, 0, "tag_family")
	_ = c.Set(ctx, "work02", g.Map{"unit": "lf"}, 0, "tag_work")

	p1, _ := c.Get(ctx, "person01")
	p2, _ := c.Get(ctx, "person02")
	fmt.Println(p1, p2)
	// 缓存标签在读取缓存数据时和直接缓存读取一样，差别只在删除时可以批量删除
	// 比如要删除 person01和person02两组对应的缓存
	// 不使用tag时
	_, _ = c.Remove(ctx, "person01")
	_, _ = c.Remove(ctx, "person02")
	// 或
	c.Removes(ctx, []string{"person01", "person02"})
	// 使用缓存标签
	c.RemoveByTag(ctx, "tag_person") // 直接就可以删除该标签下的缓存("person01","person02")
	// 甚至可以批量删除标签
	c.RemoveByTags(ctx, []string{"tag_person", "tag_family"}) // 同时删除多组标签下的数据
}

// 缓存使用redis测试
func testRedisWithTag(t *testing.T) {
	config := gredis.Config{
		Address: "127.0.0.1:6379",
		Db:      1,
	}
	ctx := context.Background()
	gredis.SetConfig(&config)
	c := cache.NewRedis("prefix")
	// tag can batch Management Cache

	_ = c.Set(ctx, "person01", g.Map{"name": "John", "age": 10}, 0, "tag_person")
	_ = c.Set(ctx, "family01", g.Map{"address": "Kan Yun street"}, 0, "tag_family")
	_ = c.Set(ctx, "work01", g.Map{"unit": "lf"}, 0, "tag_work")

	_ = c.Set(ctx, "person02", g.Map{"name": "John", "age": 10}, 0, "tag_person")
	_ = c.Set(ctx, "family02", g.Map{"address": "Kan Yun street"}, 0, "tag_family")
	_ = c.Set(ctx, "work02", g.Map{"unit": "lf"}, 0, "tag_work")

	p1, _ := c.Get(ctx, "person01")
	p2, _ := c.Get(ctx, "person02")
	fmt.Println(p1, p2)
	// 缓存标签在读取缓存数据时和直接缓存读取一样，差别只在删除时可以批量删除
	// 比如要删除 person01和person02两组对应的缓存
	// 不使用tag时
	_, _ = c.Remove(ctx, "person01")
	_, _ = c.Remove(ctx, "person02")
	// 或
	c.Removes(ctx, []string{"person01", "person02"})
	// 使用缓存标签
	c.RemoveByTag(ctx, "tag_person") // 直接就可以删除该标签下的缓存("person01","person02")
	// 甚至可以批量删除标签
	c.RemoveByTags(ctx, []string{"tag_person", "tag_family"}) // 同时删除多组标签下的数据
}

// 缓存使用redis清空缓存测试
func testRedisClear(t *testing.T) {
	config := gredis.Config{
		Address: "127.0.0.1:6379",
		Db:      1,
	}
	ctx := context.Background()
	gredis.SetConfig(&config)
	c := cache.NewRedis("prefix")
	// tag can batch Management Cache
	_ = c.Set(ctx, "person", g.Map{"name": "John", "age": 10}, 0)
	v, _ := c.Get(ctx, "person")
	fmt.Println(v)
	// 清空
	_ = c.Clear(ctx)
	v, _ = c.Get(ctx, "person")
	fmt.Println(v)
}

// 缓存使用内存清空缓存测试
func testMemoryClear(t *testing.T) {
	c := cache.New("prefix")
	ctx := context.Background()
	// tag can batch Management Cache
	_ = c.Set(ctx, "person", g.Map{"name": "John", "age": 10}, 0)
	v, _ := c.Get(ctx, "person")
	fmt.Println(v)
	// 清空
	_ = c.Clear(ctx)
	v, _ = c.Get(ctx, "person")
	fmt.Println(v)
}
