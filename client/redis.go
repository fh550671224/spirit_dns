package client

import (
	"context"
	dns "github.com/fh550671224/spirit_dns_public"
	"log"
)

var RedisClient dns.RedisClient

func InitRedis() {
	err := RedisClient.InitRedis(context.Background(), "redis-service:6379")
	if err != nil {
		log.Printf("Warning: client.InitRedis err: %v\n", err)
	}
}
