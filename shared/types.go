package shared

import "net"

type DNSPacket struct {
	Msg  DNSMessage
	Addr net.Addr
}

type DNSMessage struct {
	//RawData           []byte
	Header            DNSHeader
	Questions         []*DNSQuestion
	AnswerRecords     []*DNSResourceRecord
	NSRecords         []*DNSResourceRecord
	AdditionalRecords []*DNSResourceRecord
}

type DNSQuestion struct {
	Name  string
	Type  uint16
	Class uint16
}

type DNSHeader struct {
	ID            uint16
	Flags         DNSFlags
	QuestionCount uint16
	AnswerCount   uint16
	NSCount       uint16
	ARCount       uint16
}

type DNSFlags struct {
	QR       bool  // 是否是response
	OpCode   uint8 // 操作码
	AA       bool  // 权威答案
	TC       bool  // 是否截断
	RD       bool  // 是否期望递归
	RA       bool  // 是否支持递归
	Z        uint8 // 保留，通常为0
	RespCode uint8 // 响应码
}

type DNSResourceRecord struct {
	Name               string
	Type               uint16
	Class              uint16
	TTL                uint32
	ResourceDataLength uint16
	ResourceData       string
}
