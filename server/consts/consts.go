package consts

import (
	"os"
	"time"
)

type ContextKey string
type RedisResourceValue string /* [name, key, url, uploadedAt] joined by RedisValueDelim  */

var (
	Local          = os.Getenv("LOCAL") == "true"
	RedisCacheKey  = ContextKey("redis_cache")
	RedisClientKey = ContextKey("redis_client")
	RedisKeyPrefix = "FD1"
	RedisTTL       = 24 * time.Hour
	RedisValDelim  = "___"
	S3ClientKey    = ContextKey("s3_client")
	TokenKey       = []byte(os.Getenv("HASH_KEY"))
)
