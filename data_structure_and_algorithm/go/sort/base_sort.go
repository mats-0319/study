package sort

func BaseSort(data []int) {
	if len(data) < 2 {
		return
	}

	sortTimes := calcSortTimes(data)

	for i := range sortTimes {
		matrix := make([][]int, 10)

		for j := range data {
			item := data[j]

			v := 0
			for k := 0; k <= i && item != 0; k++ {
				v = item % 10
				item /= 10
			}

			matrix[v] = append(matrix[v], data[j])
		}

		index := 0
		for j := range matrix {
			for k := range matrix[j] {
				data[index] = matrix[j][k]
				index++
			}
		}

		compareItems := 0
		for j := 1; j < len(matrix); j++ {
			compareItems += len(matrix[j])
		}
		if compareItems < 2 {
			break // already sorted
		}
	}
}

func calcSortTimes(data []int) int {
	maxItem := data[0]
	for i := range data {
		if data[i] > maxItem {
			maxItem = data[i]
		}
	}

	sortTimes := 0
	for maxItem != 0 {
		maxItem /= 10
		sortTimes++
	}

	return sortTimes
}
