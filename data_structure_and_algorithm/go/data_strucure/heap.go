package data_strucure

// SmallHeap small-top heap, array type.
// deprecated '0' index and:
// - parent      index: i/2
// - left  child index: 2*i
// - right child index: 2*i + 1
type SmallHeap struct {
	Data []int
}

func NewSmallHeap(data []int) *SmallHeap {
	sh := &SmallHeap{}
	for i := range data {
		sh.add(data[i])
	}

	return sh
}

func (h *SmallHeap) Print() []int {
	items := len(h.Data) - 1

	res := make([]int, 0, items)
	for range items {
		res = append(res, h.pop())
	}

	return res
}

func (h *SmallHeap) add(value int) {
	if h.Data == nil {
		h.Data = make([]int, 1)
	}

	h.Data = append(h.Data, value)
	h.shiftUp()
}

func (h *SmallHeap) shiftUp() {
	if len(h.Data) < 3 { // 去掉忽略的0号索引，堆为空或只有一个元素时，认为堆是有序的，直接返回
		return
	}

	currentIndex := len(h.Data) - 1
	parentIndex := currentIndex / 2
	for currentIndex > 0 && h.Data[currentIndex] < h.Data[parentIndex] {
		h.Data[currentIndex], h.Data[parentIndex] = h.Data[parentIndex], h.Data[currentIndex]
		currentIndex = parentIndex
		parentIndex = currentIndex / 2
	}
}

func (h *SmallHeap) pop() int {
	if len(h.Data) < 2 {
		panic("can't pop from an empty heap")
	}

	res := h.Data[1]

	h.Data[1] = h.Data[len(h.Data)-1]
	h.Data = h.Data[:len(h.Data)-1]

	currentIndex := 1
	leftIndex := currentIndex * 2
	rightIndex := currentIndex*2 + 1
	for leftIndex < len(h.Data) {
		currentValue := h.Data[currentIndex]
		leftValue := h.Data[leftIndex]
		rightValue := getRightValue(h.Data, rightIndex, leftValue+1)

		if currentValue <= leftValue && currentValue <= rightValue {
			break
		}

		if currentValue > leftValue && currentValue > rightValue {
			if leftValue < rightValue {
				h.Data[currentIndex], h.Data[leftIndex] = h.Data[leftIndex], h.Data[currentIndex]
				currentIndex = leftIndex
			} else {
				h.Data[currentIndex], h.Data[rightIndex] = h.Data[rightIndex], h.Data[currentIndex]
				currentIndex = rightIndex
			}
		} else if currentValue > leftValue {
			h.Data[currentIndex], h.Data[leftIndex] = h.Data[leftIndex], h.Data[currentIndex]
			currentIndex = leftIndex
		} else if currentValue > rightValue { // 可以只写else，全写出来看着更清晰
			h.Data[currentIndex], h.Data[rightIndex] = h.Data[rightIndex], h.Data[currentIndex]
			currentIndex = rightIndex
		}

		leftIndex = currentIndex * 2
		rightIndex = currentIndex*2 + 1
	}

	return res
}

func getRightValue(data []int, rightIndex int, defaultValue int) (rightValue int) {
	if rightIndex < len(data) {
		rightValue = data[rightIndex]
	} else {
		rightValue = defaultValue
	}

	return
}
