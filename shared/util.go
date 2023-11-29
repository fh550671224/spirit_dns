package shared

import "strings"

func SplitWithoutEmpty(s string, sep string) []string {
	list := strings.Split(s, sep)
	var res []string
	for _, v := range list {
		if v != "" {
			res = append(res, v)
		}
	}

	return res
}

func SplitWithSep(s string, sub string) []string {
	if sub == "" {
		return nil
	}

	index := strings.Index(s, sub)

	a := s[:index]
	b := s[index+len(sub):]

	res := []string{}
	if a != "" {
		res = append(res, a)
	}
	res = append(res, sub)
	if b != "" {
		res = append(res, b)
	}

	return res
}

func IsOverlapping(str, substr1, substr2 string) bool {
	// 找到两个子字符串在母字符串中的位置
	pos1 := strings.Index(str, substr1)
	pos2 := strings.Index(str, substr2)

	// 如果任何一个子字符串没有在母字符串中找到，返回false
	if pos1 == -1 || pos2 == -1 {
		return false
	}

	// 计算两个子字符串的结束位置
	endPos1 := pos1 + len(substr1)
	endPos2 := pos2 + len(substr2)

	// 检查子字符串的位置是否有重叠
	// 重叠的条件是第一个子字符串的结束位置在第二个子字符串的开始位置之后
	return (pos1 < pos2 && endPos1 > pos2) || (pos2 < pos1 && endPos2 > pos1)
}
