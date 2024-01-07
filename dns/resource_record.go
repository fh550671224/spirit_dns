package dns

type RR interface {
	Header() *RR_Header

	pack(buf []byte, off int, compression map[string]uint16) (off1 int, err error)

	unpack(msg []byte, off int) (off1 int, err error)
}

type RR_Header struct {
	Name     string
	Rrtype   uint16
	Class    uint16
	Ttl      uint32
	Rdlength uint16
}

func (h *RR_Header) Header() *RR_Header {
	return h
}

func (h *RR_Header) pack(msg []byte, off int, compression map[string]uint16) (off1 int, err error) {
	return off, nil
}

func (h *RR_Header) unpack(msg []byte, off int) (off1 int, err error) {
	panic("dns: internal error: unpack should never be called on RR_Header")
}
