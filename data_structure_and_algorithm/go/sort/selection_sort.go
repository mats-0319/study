package sort

func SelectionSort(data []int) {
	if len(data) < 2 {
		return
	}

	for i := 0; i < len(data)-1; i++ { // i: next index after sorted area
		minIndex := i + 1
		for j := minIndex + 1; j < len(data); j++ {
			if data[j] < data[minIndex] {
				minIndex = j
			}
		}

		data[i], data[minIndex] = data[minIndex], data[i]
	}
}
