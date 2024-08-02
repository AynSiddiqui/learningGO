package main

import (
	"bytes"
	"fmt"
)

func main() {

	// Creating and initializing
	// the slice of bytes
	// Using shorthand declaration
	slice_1 := []byte{'!', '!', 'G', 'e', 'e', 'e', 'e', 'k', 's',
		'f', 'o', 'r', 'G', 'e', 'e', 'k', 's', '#', '#'}

	slice_2 := []byte{'A', 'p', 'p', 'l', 'e'}

	slice_3 := []byte{'%', 'g', '%', 'e', '%', 'e',
		'%', 'k', '%', 's', '%'}

	// Displaying slices
	fmt.Println("Original Slice:")
	fmt.Printf("Slice 1: %s", slice_1)
	fmt.Printf("\nSlice 2: %s", slice_2)
	fmt.Printf("\nSlice 3: %s", slice_3)

	// Splitting the slice of bytes
	// Using Split function
	res1 := bytes.Split(slice_1, []byte("eek"))
	res2 := bytes.Split(slice_2, []byte(""))
	res3 := bytes.Split(slice_3, []byte("%"))

	// Display the results
	fmt.Printf("\n\nAfter splitting:")
	fmt.Printf("\nSlice 1: %s", res1)
	fmt.Printf("\nSlice 2: %s", res2)
	fmt.Printf("\nSlice 3: %s", res3)

}
