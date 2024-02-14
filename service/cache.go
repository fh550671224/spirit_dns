package service

import (
	"github.com/fh550671224/spirit_dns_public"
	"math"
	"sync"
	"time"
)

var answerCache *AnswerCache

type AnswerList struct {
	answers []dns.RR
	ttl     uint32
}

type AnswerCache struct {
	questionAnswerMap map[dns.Question]AnswerList
	mu                sync.RWMutex
}

func InitCache() {
	answerCache = &AnswerCache{
		questionAnswerMap: make(map[dns.Question]AnswerList),
	}

	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				answerCache.mu.Lock()
				for key, answerList := range answerCache.questionAnswerMap {
					answerList.ttl--
					if answerList.ttl == 0 {
						delete(answerCache.questionAnswerMap, key)
					}
					for _, answer := range answerList.answers {
						answer.Header().Ttl--
					}
				}
				answerCache.mu.Unlock()
			}
		}
	}()
}

func (ac *AnswerCache) Store(key dns.Question, answers []dns.RR) {
	var ttl uint32
	ttl = math.MaxUint32
	for _, a := range answers {
		if ttl > a.Header().Ttl {
			ttl = a.Header().Ttl
		}
	}

	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.questionAnswerMap[key] = AnswerList{
		answers: answers,
		ttl:     ttl,
	}
}

func (ac *AnswerCache) Get(key dns.Question) (AnswerList, bool) {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	a, ok := ac.questionAnswerMap[key]
	return a, ok
}
