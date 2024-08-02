package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Constants for maximum prices
var MAX_CHICKEN_PRICE float32 = 5
var MAX_TOFU_PRICE float32 = 3

func main() {
	// Channels for sending chicken and tofu deal notifications
	var chickenChannel = make(chan string)
	var tofuChannel = make(chan string)

	// List of websites to check for deals
	var websites = []string{"walmart.com", "costco.com", "wholefoods.com"}

	// Iterate through websites to check prices concurrently
	for i := range websites {
		// Start goroutines to check prices for chicken and tofu
		go checkChickenPrices(websites[i], chickenChannel)
		go checkTofuPrices(websites[i], tofuChannel)

		// Send notifications when a deal is found
		sendMessage(chickenChannel, tofuChannel)
	}
}

// Function to check tofu prices on a website
func checkTofuPrices(website string, c chan string) {
	for {
		// Simulate delay before checking prices
		time.Sleep(time.Second * 1)

		// Generate a random tofu price
		var tofu_price = rand.Float32() * 20

		// If the price is below the maximum, send website name on the channel
		if tofu_price < MAX_TOFU_PRICE {
			c <- website
			break // Exit loop after sending the message
		}
	}
}

// Function to check chicken prices on a website
func checkChickenPrices(website string, chickenChannel chan string) {
	for {
		// Simulate delay before checking prices
		time.Sleep(time.Second * 1)

		// Generate a random chicken price
		var chickenPrice = rand.Float32() * 20

		// If the price is below the maximum, send website name on the channel
		if chickenPrice <= MAX_CHICKEN_PRICE {
			chickenChannel <- website
			break // Exit loop after sending the message
		}
	}
}

// Function to send notifications when a deal is found
func sendMessage(chickenChannel chan string, tofuChannel chan string) {
	// Listen for messages on both chicken and tofu channels
	select {
	case website := <-chickenChannel:
		// If a message is received on the chicken channel, print a text notification
		fmt.Printf("\nText Sent: Found deal on chicken at %v.", website)
	case website := <-tofuChannel:
		// If a message is received on the tofu channel, print an email notification
		fmt.Printf("\nEmail Sent: Found deal on tofu at %v. ", website)
	}
}
