package test

import (
	"fmt"
	"testing"
)

func TestReSizeInLoop(t *testing.T) {
	intSlice := []int{1, 2, 3}

	for i, v := range intSlice {
		if v == 2 {
			intSlice = append(intSlice, 4)
		}

		t.Log(fmt.Sprintf("index: %d, value: %d", i, v))
	}

	t.Log("length: ", len(intSlice))

	intSlice = []int{1, 2, 3}

	for i := 0; i < len(intSlice); i++ {
		if intSlice[i] == 2 {
			intSlice = append(intSlice, 4)
		}

		t.Log(fmt.Sprintf("index: %d, value: %d", i, intSlice[i]))
	}

	t.Log("length: ", len(intSlice))
}
