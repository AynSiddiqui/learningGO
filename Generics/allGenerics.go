package main

import "fmt"

func main() {
	// Define an empty integer slice
	var intSlice = []int{}

	// Call isEmpty function for intSlice and print the result
	fmt.Println("Is intSlice empty:", isEmpty[int](intSlice))

	// Define a non-empty float32 slice
	var float32Slice = []float32{1, 2, 3}

	// Call isEmpty function for float32Slice and print the result
	fmt.Println("Is float32Slice empty:", isEmpty[float32](float32Slice))
}

// isEmpty function checks if a slice is empty
func isEmpty[T any](slice []T) bool {
	// Return true if the length of the slice is zero, indicating it's empty
	return len(slice) == 0
}
