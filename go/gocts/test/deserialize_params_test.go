package test

import (
	"encoding/json"
	"testing"
)

func TestDeserializeFromJson(t *testing.T) {
	jsonStr := `{"operator":"mario","page":{"page_size":10,"page_num":1}}`
	//jsonStr := `{"operator":"mario","page_size":10,"page_num":1}`

	type Page struct {
		PageSize int `json:"page_size"`
		PageNum  int `json:"page_num"`
	}
	res := &struct {
		Operator string `json:"operator"`
		Page     `json:"page"`
	}{}

	// 前端传json类型的参数，后端解析的时候直接丢一个对应类型的结构体就行了，即使是嵌套的结构体也能正常创建，
	// 这样来看，前后端之间使用json传参可以在参数解析的时候更方便
	// 对于内嵌结构体，没有json tag时正常组合，将内嵌结构体字段提高到与内嵌字段同级；存在json tag时，按照具名字段处理
	err := json.Unmarshal([]byte(jsonStr), &res)
	if err != nil {
		t.Errorf("json unamrshal failed, error: %v", err)
	}

	t.Log(res, res.Page)
}
