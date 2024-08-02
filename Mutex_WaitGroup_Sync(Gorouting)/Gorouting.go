package main

import (
	"fmt"
	"sync"
	"time"
)

var m = sync.RWMutex{}    // RWMutex for protecting shared data
var wg = sync.WaitGroup{} // WaitGroup to synchronize goroutines

var dbData = []string{"id1", "id2", "id3", "id4", "id5"} // Data to be processed
var results = []string{}                                 // Slice to store results

func main() {
	// Record start time
	startTime := time.Now()

	// Iterate through the dbData
	for i := 0; i < len(dbData); i++ {
		wg.Add(1)    // Add goroutine to WaitGroup
		go dbCall(i) // Call dbCall concurrently
	}

	wg.Wait()                                                       // Wait for all goroutines to finish
	fmt.Printf("\nTotal execution time: %v", time.Since(startTime)) // Print total execution time
	fmt.Printf("\nThe results are %v", results)                     // Print results
}

// Simulates a DB call with a delay
func dbCall(i int) {
	// Simulate DB call delay
	var delay float32 = 2000
	time.Sleep(time.Duration(delay) * time.Millisecond)

	// Save result
	save(dbData[i])

	// Log result
	Log()

	wg.Done() // Notify WaitGroup that the goroutine is done
}

// Saves result to the results slice, protected by a mutex
func save(result string) {
	m.Lock() // Lock mutex before accessing shared data
	results = append(results, result)
	m.Unlock() // Unlock mutex after accessing shared data
}

// Logs current results, read-locked for concurrency safety
func Log() {
	m.RLock()         // Read lock mutex for concurrent access
	defer m.RUnlock() // Ensure mutex is unlocked when function exits
	fmt.Printf("\nThe current results are: %v", results)
}
