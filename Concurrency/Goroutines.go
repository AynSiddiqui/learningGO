package main

import (
	"fmt"
	"time"
)

func display(str string) {
	for w := 0; w < 6; w++ {
		fmt.Println(str, "HELOOOOOO:", w)
		time.Sleep(100 * time.Millisecond) // Simulate some work
	}
}

func d(str string) {
	for w := 0; w < 6; w++ {
		fmt.Println(str, "AYAAN:", w)
		time.Sleep(100 * time.Millisecond) // Simulate some work
	}
}

func main() {
	fmt.Println("Start of main function")

	// Calling Goroutine
	go display("Welcome")
	go d("second")
	// Calling normal function
	display("GeeksforGeeks")

	fmt.Println("End of main function")
	time.Sleep(1 * time.Second) // Wait to see output
}
