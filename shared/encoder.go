package shared

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

type PtrNode struct {
	key      string
	ptr      uint16
	children []*PtrNode
}

type NamePtrTree struct {
	root *PtrNode
}

type encodeContext struct {
	buffer      []byte
	offset      int
	NamePtrTree *NamePtrTree
}

func EncodeDNSMessage(msg *DNSMessage, isUdp bool) ([]byte, error) {
	ctx := encodeContext{
		buffer:      make([]byte, 1024),
		offset:      0,
		NamePtrTree: &NamePtrTree{root: &PtrNode{}},
	}

	ctx.encodeDNSHeader(&msg.Header, ctx.buffer)
	lastOffset := ctx.offset

	indices := ctx.encodeDNSQuestions(msg.Questions)
	udpSafeCount := getUdpSafeCount(indices)
	if isUdp && udpSafeCount < msg.Header.QuestionCount {
		return ctx.wrapUdpMsg(msg, lastOffset, indices, udpSafeCount, COUNT_TYPE_QUESTION)

	}
	lastOffset = ctx.offset

	// answer record
	indices, err := ctx.encodeDNSResourceRecords(msg.AnswerRecords)
	if err != nil {
		return nil, fmt.Errorf("encodeDNSResourceRecords err: %v", err)
	}
	udpSafeCount = getUdpSafeCount(indices)
	if isUdp && udpSafeCount < msg.Header.AnswerCount {
		return ctx.wrapUdpMsg(msg, lastOffset, indices, udpSafeCount, COUNT_TYPE_ANSWER)
	}

	// ns record
	indices, err = ctx.encodeDNSResourceRecords(msg.NSRecords)
	if err != nil {
		return nil, fmt.Errorf("encodeDNSResourceRecords err: %v", err)
	}
	udpSafeCount = getUdpSafeCount(indices)
	if isUdp && udpSafeCount < msg.Header.NSCount {
		return ctx.wrapUdpMsg(msg, lastOffset, indices, udpSafeCount, COUNT_TYPE_NS)
	}

	// additional record
	indices, err = ctx.encodeDNSResourceRecords(msg.AdditionalRecords)
	if err != nil {
		return nil, fmt.Errorf("encodeDNSResourceRecords err: %v", err)
	}
	udpSafeCount = getUdpSafeCount(indices)
	if isUdp && udpSafeCount < msg.Header.ARCount {
		return ctx.wrapUdpMsg(msg, lastOffset, indices, udpSafeCount, COUNT_TYPE_AR)
	}

	return ctx.buffer[:ctx.offset], nil
}

func (ctx *encodeContext) wrapUdpMsg(msg *DNSMessage, lastOffset int, indices []uint16, udpSafeCount uint16, countType int) ([]byte, error) {
	msg.Header.Flags.TC = true

	switch countType {
	case COUNT_TYPE_QUESTION:
		msg.Header.QuestionCount = udpSafeCount
	case COUNT_TYPE_ANSWER:
		msg.Header.AnswerCount = udpSafeCount
	case COUNT_TYPE_NS:
		msg.Header.NSCount = udpSafeCount
	case COUNT_TYPE_AR:
		msg.Header.ARCount = udpSafeCount
	}
	ctx.encodeDNSHeader(&msg.Header, ctx.buffer)

	if udpSafeCount == 0 {
		// 如果没有任何结果可以保留，那就使用上次的offset
		return ctx.buffer[:lastOffset], nil
	} else {
		return ctx.buffer[:indices[udpSafeCount-1]], nil
	}
}

func (ctx *encodeContext) encodeDNSResourceRecords(records []*DNSResourceRecord) (indices []uint16, err error) {
	for i, r := range records {
		fmt.Println(i)
		ctx.encodeName(r.Name)

		binary.BigEndian.PutUint16(ctx.buffer[ctx.offset:ctx.offset+2], r.Type)
		ctx.offset += 2

		binary.BigEndian.PutUint16(ctx.buffer[ctx.offset:ctx.offset+2], r.Class)
		ctx.offset += 2

		binary.BigEndian.PutUint32(ctx.buffer[ctx.offset:ctx.offset+4], r.TTL)
		ctx.offset += 4

		binary.BigEndian.PutUint16(ctx.buffer[ctx.offset:ctx.offset+2], r.ResourceDataLength)
		ctx.offset += 2

		switch r.Type {
		case TYPE_A:
			err = ctx.encodeARecordData(r.ResourceData)
			if err != nil {
				return nil, fmt.Errorf("encodeARecordData err: %v", err)
			}
		case TYPE_NS, TYPE_CNAME:
			length := ctx.encodeName(r.ResourceData)
			binary.BigEndian.PutUint16(ctx.buffer[ctx.offset-length-2:ctx.offset-length], uint16(length))
		}

		indices = append(indices, uint16(ctx.offset))
	}

	return indices, nil
}

func (ctx *encodeContext) encodeARecordData(data string) error {
	parts := SplitWithoutEmpty(data, ".")
	for _, part := range parts {
		tmp, err := strconv.Atoi(part)
		if err != nil {
			return fmt.Errorf(" strconv.ParseInt err: %v", err)
		}
		ctx.buffer[ctx.offset] = uint8(tmp)
		ctx.offset++
	}

	return nil
}

func getUdpSafeCount(indices []uint16) uint16 {
	for i := 0; i < len(indices); i++ {
		if indices[i] >= 512 {
			return uint16(i)
		}
	}
	return uint16(len(indices))
}

func (ctx *encodeContext) encodeDNSQuestions(questions []*DNSQuestion) (indices []uint16) {
	for _, q := range questions {
		ctx.encodeName(q.Name)

		binary.BigEndian.PutUint16(ctx.buffer[ctx.offset:ctx.offset+2], q.Type)
		ctx.offset += 2

		binary.BigEndian.PutUint16(ctx.buffer[ctx.offset:ctx.offset+2], q.Class)
		ctx.offset += 2

		indices = append(indices, uint16(ctx.offset))
	}

	return indices
}

func (tree *NamePtrTree) GetPtrStartIndex(names []string) int {
	if len(names) == 0 {
		return 0
	}

	res := []string{names[len(names)-1]}
	if tree.Match(tree.root, res, len(res)-1) == nil {
		return len(names)
	}

	for i := len(names) - 2; i >= 0; i-- {
		temp := append([]string{names[i]}, res...)
		if tree.Match(tree.root, temp, len(temp)-1) != nil {
			res = temp
		} else {
			return i + 1
		}
	}

	return 0
}

func (ctx *encodeContext) encodeName(name string) int {
	offset := ctx.offset
	buffer := ctx.buffer

	names := SplitWithoutEmpty(name, ".")
	// 获取指针起始index，有三种情况：
	// 1. 纯指针
	// 2. 混合情况，前面是非指针，后面是指针
	// 3. 纯非指针
	ptrStartIndex := ctx.NamePtrTree.GetPtrStartIndex(names)
	if ptrStartIndex < len(names) {
		if ptrStartIndex == 0 {
			// 纯指针
			node := ctx.NamePtrTree.Match(ctx.NamePtrTree.root, names, len(names)-1)
			// 将指针写入报文
			tmp := 0xC000 | node.ptr
			binary.BigEndian.PutUint16(buffer[offset:offset+2], tmp)
			offset += 2
		} else {
			// 混合情况，前面是非指针，后面是指针

			// 1. 非指针部分
			nonPtrNames := names[:ptrStartIndex]
			// 插入树
			ptr := uint16(offset)
			ctx.NamePtrTree.Insert(ctx.NamePtrTree.root, nonPtrNames, len(nonPtrNames)-1, ptr)
			// 将部分域名写入报文
			for _, key := range nonPtrNames {
				buffer[offset] = uint8(len(key))
				offset++
				copy(buffer[offset:offset+len(key)], key)
				offset += len(key)
			}

			// 2. 指针部分
			ptrNames := names[ptrStartIndex:]
			node := ctx.NamePtrTree.Match(ctx.NamePtrTree.root, ptrNames, len(ptrNames)-1)
			// 将指针写入报文
			tmp := 0xC000 | node.ptr
			binary.BigEndian.PutUint16(buffer[offset:offset+2], tmp)
			offset += 2
		}

	} else {
		// 纯非指针
		// 插入树
		ptr := uint16(offset)
		ctx.NamePtrTree.Insert(ctx.NamePtrTree.root, names, len(names)-1, ptr)

		// 将域名写入报文
		for _, key := range names {
			buffer[offset] = uint8(len(key))
			offset++
			copy(buffer[offset:offset+len(key)], key)
			offset += len(key)
		}

		// 已记录完毕，在最后添加终止符
		buffer[offset] = 0
		offset++

	}

	length := offset - ctx.offset
	ctx.offset = offset
	return length
}

func (tree *NamePtrTree) Insert(node *PtrNode, keys []string, curIndex int, ptr uint16) {
	if curIndex == 0 {
		node.children = append(node.children, &PtrNode{
			key: keys[curIndex],
			ptr: calculatePtr(keys, curIndex, ptr),
		})
		return
	}

	find := false
	for _, child := range node.children {
		if child.key == keys[curIndex] {
			find = true
			tree.Insert(child, keys, curIndex-1, ptr)
		}
	}

	if !find {
		n := &PtrNode{
			key: keys[curIndex],
			ptr: calculatePtr(keys, curIndex, ptr),
		}
		node.children = append(node.children, n)

		tree.Insert(n, keys, curIndex-1, ptr)
	}
}

func calculatePtr(keys []string, curIndex int, ptr uint16) uint16 {
	offset := ptr
	for i := 0; i < len(keys); i++ {
		if i == curIndex {
			break
		}
		offset += 1 + uint16(len(keys[i]))

	}

	return offset
}

func (tree *NamePtrTree) FindPtrInTree(key string) map[string]*PtrNode {
	// name的所有连续子集都需要在树中遍历
	// 获取所有连续子集
	subSets := getContinuousSubSets(key)

	res := make(map[string]*PtrNode)
	for i := len(subSets) - 1; i >= 0; i-- {
		matchPtr := tree.Match(tree.root, subSets[i], len(subSets[i])-1)
		if matchPtr != nil {
			target := strings.Join(subSets[i], ".")
			isOverlapping := false
			for str := range res {
				if IsOverlapping(key, str, target) {
					isOverlapping = true
					break
				}
			}
			if !isOverlapping {
				res[target] = matchPtr
			}
		}
	}
	return res
}

func (tree *NamePtrTree) Match(node *PtrNode, keys []string, curIndex int) *PtrNode {
	if curIndex == -1 {
		return node
	}

	for _, child := range node.children {
		if child.key == keys[curIndex] {
			return tree.Match(child, keys, curIndex-1)
		}
	}

	return nil
}

func getContinuousSubSets(name string) [][]string {
	var subSets [][]string
	names := SplitWithoutEmpty(name, ".")
	for length := 1; length <= len(names); length++ {
		for i := 0; i+length <= len(names); i++ {
			subSets = append(subSets, names[i:i+length])
		}
	}
	return subSets
}

func (ctx *encodeContext) encodeDNSHeader(header *DNSHeader, buffer []byte) {
	binary.BigEndian.PutUint16(buffer[0:2], header.ID)
	binary.BigEndian.PutUint16(buffer[2:4], encodeDNSFlags(&header.Flags))
	binary.BigEndian.PutUint16(buffer[4:6], header.QuestionCount)
	binary.BigEndian.PutUint16(buffer[6:8], header.AnswerCount)
	binary.BigEndian.PutUint16(buffer[8:10], header.NSCount)
	binary.BigEndian.PutUint16(buffer[10:12], header.ARCount)

	ctx.offset = 12
}

func encodeDNSFlags(flags *DNSFlags) (res uint16) {
	var qr uint16
	var opCode uint16
	var aa uint16
	var tc uint16
	var rd uint16
	var ra uint16
	var z uint16
	var respCode uint16

	if flags.QR {
		qr = 0x8000
	}
	opCode = uint16(flags.OpCode) << 11
	if flags.AA {
		aa = 0x0400
	}
	if flags.TC {
		tc = 0x0200
	}
	if flags.RD {
		rd = 0x0100
	}
	if flags.RA {
		ra = 0x0080
	}
	z = uint16(flags.Z) << 4
	respCode = uint16(flags.RespCode)

	return qr | opCode | aa | tc | rd | ra | z | respCode
}
