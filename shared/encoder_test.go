package shared

//
//import (
//	"encoding/json"
//	"fmt"
//	"strings"
//	"testing"
//)
//
//func TestGetContinuousSubSets(t *testing.T) {
//	subSets := getContinuousSubSets("www.baidu.com")
//	fmt.Println(subSets)
//}
//
//func TestNamePtrTree_Match(t *testing.T) {
//	tree := &NamePtrTree{root: &PtrNode{
//		key: "",
//		ptr: 0,
//		children: []*PtrNode{
//			{
//				key: "com",
//				ptr: 5,
//				children: []*PtrNode{
//					{
//						key: "baidu",
//						ptr: 1,
//						children: []*PtrNode{
//							{
//								key:      "www",
//								ptr:      123,
//								children: nil,
//							},
//						},
//					},
//				}},
//		},
//	}}
//
//	ptr := tree.Match(tree.root, []string{"baidu", "com"}, 1)
//	fmt.Println(ptr)
//}
//
////func TestNamePtrTree_FindPtrInTree(t *testing.T) {
////	tree := &NamePtrTree{root: &PtrNode{
////		key: "",
////		ptr: 0,
////		children: []*PtrNode{
////			{
////				key: "com",
////				ptr: 5,
////				children: []*PtrNode{
////					{
////						key: "baidu",
////						ptr: 1,
////						children: []*PtrNode{
////							{
////								key:      "www",
////								ptr:      123,
////								children: nil,
////							},
////						},
////					},
////				}},
////		},
////	}}
////
////	ptr, res := tree.FindPtrInTree("baidu.com")
////	fmt.Println(ptr, res)
////}
//
//func TestReservedSplit(t *testing.T) {
//	res := SplitWithSep("www.baidu.com", "baidu")
//	fmt.Println(res)
//}
////a
//func TestCalculatePtr(t *testing.T) {
//	ptr := calculatePtr([]string{"www", "baidu", "com"}, 2, 0)
//	fmt.Println(ptr)
//	// 3 w w w 5 b a i d u 3 c o m
//	// 0 1 2 3 4 5 6 7 8 9 10
//}
//
//func TestNamePtrTree_Insert(t *testing.T) {
//	tree := &NamePtrTree{root: &PtrNode{
//		key: "",
//		ptr: 0,
//		children: []*PtrNode{
//			{
//				key: "com",
//				ptr: 5,
//				children: []*PtrNode{
//					{
//						key: "baidu",
//						ptr: 1,
//						children: []*PtrNode{
//							{
//								key:      "www",
//								ptr:      123,
//								children: nil,
//							},
//						},
//					},
//				}},
//		},
//	}}
//
//	tree.Insert(tree.root, strings.Split("mail.baidu.com", "."), 2, 1234)
//	fmt.Println(tree)
//}
//
//func TestIsOverlapping(t *testing.T) {
//	is := IsOverlapping("ooooo", "o", "o")
//	fmt.Println(is)
//}
//
//func TestNamePtrTree_GetLongestBackwardSubsequence(t *testing.T) {
//	tree := &NamePtrTree{root: &PtrNode{
//		key: "",
//		ptr: 0,
//		children: []*PtrNode{
//			{
//				key: "net",
//				ptr: 5,
//				children: []*PtrNode{
//					{
//						key: "gtld-servers",
//						ptr: 1,
//						children: []*PtrNode{
//							{
//								key:      "e",
//								ptr:      123,
//								children: nil,
//							},
//						},
//					},
//				}},
//			//{
//			//	key:      "net",
//			//	ptr:      2234,
//			//	children: nil,
//			//},
//		},
//	}}
//
//	res := tree.GetPtrStartIndex(strings.Split("b.gtld-servers.net", "."))
//	fmt.Println(res)
//}
//
//func TestEncode(t *testing.T) {
//	msgStr := `{"Header":{"ID":2,"Flags":{"QR":true,"OpCode":0,"AA":false,"TC":true,"RD":true,"RA":false,"Z":0,"RespCode":0},"QuestionCount":1,"AnswerCount":0,"NSCount":13,"ARCount":11},"Questions":[{"Name":"www.baidu.com","QType":1,"QClass":1}],"AnswerRecords":null,"NSRecords":[{"Name":"com","QType":2,"QClass":1,"TTL":172800,"ResourceDataLength":20,"ResourceData":"e.gtld-servers.net."},{"Name":"com","QType":2,"QClass":1,"TTL":172800,"ResourceDataLength":4,"ResourceData":"b.gtld-servers.net"},{"Name":"com","QType":2,"QClass":1,"TTL":172800,"ResourceDataLength":4,"ResourceData":"j.gtld-servers.net"},{"Name":"com","QType":2,"QClass":1,"TTL":172800,"ResourceDataLength":4,"ResourceData":"m.gtld-servers.net"},{"Name":"com","QType":2,"QClass":1,"TTL":172800,"ResourceDataLength":4,"ResourceData":"i.gtld-servers.net"},{"Name":"com","QType":2,"QClass":1,"TTL":172800,"ResourceDataLength":4,"ResourceData":"f.gtld-servers.net"},{"Name":"com","QType":2,"QClass":1,"TTL":172800,"ResourceDataLength":4,"ResourceData":"a.gtld-servers.net"},{"Name":"com","QType":2,"QClass":1,"TTL":172800,"ResourceDataLength":4,"ResourceData":"g.gtld-servers.net"},{"Name":"com","QType":2,"QClass":1,"TTL":172800,"ResourceDataLength":4,"ResourceData":"h.gtld-servers.net"},{"Name":"com","QType":2,"QClass":1,"TTL":172800,"ResourceDataLength":4,"ResourceData":"l.gtld-servers.net"},{"Name":"com","QType":2,"QClass":1,"TTL":172800,"ResourceDataLength":4,"ResourceData":"k.gtld-servers.net"},{"Name":"com","QType":2,"QClass":1,"TTL":172800,"ResourceDataLength":4,"ResourceData":"c.gtld-servers.net"},{"Name":"com","QType":2,"QClass":1,"TTL":172800,"ResourceDataLength":4,"ResourceData":"d.gtld-servers.net"}],"AdditionalRecords":[{"Name":"e.gtld-servers.net","QType":1,"QClass":1,"TTL":172800,"ResourceDataLength":4,"ResourceData":"192.12.94.30"},{"Name":"e.gtld-servers.net","QType":28,"QClass":1,"TTL":172800,"ResourceDataLength":16,"ResourceData":""},{"Name":"b.gtld-servers.net.com","QType":1,"QClass":1,"TTL":172800,"ResourceDataLength":4,"ResourceData":"192.33.14.30"},{"Name":"b.gtld-servers.net.com","QType":28,"QClass":1,"TTL":172800,"ResourceDataLength":16,"ResourceData":""},{"Name":"j.gtld-servers.net.com","QType":1,"QClass":1,"TTL":172800,"ResourceDataLength":4,"ResourceData":"192.48.79.30"},{"Name":"j.gtld-servers.net.com","QType":28,"QClass":1,"TTL":172800,"ResourceDataLength":16,"ResourceData":""},{"Name":"m.gtld-servers.net.com","QType":1,"QClass":1,"TTL":172800,"ResourceDataLength":4,"ResourceData":"192.55.83.30"},{"Name":"m.gtld-servers.net.com","QType":28,"QClass":1,"TTL":172800,"ResourceDataLength":16,"ResourceData":""},{"Name":"i.gtld-servers.net.com","QType":1,"QClass":1,"TTL":172800,"ResourceDataLength":4,"ResourceData":"192.43.172.30"},{"Name":"i.gtld-servers.net.com","QType":28,"QClass":1,"TTL":172800,"ResourceDataLength":16,"ResourceData":""},{"Name":"f.gtld-servers.net.com","QType":1,"QClass":1,"TTL":172800,"ResourceDataLength":4,"ResourceData":"192.35.51.30"}]}`
//	var msg DNSMessage
//	json.Unmarshal([]byte(msgStr), &msg)
//
//	data, _ := EncodeDNSMessage(&msg, true)
//
//	msgg := DecodeDNSMessage(data)
//
//	fmt.Println(msgg)
//}
