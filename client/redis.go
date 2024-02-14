package client

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var RedisClient *redis.Client

func InitRedis() {
	// 创建Redis连接
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis服务器地址
		Password: "",               // Redis密码
		DB:       0,                // 使用默认的数据库
	})

	// 使用Ping命令检查连接是否正常
	ctx := context.Background()
	pong, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return
	}
	fmt.Println("Connected to Redis:", pong)
}

func CloseRedis() {
	err := RedisClient.Close()
	if err != nil {
		fmt.Println("Error closing connection:", err)
		return
	}
}

func SetRedis(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := RedisClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("setting key err: %v", err)
	}

	log.Printf("setting key(%v) success", key)
	return nil
}

func GetRedis(ctx context.Context, key string) (string, int, error) {
	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return "", 0, fmt.Errorf("getting value err: %v", err)
	}

	ttl, err := RedisClient.TTL(ctx, key).Result()
	if err != nil {
		return "", 0, fmt.Errorf("getting ttl err: %v", err)
	}

	return val, int(ttl), nil
}
