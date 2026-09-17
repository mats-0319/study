package sort

import (
	"fmt"
)

func ExampleBubbleSort() {
	data := []int{3, 1, 5, 2, 4}
	BubbleSort(data)
	fmt.Println(data)

	// Output:
	// [1 2 3 4 5]
}

func ExampleSelectionSort() {
	data := []int{3, 1, 5, 2, 4}
	SelectionSort(data)
	fmt.Println(data)

	// Output:
	// [1 2 3 4 5]
}

func ExampleInsertionSort() {
	data := []int{3, 1, 5, 2, 4}
	InsertionSort(data)
	fmt.Println(data)

	// Output:
	// [1 2 3 4 5]
}

func ExampleMergeSort() {
	data := []int{3, 1, 5, 2, 4}
	MergeSort(data)
	fmt.Println(data)

	// Output:
	// [1 2 3 4 5]
}

func ExampleQuickSort() {
	data := []int{3, 1, 5, 2, 4}
	QuickSort(data)
	fmt.Println(data)

	// Output:
	// [1 2 3 4 5]
}

func ExampleBaseSort() {
	data := []int{3, 1, 5, 2, 4}
	BaseSort(data)
	fmt.Println(data)

	// Output:
	// [1 2 3 4 5]
}

func ExampleHeapSort() {
	data := []int{3, 1, 5, 2, 4}
	HeapSort(data)
	fmt.Println(data)

	// Output:
	// [1 2 3 4 5]
}
