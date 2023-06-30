package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func CreateRedisDatabase1(user string, password string) (*redis.Client, error) {
	opt, _ := redis.ParseURL("redis://" + user + "::" + password + "@redis-18916.c1.asia-northeast1-1.gce.cloud.redislabs.com:18916")
	fmt.Println(opt)
	client := redis.NewClient(opt)
	return client, nil
}

func CreateRedisDatabase() (*redis.Client, error) {
	opt, _ := redis.ParseURL("redis://admin::iL83FvLpqHHJH!@redis-18916.c1.asia-northeast1-1.gce.cloud.redislabs.com:18916")
	client := redis.NewClient(opt)
	return client, nil
}

func Get(client *redis.Client, key string) (string, error) {
	ctx := context.Background()
	value, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func HGet(client *redis.Client, key string, field string) (string, error) {
	ctx := context.Background()
	value, err := client.HGet(ctx, key, field).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func HGet_int64(client *redis.Client, key string, field string) (int64, error) {
	ctx := context.Background()
	value, err := client.HGet(ctx, key, field).Int64()
	if err != nil {
		return 0, err
	}
	return value, nil
}

func Set(client *redis.Client, key string, value interface{}, expiration time.Duration) (string, error) {
	ctx := context.Background()

	_, err := client.Set(ctx, key, value, expiration).Result()
	if err != nil {
		return "", err
	}
	return key, nil
}

func HSet(client *redis.Client, key string, values ...interface{}) {
	ctx := context.Background()

	err := client.HSet(ctx, key, values).Err()
	if err != nil {
		panic(err)
	}
}

func SetExpired(client *redis.Client, key string) {
	ctx := context.Background()
	expiry := 10 * time.Minute
	err := client.Expire(ctx, key, expiry).Err()
	if err != nil {
		panic(err)
	}
}

func Keys(client *redis.Client, pattern string) ([]string, error) {
	ctx := context.Background()

	keys, err := client.Keys(ctx, pattern).Result()
	if err != nil {
		fmt.Println(err)
		return nil, err

	}
	return keys, nil
}

func Scan(client *redis.Client, key string) []string {
	var cursor uint64
	var keys []string
	// Fetch next batch of keys
	var err error
	keys, cursor, err = client.Scan(context.Background(), cursor, key, 0).Result()
	if err != nil {
		log.Fatalf("Failed to fetch keys: %v", err)
	}

	return keys
}
func HLen(client *redis.Client, key string) (int64, error) {
	// Fetch next batch of keys
	var err error
	res, err := client.HLen(context.Background(), key).Result()
	return res, err
}
func Del(client *redis.Client, keys string) (bool, error) {
	ctx := context.Background()

	err := client.Del(ctx, keys).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func ZAdd(client *redis.Client, key string, members ...redis.Z) error {
	ctx := context.Background()

	_, err := client.ZAdd(ctx, key, members...).Result()
	return err
}

func ZRevRangeWithScores(client *redis.Client, key string, start int64, stop int64) ([]redis.Z, error) {
	ctx := context.Background()

	results, err := client.ZRevRangeWithScores(ctx, key, start, stop).Result()
	return results, err
}
