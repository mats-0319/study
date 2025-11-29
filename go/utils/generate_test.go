package utils

import (
	"fmt"
	"testing"
)

func TestGenerateRandomIntSlice(t *testing.T) {
	for range 3 {
		fmt.Println(GenerateRandomSlice[int](10, 10))
		fmt.Println(GenerateRandomSlice[float64](10, 10.0))
	}

	// [8 5 6 -7 -6 -5 2 7 9 -6]
	// [0.65 -4.45 8.539 -2.607 1.322 5.526 7.381 5.125 -2.186 -1.234]
	// [5 2 -7 7 -4 0 -3 -1 8 -3]
	// [-9.507 2.702 3.012 -1.703 -6.763 4.521 -7.792 3.381 -0.714 6.662]
	// [3 -1 -4 -5 -3 4 -2 -1 -8 -9]
	// [6.961 2.952 -7.435 9.487 5.468 -0.123 -4.897 8.305 4.765 -6.441]
}

func TestGenerateRandomBytes(t *testing.T) {
	for range 3 {
		fmt.Println(string(GenerateRandomBytes(20)))
	}

	//kTGIilO7A3j8OMqPV7lp
	//b4dAz7JLxqgChukX7ooa
	//bA8qe0CM7KMtf3LWvqIo
}
