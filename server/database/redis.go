package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/khengsaurus/file-drop/server/consts"
)

type RedisClient struct {
	instance *redis.Client
}

func InitRedisClient() *RedisClient {
	var opts *redis.Options

	if consts.Local {
		fmt.Println("Redis config: local")
		opts = (&redis.Options{
			Addr:     os.Getenv("REDIS_URI_DEV"),
			Password: "",
			DB:       0,
		})
	} else {
		fmt.Println("Redis config: remote")
		opts = (&redis.Options{
			Addr:     os.Getenv("REDIS_URL"),
			Password: os.Getenv("REDIS_PW"),
			DB:       0,
		})
	}

	return &RedisClient{instance: redis.NewClient(opts)}
}

func GetRedisClient(ctx context.Context) (*RedisClient, error) {
	redisClient, ok := ctx.Value(consts.RedisClientKey).(*RedisClient)
	if !ok {
		return nil, fmt.Errorf("couldn't find %s in request context", consts.RedisClientKey)
	}
	return redisClient, nil
}

/* -------------------- Methods --------------------*/

func (redisClient *RedisClient) CheckExists(ctx context.Context, key string) (bool, error) {
	prefixedKey := fmt.Sprintf("%s_%s", consts.RedisKeyPrefix, key)
	res, err := redisClient.instance.Exists(ctx, prefixedKey).Result()
	return res == 1, err
}

func (redisClient *RedisClient) SetValue(
	ctx context.Context,
	key string,
	value string,
	ttl time.Duration,
) error {
	prefixedKey := fmt.Sprintf("%s_%s", consts.RedisKeyPrefix, key)
	return redisClient.instance.Set(
		ctx,
		prefixedKey,
		value,
		ttl,
	).Err()
}

func (redisClient *RedisClient) GetValue(
	ctx context.Context,
	key string,
) string {
	prefixedKey := fmt.Sprintf("%s_%s", consts.RedisKeyPrefix, key)
	redisValue, err := redisClient.instance.Get(ctx, prefixedKey).Result()
	if err == redis.Nil {
		return ""
	} else if err != nil {
		fmt.Println(err)
		return ""
	} else {
		return redisValue
	}
}

func (redisClient *RedisClient) DeleteKey(ctx context.Context, key string) {
	redisClient.instance.Del(ctx, key)
}
