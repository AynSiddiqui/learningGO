package main

import (
	"fmt"
)

// producer function sends data to the channel
func producer(data chan<- int) {
	for i := 0; i < 5; i++ {
		data <- i // Send data to the channel
	}
	close(data) // Close the channel when done
}

// consumer function receives data from the channel
func consumer(data <-chan int, done chan<- bool) {
	for num := range data {
		fmt.Println("Received:", num)
	}
	done <- true // Signal that consumer is done
}

func main() {
	// Create a unidirectional channel for sending data from producer to consumer
	dataChannel := make(chan int)

	// Create a channel for signaling when consumer is done
	done := make(chan bool)

	// Start the producer goroutine
	go producer(dataChannel)

	// Start the consumer goroutine
	go consumer(dataChannel, done)

	// Wait for the consumer to finish
	<-done

	fmt.Println("Main goroutine exits")
}
