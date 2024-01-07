package dns

func (rr *A) Header() *RR_Header {
	return &rr.Hdr
}

func (rr *A) pack(msg []byte, off int, compression map[string]uint16) (off1 int, err error) {
	off, err = packDataA(rr.A, msg, off)
	if err != nil {
		return off, err
	}
	return off, nil
}

func (rr *A) unpack(msg []byte, off int) (off1 int, err error) {
	rr.A, off, err = unpackDataA(msg, off)
	if err != nil {
		return off, err
	}
	return off, nil
}

func (rr *AAAA) Header() *RR_Header {
	return &rr.Hdr
}

func (rr *AAAA) pack(msg []byte, off int, compression map[string]uint16) (off1 int, err error) {
	off, err = packDataAAAA(rr.AAAA, msg, off)
	if err != nil {
		return off, err
	}
	return off, nil
}

func (rr *AAAA) unpack(msg []byte, off int) (off1 int, err error) {
	rr.AAAA, off, err = unpackDataAAAA(msg, off)
	if err != nil {
		return off, err
	}
	return off, nil
}

func (rr *CNAME) Header() *RR_Header {
	return &rr.Hdr
}

func (rr *CNAME) pack(msg []byte, off int, compression map[string]uint16) (off1 int, err error) {
	off, err = packDomainName(rr.Target, msg, off, compression)
	if err != nil {
		return off, err
	}

	return off, nil
}

func (rr *CNAME) unpack(msg []byte, off int) (off1 int, err error) {
	rr.Target, off, err = unpackDomainName(msg, off)
	if err != nil {
		return off, err
	}
	return off, nil
}

func packRRSlice(rrs []RR, buf []byte, off int, compression map[string]uint16) (off1 int, err error) {
	for _, rr := range rrs {
		off, err = packDomainName(rr.Header().Name, buf, off, compression)
		if err != nil {
			return off, err
		}

		off, err = packUint16(rr.Header().Rrtype, buf, off)
		if err != nil {
			return off, err
		}

		off, err = packUint16(rr.Header().Class, buf, off)
		if err != nil {
			return off, err
		}

		off, err = packUint32(rr.Header().Ttl, buf, off)
		if err != nil {
			return off, err
		}

		// set rdLength later
		offAtRdLengh := off
		off += 2

		off, err = rr.pack(buf, off, compression)
		if err != nil {
			return off, err
		}

		rdLength := off - offAtRdLengh - 2
		_, err = packUint16(uint16(rdLength), buf, offAtRdLengh)
		if err != nil {
			return off, err
		}
	}

	return off, nil
}
