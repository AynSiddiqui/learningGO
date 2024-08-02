package main

import (
	"encoding/json" // Package for encoding and decoding JSON data
	"fmt"           // Package for formatted I/O
)

// Define the struct for contact information
type contactInfo struct {
	Name  string `json:"name"`  // Name of the contact
	Email string `json:"email"` // Email address of the contact
}

// Define the struct for purchase information
type purchaseInfo struct {
	Name   string  `json:"name"`   // Name of the purchased product
	Price  float32 `json:"price"`  // Price of the purchased product
	Amount int     `json:"amount"` // Quantity of the purchased product
}

func main() {
	// JSON data for contact information
	contactData := `
    [
        {"name": "John Doe", "email": "john@example.com"},
        {"name": "Jane Smith", "email": "jane@example.com"}
    ]
    `

	// JSON data for purchase information
	purchaseData := `
    [
        {"name": "Product A", "price": 10.99, "amount": 2},
        {"name": "Product B", "price": 20.50, "amount": 1}
    ]
    `

	// Load contact information from JSON string
	var contacts []contactInfo
	if err := json.Unmarshal([]byte(contactData), &contacts); err != nil {
		fmt.Println("Error unmarshalling contact JSON:", err)
		return
	}
	fmt.Printf("\nContact Info: %+v\n", contacts)

	// Load purchase information from JSON string
	var purchases []purchaseInfo
	if err := json.Unmarshal([]byte(purchaseData), &purchases); err != nil {
		fmt.Println("Error unmarshalling purchase JSON:", err)
		return
	}
	fmt.Printf("\nPurchase Info: %+v\n", purchases)

	// write json marshalling simple example
	i, _ := json.Marshal(contacts)
	fmt.Println("Marshal->>>>>>", string(i))
}
