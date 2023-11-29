package shared

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetContinuousSubSets(t *testing.T) {
	subSets := getContinuousSubSets("www.baidu.com")
	fmt.Println(subSets)
}

func TestNamePtrTree_Match(t *testing.T) {
	tree := &NamePtrTree{root: &PtrNode{
		key: "",
		ptr: 0,
		children: []*PtrNode{
			{
				key: "com",
				ptr: 5,
				children: []*PtrNode{
					{
						key: "baidu",
						ptr: 1,
						children: []*PtrNode{
							{
								key:      "www",
								ptr:      123,
								children: nil,
							},
						},
					},
				}},
		},
	}}

	ptr := tree.Match(tree.root, []string{"baidu", "com"}, 1)
	fmt.Println(ptr)
}

//func TestNamePtrTree_FindPtrInTree(t *testing.T) {
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
//	ptr, res := tree.FindPtrInTree("baidu.com")
//	fmt.Println(ptr, res)
//}

func TestReservedSplit(t *testing.T) {
	res := SplitWithSep("www.baidu.com", "baidu")
	fmt.Println(res)
}

func TestCalculatePtr(t *testing.T) {
	ptr := calculatePtr([]string{"www", "baidu", "com"}, 2, 0)
	fmt.Println(ptr)
	// 3 w w w 5 b a i d u 3 c o m
	// 0 1 2 3 4 5 6 7 8 9 10
}

func TestNamePtrTree_Insert(t *testing.T) {
	tree := &NamePtrTree{root: &PtrNode{
		key: "",
		ptr: 0,
		children: []*PtrNode{
			{
				key: "com",
				ptr: 5,
				children: []*PtrNode{
					{
						key: "baidu",
						ptr: 1,
						children: []*PtrNode{
							{
								key:      "www",
								ptr:      123,
								children: nil,
							},
						},
					},
				}},
		},
	}}

	tree.Insert(tree.root, strings.Split("mail.baidu.com", "."), 2, 1234)
	fmt.Println(tree)
}

func TestIsOverlapping(t *testing.T) {
	is := IsOverlapping("ooooo", "o", "o")
	fmt.Println(is)
}

func TestNamePtrTree_GetLongestBackwardSubsequence(t *testing.T) {
	tree := &NamePtrTree{root: &PtrNode{
		key:      "",
		ptr:      0,
		children: []*PtrNode{
			//{
			//	key: "com",
			//	ptr: 5,
			//	children: []*PtrNode{
			//		{
			//			key: "baidu",
			//			ptr: 1,
			//			children: []*PtrNode{
			//				{
			//					key:      "www",
			//					ptr:      123,
			//					children: nil,
			//				},
			//			},
			//		},
			//	}},
			//{
			//	key:      "net",
			//	ptr:      2234,
			//	children: nil,
			//},
		},
	}}

	res := tree.GetLongestBackwardSubsequence(strings.Split("www.baidu.com", "."))
	fmt.Println(res)
}
