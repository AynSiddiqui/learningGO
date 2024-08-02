package main

import "fmt"

func main() {
	// Create a channel of type int
	var c = make(chan int)

	// Start a goroutine to process the channel
	go process(c)

	// Receive values from the channel until it's closed
	for i := range c {
		fmt.Println(i)
	}
}

// process is a function that sends values to a channel
func process(c chan int) {
	// Defer closing the channel when this function exits
	defer close(c)

	// Send integers 0 to 4 to the channel
	for i := 0; i < 5; i++ {
		c <- i
	}
}
