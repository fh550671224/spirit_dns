package service

import (
	"context"
	"encoding/json"
	"github.com/fh550671224/spirit_dns_public"
	"log"
	"math"
	"spiritDNS/client"
	"sync"
	"time"
)

type wrappedObj struct {
	Type    uint16
	Payload interface{}
}

var mu sync.RWMutex

func StoreRedisCache(key dns.Question, answers []dns.RR) {
	var ttl uint32
	ttl = math.MaxUint32
	for _, a := range answers {
		if ttl > a.Header().Ttl {
			ttl = a.Header().Ttl
		}
	}

	mu.Lock()
	defer mu.Unlock()

	keyStr, err := json.Marshal(key)
	if err != nil {
		log.Printf("marshaling key err: %v", err)
		return
	}

	var woList []wrappedObj
	for _, a := range answers {
		woList = append(woList, wrappedObj{
			Type:    a.Header().Rrtype,
			Payload: a,
		})
	}

	woListStr, err := json.Marshal(&woList)
	if err != nil {
		log.Printf("marshaling answers err: %v", err)
	}

	err = client.SetRedis(context.Background(), string(keyStr), string(woListStr), time.Duration(ttl)*time.Second)
	if err != nil {
		log.Printf("client.SetRedis err: %v", err)
		return
	}
}

func GetRedisCache(key dns.Question) (AnswerList, bool) {
	mu.RLock()
	defer mu.RUnlock()

	keyStr, err := json.Marshal(key)
	if err != nil {
		log.Printf("marshaling key err: %v", err)
		return AnswerList{}, false
	}

	woListStr, ttl, err := client.GetRedis(context.Background(), string(keyStr))
	if err != nil {
		log.Printf("client.GetRedis err: %v", err)
		return AnswerList{}, false
	}

	var woList []wrappedObj
	err = json.Unmarshal([]byte(woListStr), &woList)
	if err != nil {
		log.Printf("unmarshaling key err: %v", err)
		return AnswerList{}, false
	}

	var answers []dns.RR
	for _, wo := range woList {
		if rrFunc, ok := dns.TypeToRR[wo.Type]; ok {
			payloadStr, err := json.Marshal(wo.Payload)
			if err != nil {
				log.Printf("marshaling payload err: %v", err)
			}
			rr := rrFunc()
			err = json.Unmarshal(payloadStr, &rr)
			if err != nil {
				log.Printf("unmarshaling payload err: %v", err)
			}
			answers = append(answers, rr)
		} else {
			log.Printf("unsupported rr type %d", wo.Type)
			return AnswerList{}, false
		}
	}

	return AnswerList{
		answers: answers,
		ttl:     uint32(ttl),
	}, true
}
