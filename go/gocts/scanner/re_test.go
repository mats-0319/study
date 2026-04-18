package scanner

import (
	"strings"
	"testing"
)

// FindAllSubmatch 的匹配结果：[][][]byte
// - 第三维：表示匹配到了多少次
// - 第二维：保存每一次匹配到的具体信息，包括使用`()`捕捉的部分
// - 第一维：[]byte，字节数组，可以视作string
//
// 第二维具体保存了哪些信息：
// - 第一个元素：匹配到的完整字符串
// - 后续元素：使用`()`捕捉的内容，按照括号从左到右的顺序

func TestRequestRE(t *testing.T) {
	str := []byte(`const URI_ListUser = "/user/list"
			const URI_CheckToken= "/user/check-token"
			const URI_CheckToken= '/user/check_token'`) // '\w'包括下划线

	matched := requestRE.FindAllSubmatch(str, -1)
	printMatchedSlice(t, matched)
}

func TestEnumRE(t *testing.T) {
	str := []byte(`type EnumTypeName int8
		const (
			EnumTypeName_Placeholder EnumTypeName = -1 
			EnumTypeName_First                    = 1
		)`)

	matched := enumRE.FindAllSubmatch(str, -1)
	printMatchedSlice(t, matched)

	if len(matched) < 1 {
		return
	}

	for _, m := range matched {
		matched = enumUnitRE.FindAllSubmatch(m[2], -1)
		printMatchedSlice(t, matched)
	}
}

func TestStructRE(t *testing.T) {
	str := []byte(`
// This is a cross-
// lines comment
// before a struct
type StructName struct {
	Field1 int "json:'field1'" // 这里json tag周围应当使用反引号，使用引号是因为测试用例外围需要反引号
	TypeName1 // embedded field
	TypeName2 "json:'type_name_2'" // anonymous field
}`)

	matched := structRE.FindAllSubmatch(str, -1)
	printMatchedSlice(t, matched)

	if len(matched) < 1 {
		return
	}

	for _, m := range matched {
		matched = structFieldRE.FindAllSubmatch(m[3], -1)
		printMatchedSlice(t, matched)
	}
}

func printMatchedSlice(t *testing.T, matched [][][]byte) {
	t.Logf("Matched: (len = %d)\n", len(matched))

	for i, subMatch := range matched {
		for j, capture := range subMatch {
			// compress and print
			vFields := strings.Fields(string(capture))
			v := strings.Join(vFields, " ")
			t.Logf("  %d %d: %s", i, j, v)
		}
	}

	t.Log()
}
