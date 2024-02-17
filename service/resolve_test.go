package service

import (
	"fmt"
	dns "github.com/fh550671224/spirit_dns_public"
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
