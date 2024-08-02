package main

import (
	"fmt"
	"sync"
	"time"
)

// producer is a function that simulates a producer of messages.
// It sends messages to the messages channel and signals when it's done via the done channel.
func producer(name string, messages chan<- string, done chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		message := fmt.Sprintf("%s: Message %d", name, i)
		messages <- message // Send message to the channel
		fmt.Printf("%s sent: %s\n", name, message)
		time.Sleep(time.Millisecond * 500) // Simulate some work
	}
	done <- name // Signal that producer is done
}

// consumer is a function that consumes messages from the messages channel.
func consumer(messages <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range messages {
		fmt.Println("Received:", msg)
		time.Sleep(time.Millisecond * 1000) // Simulate processing time
	}
}

func main() {
	// Create buffered channels for communication
	channelA := make(chan string, 3)
	channelB := make(chan string, 2)
	channelC := make(chan string, 4)

	// Channel to synchronize producer completion
	done := make(chan string)

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Start producer goroutines
	wg.Add(3)
	go producer("Producer A", channelA, done, &wg)
	go producer("Producer B", channelB, done, &wg)
	go producer("Producer C", channelC, done, &wg)

	// Start consumer goroutines
	wg.Add(3)
	go consumer(channelA, &wg)
	go consumer(channelB, &wg)
	go consumer(channelC, &wg)

	// Wait for all producers to finish
	go func() {
		wg.Wait()
		close(done)
	}()

	// Process completion signals from producers
	for i := 0; i < 3; i++ {
		producerName := <-done
		fmt.Println(producerName, "has finished producing messages.")
	}

	// Print message indicating the main goroutine exits
	fmt.Println("Main goroutine exits")
}
