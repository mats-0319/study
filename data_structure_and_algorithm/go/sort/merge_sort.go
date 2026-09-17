package sort

func MergeSort(data []int) {
	if len(data) < 2 {
		return
	}

	middle := getMiddleIndex(data)

	sorted := merge(data[:middle], data[middle:])

	for i := range data {
		data[i] = sorted[i]
	}

	return
}

func merge(left []int, right []int) []int {
	if len(left) > 1 {
		middle := getMiddleIndex(left)
		left = merge(left[:middle], left[middle:])
	}
	if len(right) > 1 {
		middle := getMiddleIndex(right)
		right = merge(right[:middle], right[middle:])
	}

	res := make([]int, 0, len(left)+len(right))
	for len(left) > 0 && len(right) > 0 {
		if left[0] <= right[0] {
			res = append(res, left[0])
			left = left[1:]
		} else {
			res = append(res, right[0])
			right = right[1:]
		}
	}

	for _, v := range left {
		res = append(res, v)
	}
	for _, v := range right {
		res = append(res, v)
	}

	return res
}

func getMiddleIndex(v []int) int {
	return (len(v) + 1) / 2
}
