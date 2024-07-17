// Package redis contains functions for using redis
package redis

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go-event-management/conf"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v8"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis.v8"
)

var (
	rdb        redis.UniversalClient
	rdbReplica redis.UniversalClient
)

// Init needs to be called first to set up the rdb and rdbReplica values
func Init(enableSSL bool, endpoint string, replicaEndpoint string) {
	options := redis.Options{
		Addr:         endpoint,
		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		PoolSize:     poolSize,
		PoolTimeout:  poolTimeout,
	}
	replicaOptions := redis.Options{
		Addr:         replicaEndpoint,
		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		PoolSize:     poolSize,
		PoolTimeout:  poolTimeout,
	}
	if enableSSL {
		options.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		replicaOptions.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}
	rdb = redistrace.NewClient(&options)
	rdbReplica = redistrace.NewClient(&replicaOptions)
}

// Set sets a string value with given ttl against a key
// 0 ttl means no expiry
func Set(key string, value string, ttl time.Duration) error {
	if ttl < 0 {
		ttl = 0
	}
	if key == "" {
		return ErrorEmptyKey
	}
	key += suffix
	err := rdb.Set(context.Background(), key, value, ttl)
	if err != nil {
		return err.Err()
	}
	return nil
}

// Keys returns all keys matching with pattern
func Keys(ctx context.Context, pattern string) ([]string, error) {
	if pattern == "" {
		return []string{}, ErrorEmptyPattern
	}

	ctx, cancelCtx := context.WithTimeout(ctx, contextTimeout)
	defer cancelCtx()
	val, err := rdbReplica.Keys(ctx, pattern).Result()
	if err != nil {
		return []string{}, err
	}

	result := make([]string, len(val))
	for _, key := range val {
		key = strings.TrimSuffix(key, suffix)
		result = append(result, key)
	}
	return result, nil
}

// GetTTL returns ttl for a key
func GetTTL(ctx context.Context, key string) (time.Duration, error) {
	key += suffix
	ttl, err := rdbReplica.TTL(context.TODO(), key).Result()
	return ttl, err
}

// SetStruct sets a struct object with given ttl against a key
func SetStruct(key string, obj interface{}, ttl time.Duration) error {
	valueBytes, err := json.Marshal(obj)
	if err != nil {
		return ErrorUnsupportedValue
	}
	return Set(key, string(valueBytes), ttl)
}

// SetStructWithLongTTL sets a struct object with a predefined long ttl value
func SetStructWithLongTTL(key string, obj interface{}) error {
	return SetStruct(key, obj, LongRedisTTL)
}

// Get returns value and error if any
// In case, key is not found it returns redis.Nil as error
func Get(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", ErrorEmptyKey
	}
	key += suffix
	ctx, cancelCtx := context.WithTimeout(ctx, contextTimeout)
	defer cancelCtx()
	val, err := rdbReplica.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// Delete deletes the key from redis
// It does not return an error if key is not found
func Delete(ctx context.Context, key string) error {
	if key == "" {
		return ErrorEmptyKey
	}
	key += suffix
	ctx, cancelCtx := context.WithTimeout(ctx, contextTimeout)
	defer cancelCtx()
	_, err := rdb.Del(ctx, key).Result()
	return err
}

func getSuffix(env string) string {
	if env == conf.ENV_PROD {
		return "-lendingapp"
	}
	return fmt.Sprintf("-lendingapp-%s", env)
}

func SetArgs(ctx context.Context, key string, value interface{}, a redis.SetArgs) error {
	_, err := rdb.SetArgs(ctx, key+suffix, value, a).Result()
	return err
}
