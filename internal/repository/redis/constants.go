package redis

import (
	"errors"
	"go-event-management/conf"

	"time"

	redis "github.com/go-redis/redis/v8"
)

var suffix = getSuffix(conf.ENV)

var (
	Nil                   = redis.Nil
	ErrorEmptyKey         = errors.New("redis: key cannot be blank string")
	ErrorEmptyPattern     = errors.New("redis: pattern cannot be blank string")
	ErrorUnsupportedValue = errors.New("redis: unsupported value passed")
)

const (
	LongRedisTTL   = time.Hour * 24 * 7 // 1 week
	contextTimeout = 20 * time.Second
	dialTimeout    = 10 * time.Second
	readTimeout    = 30 * time.Second
	writeTimeout   = 30 * time.Second
	poolSize       = 20
	poolTimeout    = 30 * time.Second
)
