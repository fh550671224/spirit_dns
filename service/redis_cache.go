package service

import (
	"github.com/fh550671224/spirit_dns_public"
	"log"
	"spiritDNS/client"
	"sync"
)

var mu sync.RWMutex

func StoreRedisCache(key dns.Question, answers []dns.RR) {
	mu.Lock()
	defer mu.Unlock()

	err := client.RedisClient.StoreRedisCache(key, answers)
	if err != nil {
		log.Printf("Warning: client.RedisClient.StoreRedisCache err: %v\n", err)
	}
}

func GetRedisCache(key dns.Question) ([]dns.RR, bool) {
	mu.RLock()
	defer mu.RUnlock()

	answers, err := client.RedisClient.GetRedisCache(key)
	if err != nil {
		log.Printf("Warning: client.RedisClient.GetRedisCache err: %v\n", err)
		return nil, false
	}
	return answers, true
}
