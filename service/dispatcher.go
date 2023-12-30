package service

import (
	"github.com/miekg/dns"
	"net"
	"sync"
)

type Packet struct {
	DnsMsg dns.Msg
	Ip     net.IP
}

var PackDispatcher *PacketDispatcher

type PacketDispatcher struct {
	listeners map[string]chan Packet
	mu        sync.RWMutex
}

func InitPacketDispatcher() {
	PackDispatcher = &PacketDispatcher{listeners: make(map[string]chan Packet)}
}

func (pd *PacketDispatcher) Register(key string, ch chan Packet) {
	pd.mu.Lock()
	defer pd.mu.Unlock()
	pd.listeners[key] = ch
}

func (pd *PacketDispatcher) UnRegister(key string) {
	pd.mu.Lock()
	defer pd.mu.Unlock()
	delete(pd.listeners, key)
}

func (pd *PacketDispatcher) Dispatch(pack Packet) {
	q := pack.DnsMsg.Question[0]
	pd.mu.RLock()
	defer pd.mu.RUnlock()
	if ch, ok := pd.listeners[q.Name]; ok {
		ch <- pack
	}
}
