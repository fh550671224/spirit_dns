package service

import (
	"context"
	"github.com/fh550671224/spirit_dns_public"
	"log"
	"spiritDNS/client"
	"sync"
)

var mu sync.RWMutex

func StoreRedisCache(q dns.Question, answers []dns.RR) {
	if !client.RedisClient.IsOk() {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	err := client.RedisClient.StoreRedisCache(context.Background(), q, answers)
	if err != nil {
		log.Printf("Warning: client.RedisClient.StoreRedisCache err: %v\n", err)
	}
}

func GetRedisCache(q dns.Question) ([]dns.RR, bool) {
	if !client.RedisClient.IsOk() {
		return nil, false
	}

	mu.RLock()
	defer mu.RUnlock()

	answers, err := client.RedisClient.GetRedisCacheByKey(context.Background(), q)
	if err != nil {
		log.Printf("Warning: client.RedisClient.GetRedisCache err: %v\n", err)
		return nil, false
	}
	return answers, len(answers) > 0
}
