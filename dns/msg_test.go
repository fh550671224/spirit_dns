package dns

import (
	"github.com/miekg/dns"
	"log"
	"net"
	"testing"
)

func TestPack(t *testing.T) {
	msg := Msg{
		MsgHdr: MsgHdr{
			Id:                 51158,
			Response:           true,
			Opcode:             0,
			Authoritative:      false,
			Truncated:          false,
			RecursionDesired:   true,
			RecursionAvailable: true,
			Zero:               false,
			Rcode:              0,
		},
		Question: []Question{
			{
				Name:   "www.baidu.com.",
				QType:  1,
				QClass: 1,
			},
		},
		Answer: []RR{
			&CNAME{
				Hdr: RR_Header{
					Name:     "www.baidu.com.",
					Rrtype:   5,
					Class:    1,
					Ttl:      1200,
					Rdlength: 18,
				},
				Target: "www.a.shifen.com.",
			},
			&A{
				Hdr: RR_Header{
					Name:     "www.a.shifen.com.",
					Rrtype:   1,
					Class:    1,
					Ttl:      120,
					Rdlength: 4,
				},
				A: net.IP{
					36, 155, 132, 3,
				},
			},
		},
	}

	_msg := dns.Msg{
		MsgHdr: dns.MsgHdr{
			Id:                 51158,
			Response:           true,
			Opcode:             0,
			Authoritative:      false,
			Truncated:          false,
			RecursionDesired:   true,
			RecursionAvailable: true,
			Zero:               false,
			Rcode:              0,
		},
		Question: []dns.Question{
			{
				Name:   "www.baidu.com.",
				Qtype:  1,
				Qclass: 1,
			},
		},
		Answer: []dns.RR{
			&dns.CNAME{
				Hdr: dns.RR_Header{
					Name:     "www.baidu.com.",
					Rrtype:   5,
					Class:    1,
					Ttl:      1200,
					Rdlength: 18,
				},
				Target: "www.a.shifen.com.",
			},
			&dns.A{
				Hdr: dns.RR_Header{
					Name:     "www.a.shifen.com.",
					Rrtype:   1,
					Class:    1,
					Ttl:      120,
					Rdlength: 4,
				},
				A: net.IP{
					36, 155, 132, 3,
				},
			},
		},
	}
	_ = msg
	_ = _msg

	data, err := msg.Pack()
	if err != nil {
		t.Error(err)
	}
	t.Log(data)

	checkMsg := new(dns.Msg)
	err = checkMsg.Unpack(data)
	if err != nil {
		log.Printf("dns.UnPack err: %v\n", err)
	}

	_data, err := _msg.Pack()
	if err != nil {
		t.Error(err)
	}
	t.Log(_data)

}
