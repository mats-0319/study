package sort

import "github.com/mats0319/study/suanfadaolun/data_strucure"

func HeapSort(data []int) {
	h := data_strucure.NewSmallHeap(data)
	sorted := h.Print()

	for i := range data {
		data[i] = sorted[i]
	}
}
