package data_strucure

import (
	"testing"
)

func TestHeapSort(t *testing.T) {
	heap := NewSmallHeap([]int{3, 1, 5, 2, 4})
	res := heap.Print()
	if !isEqual(res, []int{1, 2, 3, 4, 5}) {
		t.Fatal("heap sort failed, \nwant: [1,2,3,4,5]\nget : ", res)
	}
}

func isEqual(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	isEqualFlag := true
	for i := range a {
		if a[i] != b[i] {
			isEqualFlag = false
			break
		}
	}

	return isEqualFlag
}
