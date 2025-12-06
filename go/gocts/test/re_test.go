package test

import (
	"fmt"
	"testing"

	"github.com/mats9693/study/go/gocts/parse"
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
			fmt.Println("res", i, j, string(reRes[i][j]))
		}
	}
}

func TestStructureRE(t *testing.T) {
	//type Playlist struct {
	//	ID          string  `json:"id" yaml:"id"`     // this is a comment
	//	FileName    string  `json:"name" yaml:"name"` // this is another comment
	//	Description string  `json:"description" yaml:"description"`
	//	MusicList   []Music `json:"music_list" yaml:"music_list"`
	//}

	var str = []byte("type Playlist struct {\n\t\tID          string  `json:\"id\" yaml:\"id\"`     // this is a comment\n\t\tFileName    string  `json:\"name\" yaml:\"name\"` // this is another comment\n\t\tDescription string  `json:\"description\" yaml:\"description\"`\n\t\tMusicList   []Music `json:\"music_list\" yaml:\"music_list\"`\n\t}")

	re := parse.StructRE
	reRes := re.FindAllSubmatch(str, -1)

	for i := range reRes {
		for j := range reRes[i] {
			fmt.Println("res", i, j, string(reRes[i][j]))
		}
	}

	fmt.Println("---")

	for i := range reRes {
		if len(reRes[i]) < 3 {
			continue
		}

		re = parse.StructFieldRE
		reRes = re.FindAllSubmatch(reRes[i][2], -1)
		break
	}

	for i := range reRes {
		for j := range reRes[i] {
			fmt.Println("res", i, j, string(reRes[i][j]))
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
			fmt.Println("res", i, j, string(reRes[i][j]))
		}
	}

	fmt.Println("---")

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
			fmt.Println("res", i, j, string(reRes[i][j]))
		}
	}
}
