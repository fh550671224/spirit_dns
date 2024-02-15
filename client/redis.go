package client

import (
	"context"
	dns "github.com/fh550671224/spirit_dns_public"
	"log"
)

var RedisClient dns.RedisClient

func InitRedis() {
	err := RedisClient.InitRedis(context.Background(), "localhost:6379")
	if err != nil {
		log.Fatalf("client.InitRedis err: %v", err)
	}
}
