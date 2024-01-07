package service

import "spiritDNS/dns"

type Packet struct {
	DnsMsg *dns.Msg
	Ip     string
	Port   int
}
