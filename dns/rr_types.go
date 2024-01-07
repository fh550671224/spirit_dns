package dns

import (
	"net"
)

type A struct {
	Hdr RR_Header
	A   net.IP
}

type AAAA struct {
	Hdr  RR_Header
	AAAA net.IP
}

type CNAME struct {
	Hdr    RR_Header
	Target string
}

type NS struct {
	Hdr RR_Header
	Ns  string
}

func (rr *NS) Header() *RR_Header {
	return &rr.Hdr
}

func (rr *NS) pack(msg []byte, off int, compression map[string]uint16) (off1 int, err error) {
	off, err = packDomainName(rr.Ns, msg, off, compression)
	if err != nil {
		return 0, err
	}

	return off, nil
}

func (rr *NS) unpack(msg []byte, off int) (off1 int, err error) {
	rr.Ns, off, err = unpackDomainName(msg, off)
	if err != nil {
		return 0, err
	}

	return off, nil
}
