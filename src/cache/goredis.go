package cache

import (
	"awesomeProject/src/config"
	log "awesomeProject/src/logs"
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

//var ctx = context.Background()
//var redisCache *redis.Client

type RedisCache struct {
	ctx         context.Context
	redisClient *redis.Client
}

// Redis连接池
func NewRedisCache(config config.RedisConfig) *RedisCache {
	return &RedisCache{
		ctx: context.Background(),
		redisClient: redis.NewClient(&redis.Options{
			Addr:         config.Address,
			DialTimeout:  time.Duration(config.DialConnectTimeout) * time.Second,
			ReadTimeout:  time.Duration(config.DialReadTimeout) * time.Second,
			PoolSize:     config.MaxActive,
			PoolTimeout:  time.Duration(config.IdleTimeout) * time.Second,
			MinIdleConns: config.MaxIdle,
			IdleTimeout:  time.Duration(config.IdleTimeout) * time.Second,
		}),
	}
}

func (cache *RedisCache) Close() {
	err := cache.redisClient.Close()
	if err != nil {
		log.Logger.WithError(err).Error("close redis cache error!")
	}
}

// smembers实现
func (cache *RedisCache) SMembers(key string) ([]string, error) {
	result, err := cache.redisClient.SMembers(cache.ctx, key).Result()
	if err != nil {
		log.Logger.WithError(err).Error("smember error!")
		return nil, err
	}
	return result, err
}

// set 实现
func (cache *RedisCache) Set(key string, value interface{}, expireTime time.Duration) error {
	err := cache.redisClient.Set(cache.ctx, key, value, expireTime).Err()
	if err != nil {
		return err
	}
	return nil
}

// get 实现
func (cache *RedisCache) Get(key string) (result string, err error) {
	result, err = cache.redisClient.Get(cache.ctx, key).Result()
	return
}

// setnx 实现
func (cache *RedisCache) SetNX(key string, value interface{}, expiration time.Duration) error {
	err := cache.redisClient.SetNX(cache.ctx, key, value, expiration).Err()
	return err
}

// setex
func (cache *RedisCache) SetXX(key string, value interface{}, expiration time.Duration) error {
	err := cache.redisClient.SetXX(cache.ctx, key, value, expiration).Err()
	return err
}

// lpush
func (cache *RedisCache) LPush(key string, values ...interface{}) error {
	err := cache.redisClient.LPush(cache.ctx, key, values...).Err()
	return err
}

// lpop
func (cache *RedisCache) LPop(key string) (string, error) {
	result, err := cache.redisClient.LPop(cache.ctx, key).Result()
	return result, err
}

// sadd
func (cache *RedisCache) SAdd(key string, members ...interface{}) error {
	err := cache.redisClient.SAdd(cache.ctx, key, members...).Err()
	return err
}

// srem
func (cache *RedisCache) SRem(key string, members ...interface{}) (int64, error) {
	ret, err := cache.redisClient.SRem(cache.ctx, key, members...).Result()
	return ret, err
}

// HSet
func (cache *RedisCache) HSet(key string, values ...interface{}) error {
	err := cache.redisClient.HSet(cache.ctx, key, values...).Err()
	return err
}

// HGet
func (cache *RedisCache) HGet(key string, field string) (string, error) {
	result, err := cache.redisClient.HGet(cache.ctx, key, field).Result()
	return result, err
}

func (cache *RedisCache) HLen(key string) (int64, error) {
	return cache.redisClient.HLen(cache.ctx, key).Result()
}

// hdel
func (cache *RedisCache) HDel(key string, fields ...string) error {
	err := cache.redisClient.HDel(cache.ctx, key, fields...).Err()
	return err
}
