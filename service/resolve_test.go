package service

import (
	"fmt"
	dns "github.com/fh550671224/spirit_dns_public"
	"net"
	"spiritDNS/shared"
	"testing"
)

func TestResolve(t *testing.T) {
	InitCache()

	m := new(dns.Msg)
	m.SetQuestion("www.tencent.com.", dns.TypeA)

	answer, err := Resolve(m, shared.ROOT_DNS_SERVERS)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(answer)
}

func TestTrySendUDP(t *testing.T) {
	m := new(dns.Msg)
	m.SetQuestion("baidu.com", dns.TypeA)

	ms, _ := m.Pack()

	TrySendUDP([]*net.UDPAddr{
		{
			IP:   net.ParseIP("198.41.0.4"),
			Port: 53,
			Zone: "",
		},
	}, ms)
}
