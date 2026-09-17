package sort

func InsertionSort(data []int) {
	if len(data) < 2 {
		return
	}

	for i := 1; i < len(data); i++ { // i: next index after sorted area
		value := data[i]

		j := i - 1 // left index of right position
		for ; j >= 0 && value < data[j]; j-- {
			data[j+1] = data[j]
		}
		data[j+1] = value
	}
}
