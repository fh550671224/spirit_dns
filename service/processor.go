package service

import (
	"fmt"
	"net"
	"spiritDNS/shared"
)

type Processor interface {
	Listen()
	HandleConnetion(conn net.Conn)
	FindAnswers(ipList []string) ([]byte, error)
}

type Context struct {
	RawData           []byte
	Header            *DNSHeader
	Questions         []*DNSQuestion
	AnswerRecords     []*DNSResourceRecord
	NSRecords         []*DNSResourceRecord
	AdditionalRecords []*DNSResourceRecord
}

func NewContext(data []byte) *Context {
	return &Context{
		RawData: data,
	}
}

func StarListening() {
	udpp := UDPProcessor{}
	udpp.Listen()
}

func ProcessReq(p Processor) {
	ctx := Context{}
	// 处理数据
	ctx.Header = parseDNSHeader(ctx.RawData)

	ctx.Questions, _ = parseDNSQuestion(ctx.RawData, ctx.Header.QuestionCount)

	if !ctx.Header.Flags.QR {
		for _, q := range ctx.Questions {
			fmt.Println(q.Name)
			switch q.Type {
			case shared.TYPE_A:
				fmt.Println("A record")
				// establish conn to root server
				p.FindAnswers(shared.ROOT_DNS_SERVERS)
			}
		}
	}
}
