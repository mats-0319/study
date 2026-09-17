package sort

func QuickSort(data []int) {
	if len(data) < 2 {
		return
	}

	index := 1 // set next small item into this index
	for i := index; i < len(data); i++ {
		if data[i] < data[0] {
			data[i], data[index] = data[index], data[i]
			index++
		}
	}

	index-- // last item smaller than 'data[0]'
	data[0], data[index] = data[index], data[0]

	// 直接递归，蹭开头的长度返回条件
	QuickSort(data[:index])
	QuickSort(data[index+1:])
}
