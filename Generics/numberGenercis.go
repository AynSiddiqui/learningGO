package main

import "fmt"

func main() {
	// Define an integer slice
	var intSlice = []int{1, 2, 3}

	// Call sumSlice function for integers and print the result
	fmt.Println("Sum of integers:", sumSlice[int](intSlice))

	// Define a float32 slice
	var float32Slice = []float32{1.1, 2.2, 3.3}

	// Call sumSlice function for float32 and print the result
	fmt.Println("Sum of float32:", sumSlice[float32](float32Slice))
}

// sumSlice function calculates the sum of elements in a slice of any numeric type
func sumSlice[T int | float32 | float64](slice []T) T {
	// Initialize sum variable with zero value of type T
	var sum T

	// Iterate over the slice elements
	for _, v := range slice {
		// Add each element to the sum
		sum += v
	}

	// Return the calculated sum
	return sum
}
