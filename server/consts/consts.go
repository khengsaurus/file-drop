package consts

import (
	"os"
)

type ContextKey string
type RedisResourceValue string /* [name, key, url, uploadedAt] joined by RedisValueDelim  */

var (
	Local          = os.Getenv("LOCAL") == "true"
	RedisCacheKey  = ContextKey("redis_cache")
	RedisKeyPrefix = "FD1"
	RedisValDelim  = "___"
	//
	TokenKey         = []byte(os.Getenv("HASH_KEY"))
	ClientCookieName = "X-FD-Client"
	//
	MySqlClientKey = ContextKey("mysql_client")
	RedisClientKey = ContextKey("redis_client")
	S3ClientKey    = ContextKey("s3_client")
	//
	WriteTries                = 5
	RedisValPlaceholderPrefix = "placeholder"
)
