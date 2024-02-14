package service

import "github.com/fh550671224/spirit_dns_public"

type Packet struct {
	DnsMsg *dns.Msg
	Ip     string
	Port   int
}
