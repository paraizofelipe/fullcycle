package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(host, port, password string) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
	})
	return &RedisStorage{client: client}
}

func (r *RedisStorage) Increment(ctx context.Context, key string, window time.Duration) (int64, error) {
	script := redis.NewScript(`
		local count = redis.call('INCR', KEYS[1])
		if count == 1 then
			redis.call('PEXPIRE', KEYS[1], ARGV[1])
		end
		return count
	`)
	result, err := script.Run(ctx, r.client, []string{key}, window.Milliseconds()).Int64()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (r *RedisStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	blockedKey := "blocked:" + key
	exists, err := r.client.Exists(ctx, blockedKey).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *RedisStorage) Block(ctx context.Context, key string, duration time.Duration) error {
	pipe := r.client.TxPipeline()
	pipe.Set(ctx, "blocked:"+key, 1, duration)
	pipe.Del(ctx, key)
	_, err := pipe.Exec(ctx)
	return err
}
