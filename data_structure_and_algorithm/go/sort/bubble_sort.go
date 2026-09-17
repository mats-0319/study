package sort

func BubbleSort(data []int) {
	if len(data) < 2 { // already sorted
		return
	}

	for i := range data { // i: loop times
		for j := 1; j < len(data)-i; j++ { // ignore last 'i' items, they are sorted
			if data[j-1] > data[j] {
				data[j-1], data[j] = data[j], data[j-1]
			}
		}
	}
}
