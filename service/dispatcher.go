package service

import (
	"net"
	"sync"
)

var dispatcher *Dispatcher

type Dispatcher struct {
	listeners map[*net.UDPAddr]chan []byte
	mu        sync.RWMutex
}

func InitPacketDispatcher() {
	dispatcher = &Dispatcher{listeners: make(map[*net.UDPAddr]chan []byte)}
}

func (pd *Dispatcher) Register(key *net.UDPAddr, ch chan []byte) {
	pd.mu.Lock()
	defer pd.mu.Unlock()
	pd.listeners[key] = ch
}

func (pd *Dispatcher) UnRegister(key *net.UDPAddr) {
	pd.mu.Lock()
	defer pd.mu.Unlock()
	delete(pd.listeners, key)
}

func (pd *Dispatcher) Dispatch(key *net.UDPAddr, data []byte) {
	pd.mu.RLock()
	defer pd.mu.RUnlock()
	if ch, ok := pd.listeners[key]; ok {
		ch <- data
	}
}
