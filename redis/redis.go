package redis

import (
    "time"
    "fmt"
    "context"
	"github.com/redis/go-redis/v9"
)

type RedisDatabase interface {
    Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
    Get(key string) *redis.StringCmd
}

type RedisClientImpl struct {
    client *redis.Client
}

func NewRedisClient() RedisDatabase {
    opt, _ := redis.ParseURL("redis://admin::iL83FvLpqHHJH!@redis-18916.c1.asia-northeast1-1.gce.cloud.redislabs.com:18916")
	client := redis.NewClient(opt)
    return &RedisClientImpl{client: client}
}

func (r *RedisClientImpl) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
    fmt.Println("Success")
    ctx := context.Background()
    return r.client.Set(ctx, key, value, expiration)
}

func (r *RedisClientImpl) Get(key string) *redis.StringCmd {
    ctx := context.Background()
    return r.client.Get(ctx, key)
}