package utils

import "slices"

// Set 字符串的集合，因为go暂时不支持泛型方法（已采纳提案、未实装）所以没有写成泛型结构
// 允许往集合中添加数据，集合中没有的会被保存
type Set struct {
	Data []string
}

// Add 因为在未初始化的情况下会返回新的实例，所以应将函数返回值重新赋值给调用者，形如`s = s.add("123")`
func (s *Set) Add(value string) *Set {
	if s == nil || s.Data == nil { // not initialized
		return &Set{
			Data: []string{value},
		}
	}

	if !slices.Contains(s.Data, value) {
		s.Data = append(s.Data, value)
	}

	return s
}
