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

	// [-5 7 1 -1 -6 -2 7 -6 -9 -1]
	// [-9.1 6.891 -2.512 -8.835 6.244 -9.5 2.727 4.054 -3.217 -5.947]
	// [-1 6 -5 -4 2 -3 -6 3 6 0]
	// [5.765 2.854 -6.428 2.232 -3.812 -6.636 -8.834 -3.799 7.508 1.391]
	// [7 -1 -7 -1 4 6 -4 4 2 2]
	// [-6.696 -0.248 -3.602 4.661 7.863 -0.165 5.557 2.313 -8.746 4.214]
}

func TestGenerateRandomBytes(t *testing.T) {
	for range 3 {
		fmt.Println(string(GenerateRandomBytes(20)))

		//kTGIilO7A3j8OMqPV7lp
		//b4dAz7JLxqgChukX7ooa
		//bA8qe0CM7KMtf3LWvqIo
	}
}
