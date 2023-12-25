package consts

import (
	"os"
	"time"
)

type ContextKey string

var (
	Local          = os.Getenv("LOCAL") == "true"
	S3ClientKey    = ContextKey("s3_client")
	RedisClientKey = ContextKey("redis_client")
	RedisKeyPrefix = "FD1"
	RedisValDelim  = "___"
	RedisTTL       = 24 * time.Hour
)
