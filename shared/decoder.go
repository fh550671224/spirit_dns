package shared

import (
	"encoding/binary"
	"fmt"
	"strings"
)

func decodeDNSQuestion(data []byte, Qdcount uint16) ([]*DNSQuestion, int) {
	offset := 12
	var questions []*DNSQuestion
	for Qdcount > 0 {
		q := DNSQuestion{}
		q.Name, offset = decodeNameByOffset(data, offset)

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

func decodeNameWithLength(data []byte, offset int, dataLength int) string {
	var names []string
	originalOffset := offset
	for {
		if offset-originalOffset == dataLength {
			break
		}

		if isPointer, ptr := checkIsPointer(data, offset); isPointer {
			name, _ := decodeNameByOffset(data, int(ptr))
			names = append(names, name)
			offset += 2
		} else {
			length := int(data[offset])
			offset++
			name := data[offset : offset+length]
			names = append(names, string(name))
			offset += length
			if data[offset] == 0 {
				offset++
				break
			}
		}
	}

	return strings.Join(names, ".")
}

func decodeNameByOffset(data []byte, offset int) (string, int) {
	var names []string
	for {
		if isPointer, ptr := checkIsPointer(data, offset); isPointer {
			name, _ := decodeNameByOffset(data, int(ptr))
			names = append(names, name)
			offset += 2
			break
		} else {
			length := int(data[offset])
			offset++
			name := data[offset : offset+length]
			names = append(names, string(name))
			offset += length
			if data[offset] == 0 {
				offset++
				break
			}
		}
	}

	return strings.Join(names, "."), offset
}

func decodeDNSResourceRecord(data []byte, offset int, count uint16) ([]*DNSResourceRecord, int) {
	var records []*DNSResourceRecord
	for count > 0 {
		r := DNSResourceRecord{}

		r.Name, offset = decodeNameByOffset(data, offset)

		r.Type = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		r.Class = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		r.TTL = binary.BigEndian.Uint32(data[offset : offset+4])
		offset += 4

		r.ResourceDataLength = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		switch r.Type {
		case TYPE_A:
			r.ResourceData = decodeARecordData(data, offset, r.ResourceDataLength)
		case TYPE_NS, TYPE_CNAME:
			r.ResourceData = decodeNameWithLength(data, offset, int(r.ResourceDataLength))
		}
		offset += int(r.ResourceDataLength)

		records = append(records, &r)
		count--
	}

	return records, offset
}

func decodeARecordData(data []byte, offset int, length uint16) string {
	var parts []string
	for length > 0 {
		parts = append(parts, fmt.Sprintf("%d", data[offset]))
		offset++
		length--
	}
	return strings.Join(parts, ".")
}

func decodeDNSHeader(data []byte) *DNSHeader {
	return &DNSHeader{
		ID:            binary.BigEndian.Uint16(data[0:2]),
		Flags:         *decodeDNSFlags(binary.BigEndian.Uint16(data[2:4])),
		QuestionCount: binary.BigEndian.Uint16(data[4:6]),
		AnswerCount:   binary.BigEndian.Uint16(data[6:8]),
		NSCount:       binary.BigEndian.Uint16(data[8:10]),
		ARCount:       binary.BigEndian.Uint16(data[10:12]),
	}
}

func decodeDNSFlags(flags uint16) *DNSFlags {
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

func DecodeDNSMessage(data []byte) DNSMessage {
	msg := DNSMessage{}
	msg.Header = *decodeDNSHeader(data)
	var offset int
	msg.Questions, offset = decodeDNSQuestion(data, msg.Header.QuestionCount)
	msg.AnswerRecords, offset = decodeDNSResourceRecord(data, offset, msg.Header.AnswerCount)
	msg.NSRecords, offset = decodeDNSResourceRecord(data, offset, msg.Header.NSCount)
	msg.AdditionalRecords, offset = decodeDNSResourceRecord(data, offset, msg.Header.ARCount)

	return msg
}
