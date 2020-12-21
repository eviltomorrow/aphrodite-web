package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/eviltomorrow/aphrodite-base/zlog"
)

// RedisDSN redis dsn
var RedisDSN string

// Redis redis
var Redis *redis.Client

// BuildRedis build redis
func BuildRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     RedisDSN,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if status := rdb.Ping(ctx); status.Err() != nil {
		zlog.Fatal("Build redis connection failure", zap.Error(status.Err()))
	}

	Redis = rdb
}

// CloseRedis close redis
func CloseRedis() error {
	zlog.Info("Close redis connection", zap.String("dsn", RedisDSN))

	if Redis != nil {
		return Redis.Close()
	}
	return nil
}
