package service

import "github.com/miekg/dns"

type Packet struct {
	DnsMsg dns.Msg
	Ip     string
	Port   int
}
