package service

import (
	"github.com/miekg/dns"
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

func (ac *AnswerCache) Store(key dns.Question, value AnswerList) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.questionAnswerMap[key] = value
}

func (ac *AnswerCache) Get(key dns.Question) (AnswerList, bool) {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	//if a, ok := ac.questionAnswerMap[key]; ok {
	//	return &a
	//} else {
	//	return nil
	//}

	a, ok := ac.questionAnswerMap[key]
	return a, ok
}
