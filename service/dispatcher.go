package service

import (
	"net"
	"spiritDNS/shared"
	"sync"
)

type Packet struct {
	DnsMsg shared.DNSMessage
	Ip     net.IP
}

var PackDispatcher *PacketDispatcher

type PacketDispatcher struct {
	listeners map[shared.DNSQuestion]chan Packet
	mu        sync.RWMutex
}

func (d *PacketDispatcher) Init() {
	PackDispatcher = &PacketDispatcher{listeners: make(map[shared.DNSQuestion]chan Packet)}
}

func (d *PacketDispatcher) Register(key shared.DNSQuestion, ch chan Packet) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.listeners[key] = ch
}

func (d *PacketDispatcher) UnRegister(key shared.DNSQuestion) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.listeners, key)
}

func (d *PacketDispatcher) Dispatch(pack Packet) {
	// TODO 支持multi-question
	q := pack.DnsMsg.Questions[0]
	d.mu.RLock()
	defer d.mu.RUnlock()
	if ch, ok := d.listeners[*q]; ok {
		ch <- pack
	}
}
