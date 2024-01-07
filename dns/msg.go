package dns

type Msg struct {
	MsgHdr
	Question []Question
	Answer   []RR
	Ns       []RR
	Extra    []RR
}

type Question struct {
	Name   string
	QType  uint16
	QClass uint16
}

type MsgHdr struct {
	Id                 uint16
	Response           bool
	Opcode             int
	Authoritative      bool
	Truncated          bool
	RecursionDesired   bool
	RecursionAvailable bool
	Zero               bool
	Rcode              int
}

type Header struct {
	Id                                 uint16
	Bits                               uint16
	Qdcount, Ancount, Nscount, Arcount uint16
}

func (q *Question) pack(buf []byte, off int, compression map[string]uint16) (int, error) {
	off, err := packDomainName(q.Name, buf, off, compression)
	if err != nil {
		return off, err
	}
	off, err = packUint16(q.QType, buf, off)
	if err != nil {
		return off, err
	}
	off, err = packUint16(q.QClass, buf, off)
	if err != nil {
		return off, err
	}
	return off, nil
}

func (h *Header) pack(buf []byte, off int) (int, error) {
	off, err := packUint16(h.Id, buf, off)
	if err != nil {
		return off, err
	}
	off, err = packUint16(h.Bits, buf, off)
	if err != nil {
		return off, err
	}
	off, err = packUint16(h.Qdcount, buf, off)
	if err != nil {
		return off, err
	}
	off, err = packUint16(h.Ancount, buf, off)
	if err != nil {
		return off, err
	}
	off, err = packUint16(h.Nscount, buf, off)
	if err != nil {
		return off, err
	}
	off, err = packUint16(h.Arcount, buf, off)
	if err != nil {
		return off, err
	}

	return off, nil
}

func (msg *Msg) Pack() (buf []byte, err error) {
	var dh Header
	dh.Id = msg.Id
	dh.Bits = uint16(msg.Opcode)<<11 | uint16(msg.Rcode&0xF)
	if msg.Response {
		dh.Bits |= BIT_QR
	}
	if msg.Authoritative {
		dh.Bits |= BIT_AA
	}
	if msg.Truncated {
		dh.Bits |= BIT_TC
	}
	if msg.RecursionDesired {
		dh.Bits |= BIT_RD
	}
	if msg.RecursionAvailable {
		dh.Bits |= BIT_RA
	}
	dh.Qdcount = uint16(len(msg.Question))
	dh.Ancount = uint16(len(msg.Answer))
	dh.Nscount = uint16(len(msg.Ns))
	dh.Arcount = uint16(len(msg.Extra))

	// TODO calculate length without compression
	buf = make([]byte, 512)

	off := 0

	off, err = dh.pack(buf, off)
	if err != nil {
		return nil, err
	}

	compression := make(map[string]uint16)
	for _, q := range msg.Question {
		off, err = q.pack(buf, off, compression)
		if err != nil {
			return nil, err
		}
	}

	off, err = packRRSlice(msg.Answer, buf, off, compression)
	if err != nil {
		return nil, err
	}

	off, err = packRRSlice(msg.Ns, buf, off, compression)
	if err != nil {
		return nil, err
	}

	off, err = packRRSlice(msg.Extra, buf, off, compression)
	if err != nil {
		return nil, err
	}

	return buf[:off], nil
}

func (msg *Msg) Unpack(data []byte) (err error) {
	off := 0
	var dh Header
	dh.Id, off, err = unpackUint16(data, off)
	if err != nil {
		return err
	}
	dh.Bits, off, err = unpackUint16(data, off)
	if err != nil {
		return err
	}
	dh.Qdcount, off, err = unpackUint16(data, off)
	if err != nil {
		return err
	}
	dh.Ancount, off, err = unpackUint16(data, off)
	if err != nil {
		return err
	}
	dh.Nscount, off, err = unpackUint16(data, off)
	if err != nil {
		return err
	}
	dh.Arcount, off, err = unpackUint16(data, off)
	if err != nil {
		return err
	}
	msg.Id = dh.Id
	msg.Response = dh.Bits&BIT_QR != 0
	msg.Opcode = int(dh.Bits>>11) & 0xF
	msg.Authoritative = dh.Bits&BIT_AA != 0
	msg.Truncated = dh.Bits&BIT_TC != 0
	msg.RecursionDesired = dh.Bits&BIT_RD != 0
	msg.RecursionAvailable = dh.Bits&BIT_RA != 0
	msg.Rcode = int(dh.Bits & 0xF)

	// question
	for i := 0; i < int(dh.Qdcount); i++ {
		var q Question
		q.Name, off, err = unpackDomainName(data, off)
		if err != nil {
			return err
		}
		q.QType, off, err = unpackUint16(data, off)
		if err != nil {
			return err
		}
		q.QClass, off, err = unpackUint16(data, off)
		if err != nil {
			return err
		}
		msg.Question = append(msg.Question, q)
	}

	// rr
	msg.Answer, off, err = unpackRRSlice(data, off, int(dh.Ancount))
	if err != nil {
		return err
	}
	msg.Ns, off, err = unpackRRSlice(data, off, int(dh.Nscount))
	if err != nil {
		return err
	}
	msg.Extra, off, err = unpackRRSlice(data, off, int(dh.Arcount))
	if err != nil {
		return err
	}

	return nil
}

func (msg *Msg) SetQuestion(name string, qtype uint16) {
	msg.Question = append(msg.Question, Question{
		Name:   name,
		QType:  qtype,
		QClass: ClassINET,
	})
}

func (msg *Msg) IsInvalid() bool {
	if msg.MsgHdr.Response {
		if len(msg.Answer) == 0 && len(msg.Ns) == 0 {
			return true
		}
	}
	return false
}
