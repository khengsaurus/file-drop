package database

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/khengsaurus/file-drop/server/consts"
)

type RedisClient struct {
	instance *redis.Client
}

func InitRedisClient() *RedisClient {
	var opts *redis.Options

	if consts.Local {
		fmt.Println("Redis config: local")
		redisAddress := fmt.Sprintf("%s:6379", os.Getenv("REDIS_URI_DEV"))
		opts = (&redis.Options{
			Addr:     redisAddress,
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

func (redisClient *RedisClient) CheckExists(ctx context.Context, key string) *redis.IntCmd {
	return redisClient.instance.Exists(ctx, key)
}

func (redisClient *RedisClient) GetShortestNewKey(ctx context.Context, key string) string {
	shortenedKey := ""

	// max length 5
	for i := 3; i <= 6; i++ {
		shortenedKey = key[:i]
		cmd := redisClient.CheckExists(ctx, shortenedKey)
		exists, err := cmd.Result()
		if err == nil && exists == 0 {
			return shortenedKey
		}
	}

	return redisClient.GetShortestNewKey(ctx, uuid.New().String())
}

func (redisClient *RedisClient) GetRedisValue(
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

func (redisClient *RedisClient) SetRedisValue(
	ctx context.Context,
	key string,
	value string,
) error {
	prefixedKey := fmt.Sprintf("%s_%s", consts.RedisKeyPrefix, key)
	return redisClient.instance.Set(
		ctx,
		prefixedKey,
		value,
		consts.RedisTTL,
	).Err()
}
