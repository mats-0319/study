package test

import (
	"testing"

	"github.com/mats9693/study/go/gocts/parse"
	"github.com/mats9693/study/go/gocts/utils"
)

func Test(t *testing.T) {
	for i := range -3 {
		print(i)
	}
}

// 正则匹配结果的结构：
// 三维数组：[][][]byte，其中最后的[]byte表示匹配到的字符数组
// 而第一维表示匹配到的每一次（其长度表示总共匹配到多少次）
// 第二维存储了单次匹配结果：第一个元素是匹配到的整个字符数组，第二个元素开始是正则中的括号提取的元素
func TestRequestRE(t *testing.T) {
	str := []byte(`const URI_GetList = "/list/get"`)

	re := parse.RequestRE
	reRes := re.FindAllSubmatch(str, -1)
	// reRes[0][1]: first '()' reRes[0][2]: second '()'

	for i := range reRes {
		for j := range reRes[i] {
			t.Log("res", i, j, string(reRes[i][j]))
		}
	}
}

func TestStructureRE(t *testing.T) {
	type ResBase struct {
		IsSuccess bool   `json:"is_success"`
		Err       string `json:"err"`
	}

	type User struct{}

	// ListUserRes struct comment
	// comment line 2
	type ListUserRes struct {
		Res    *ResBase `json:"res"`              // this is a comment
		Amount string   `json:'amount' yaml:"id"` // this is another comment
		Users  []*User  `json:"users"`
		*User  `json:"u"`
		ResBase
	}

	var str = []byte("// ListUserRes struct comment\n\t// comment line 2\n\ttype ListUserRes struct {\n\t\tRes    *ResBase `json:\"res\"`              // this is a comment\n\t\tAmount string   `json:'amount' yaml:\"id\"` // this is another comment\n\t\tUsers  []*User  `json:\"users\"`\n\t\t*User  `json:\"u\"`\n\t\tResBase\n\t}")

	re := parse.StructRE
	reRes := re.FindAllSubmatch(str, -1)

	for i := range reRes {
		for j := range reRes[i] {
			t.Log("res", i, j, string(reRes[i][j]))
		}
	}

	t.Log("---")

	fieldSlice := utils.BytesSplit(reRes[0][3], '\n', ';')

	for i, v := range fieldSlice {
		re = parse.StructFieldRE
		r := re.FindSubmatch(v)

		for j := range r {
			t.Log("res", i, j, string(r[j]))
		}
	}
}

func TestEnumRE(t *testing.T) {
	str := []byte(`
type UserIdentify int8

const (
	UserIdentify_Administrator UserIdentify = 0
	UserIdentify_VIP           UserIdentify = 1
	UserIdentify_Visitor       UserIdentify = 2
)
`)

	re := parse.EnumRE
	reRes := re.FindAllSubmatch(str, -1)

	for i := range reRes {
		for j := range reRes[i] {
			t.Log("res", i, j, string(reRes[i][j]))
		}
	}

	t.Log("---")

	for i := range reRes {
		if len(reRes[i]) < 4 {
			continue
		}

		re = parse.EnumUnitRE
		reRes = re.FindAllSubmatch(reRes[i][3], -1)
		break
	}

	for i := range reRes {
		for j := range reRes[i] {
			t.Log("res", i, j, string(reRes[i][j]))
		}
	}
}
