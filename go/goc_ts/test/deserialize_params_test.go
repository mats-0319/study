package test

import (
	"encoding/json"
	"testing"
)

func TestDeserializeFromJson(t *testing.T) {
	jsonStr := `{"operator":"mario","page":{"page_size":10,"page_num":1}}`
	res := &struct {
		Operator string `json:"operator"`
		Page     *struct {
			PageSize int `json:"page_size"`
			PageNum  int `json:"page_num"`
		} `json:"page"`
	}{}

	// 前端传json类型的参数，后端解析的时候直接丢一个对应类型的0值结构体就行了，结构体里嵌套的结构体程序会自行创建，不需要我们提前创建好
	err := json.Unmarshal([]byte(jsonStr), &res)
	if err != nil {
		t.Errorf("json unamrshal failed, error: %v", err)
	}

	t.Log(res, res.Page)
}
