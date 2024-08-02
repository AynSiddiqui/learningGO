package main

import (
	"fmt"
	"time"
)

// producer is a function that simulates a producer of messages.
// It sends messages to the messages channel and signals when it's done via the done channel.
func producer(name string, messages chan<- string, done chan<- bool) {
	for i := 0; i < 5; i++ {
		// Create a message
		message := fmt.Sprintf("%s: Message %d", name, i)
		// Send message to the channel
		messages <- message
		// Print the message sent
		fmt.Printf("%s sent: %s\n", name, message)
		// Simulate some work
		time.Sleep(time.Millisecond * 500)
	}
	// Signal that producer is done
	done <- true
}

// consumer is a function that consumes messages from the messages channel.
func consumer(messages <-chan string) {
	// Loop indefinitely, reading messages from the channel
	for msg := range messages {
		// Process the message
		fmt.Println("Received:", msg)
		// Simulate processing time
		time.Sleep(time.Millisecond * 1000)
	}
}

func main() {
	// Create two buffered channels
	channelA := make(chan string, 3) // Buffer size 3
	channelB := make(chan string, 2) // Buffer size 2

	// Channel to signal when producer finishes
	done := make(chan bool)

	// Start producer goroutines
	go producer("Producer A", channelA, done)
	go producer("Producer B", channelB, done)

	// Start consumer goroutine
	go consumer(channelA)
	go consumer(channelB)

	// Wait for producers to finish
	<-done // Wait for first producer
	<-done // Wait for second producer

	// Close channels to signal no more messages
	close(channelA)
	close(channelB)

	// Wait for consumer to finish processing
	time.Sleep(time.Second)

	// Print message indicating the main goroutine exits
	fmt.Println("Main goroutine exits")
}
