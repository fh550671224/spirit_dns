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
		return nil, fmt.Errorf("encodeDNSResourceRecords err: %s", err.Error())
	}
	udpSafeCount = getUdpSafeCount(indices)
	if isUdp && udpSafeCount < msg.Header.AnswerCount {
		return ctx.wrapUdpMsg(msg, lastOffset, indices, udpSafeCount, COUNT_TYPE_ANSWER)
	}

	// ns record
	indices, err = ctx.encodeDNSResourceRecords(msg.NSRecords)
	if err != nil {
		return nil, fmt.Errorf("encodeDNSResourceRecords err: %s", err.Error())
	}
	udpSafeCount = getUdpSafeCount(indices)
	if isUdp && udpSafeCount < msg.Header.NSCount {
		return ctx.wrapUdpMsg(msg, lastOffset, indices, udpSafeCount, COUNT_TYPE_NS)
	}

	// additional record
	indices, err = ctx.encodeDNSResourceRecords(msg.AdditionalRecords)
	if err != nil {
		return nil, fmt.Errorf("encodeDNSResourceRecords err: %s", err.Error())
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
	for _, r := range records {
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
				return nil, fmt.Errorf("encodeARecordData err: %s", err.Error())
			}
		case TYPE_NS, TYPE_CNAME:
			ctx.encodeName(r.ResourceData)
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
			return fmt.Errorf(" strconv.ParseInt err: %s", err.Error())
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

func (tree *NamePtrTree) GetLongestBackwardSubsequence(names []string) [][]string {
	var res [][]string
	var end = len(names)

	flip := false
	flag := false

	// 对于每个sub name，获取从当前位置出发的最长命中列表
	for i := len(names) - 1; i >= 0; i-- {
		tmp := names[i:end]
		match := tree.Match(tree.root, tmp, len(tmp)-1)
		if match != nil {
			continue
		} else {
			if end-i == 1 {
				res = append(res, names[i:end])
				end = i
			} else {
				res = append(res, names[i+1:end])
				end = i + 1
			}
		}
	}
	if end > 0 {
		res = append(res, names[0:end])
	}

	return res
}

func (ctx *encodeContext) encodeName(name string) {
	offset := ctx.offset
	buffer := ctx.buffer

	// 对于每个sub name，获取从当前位置出发的最长命中列表
	lbs := ctx.NamePtrTree.GetLongestBackwardSubsequence(strings.Split(name, "."))
	for i := len(lbs) - 1; i >= 0; i-- {
		keys := lbs[i]

		node := ctx.NamePtrTree.Match(ctx.NamePtrTree.root, keys, len(keys)-1)
		if node != nil {
			// 将指针写入报文
			tmp := 0xC000 | node.ptr
			binary.BigEndian.PutUint16(buffer[offset:offset+2], tmp)
			offset += 2
			continue
		} else {
			// 非指针
			// 插入树
			ptr := uint16(offset)
			ctx.NamePtrTree.Insert(ctx.NamePtrTree.root, keys, len(keys)-1, ptr)

			// 将域名写入报文
			for _, key := range keys {
				buffer[offset] = uint8(len(key))
				offset++
				copy(buffer[offset:offset+len(key)], key)
				offset += len(key)
			}

			// 如果非指针，且已记录完毕，则在最后添加终止符
			if i == 0 {
				buffer[offset] = 0
				offset++
			}
		}
	}
	ctx.offset = offset
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
