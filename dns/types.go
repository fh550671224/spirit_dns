package dns

const (
	headerSize = 12

	// Header.Bits
	BIT_QR = 1 << 15 // query/response (response=1)
	BIT_AA = 1 << 10 // authoritative
	BIT_TC = 1 << 9  // truncated
	BIT_RD = 1 << 8  // recursion desired
	BIT_RA = 1 << 7  // recursion available
)

const (
	// valid RR_Header.Rrtype and Question.qtype

	TypeNone  uint16 = 0
	TypeA     uint16 = 1
	TypeNS    uint16 = 2
	TypeMD    uint16 = 3
	TypeMF    uint16 = 4
	TypeCNAME uint16 = 5
	TypeAAAA  uint16 = 28
	TypeOPT   uint16 = 41
)

const (
	ClassINET = 1
)

const (
	RcodeSuccess     = 0 // NoError   - No Error                          [DNS]
	RcodeFormatError = 1 // FormErr   - Format Error                      [DNS]
)

var TypeToRR = map[uint16]func() RR{
	TypeA:     func() RR { return new(A) },
	TypeAAAA:  func() RR { return new(AAAA) },
	TypeNS:    func() RR { return new(NS) },
	TypeCNAME: func() RR { return new(CNAME) },
}
