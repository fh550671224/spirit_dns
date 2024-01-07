package service

import (
	"fmt"
	"spiritDNS/dns"
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
