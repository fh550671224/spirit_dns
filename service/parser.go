package service

import (
	"encoding/binary"
	"fmt"
	"spiritDNS/shared"
	"strings"
)

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

func parseDNSQuestion(data []byte, Qdcount uint16) ([]*DNSQuestion, int) {
	offset := 12
	var questions []*DNSQuestion
	for Qdcount > 0 {
		q := DNSQuestion{}
		q.Name, offset = parseNameByOffset(data, offset)

		q.Type = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		q.Class = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		questions = append(questions, &q)
		Qdcount--
	}

	return questions, offset
}

func checkIsPointer(data []byte, offset int) (bool, uint16) {
	if data[offset]&0xC0 == 0xC0 {
		ptr := binary.BigEndian.Uint16(data[offset:offset+2]) - 0xC000
		return true, ptr
	} else {
		return false, 0
	}
}

func parseNameByOffset(data []byte, offset int) (string, int) {
	var names []string
	for {
		if data[offset] == 0 {
			offset++
			break
		}

		if isPointer, ptr := checkIsPointer(data, offset); isPointer {
			name, _ := parseNameByOffset(data, int(ptr))
			names = append(names, name)
			offset += 2
			if data[offset] == 0 {
				break
			}
		} else {
			length := int(data[offset])
			offset++
			name := data[offset : offset+length]
			names = append(names, string(name))
			offset += length
		}
	}

	return strings.Join(names, "."), offset
}

func parseDNSResourceRecord(data []byte, offset int, count uint16) ([]*DNSResourceRecord, int) {
	var records []*DNSResourceRecord
	for count > 0 {
		r := DNSResourceRecord{}

		r.Name, offset = parseNameByOffset(data, offset)

		r.Type = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		r.Class = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		r.TTL = binary.BigEndian.Uint32(data[offset : offset+4])
		offset += 4

		r.ResourceDataLength = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		switch r.Type {
		case shared.TYPE_A:
			r.ResourceData = parseARecordData(data, offset, r.ResourceDataLength)
		case shared.TYPE_NS, shared.TYPE_CNAME:
			r.ResourceData, _ = parseNameByOffset(data, offset)
		}
		offset += int(r.ResourceDataLength)

		records = append(records, &r)
		count--
	}

	return records, offset
}

func parseARecordData(data []byte, offset int, length uint16) string {
	var parts []string
	for length > 0 {
		parts = append(parts, fmt.Sprintf("%d", data[offset]))
		offset++
		length--
	}
	return strings.Join(parts, ".")
}

func parseDNSHeader(data []byte) *DNSHeader {
	return &DNSHeader{
		ID:            binary.BigEndian.Uint16(data[0:2]),
		Flags:         *parseDNSFlags(binary.BigEndian.Uint16(data[2:4])),
		QuestionCount: binary.BigEndian.Uint16(data[4:6]),
		AnswerCount:   binary.BigEndian.Uint16(data[6:8]),
		NSCount:       binary.BigEndian.Uint16(data[8:10]),
		ARCount:       binary.BigEndian.Uint16(data[10:12]),
	}
}

func parseDNSFlags(flags uint16) *DNSFlags {
	return &DNSFlags{
		QR:       (flags & 0x8000) != 0,
		OpCode:   uint8((flags & 0x7800) >> 11),
		AA:       (flags & 0x0400) != 0,
		TC:       (flags & 0x0200) != 0,
		RD:       (flags & 0x0100) != 0,
		RA:       (flags & 0x0080) != 0,
		Z:        uint8((flags & 0x0070) >> 4),
		RespCode: uint8(flags & 0x000f),
	}
}
