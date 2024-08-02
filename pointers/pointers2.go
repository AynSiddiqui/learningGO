package main

import "fmt"

func main() {
	var p *int32 = new(int32)
	var q **int32 // Pointer to a pointer
	var i int32 = 32

	q = &p // Assigning the address of p to q
	intput := "This is a string"
	output := []string{intput}
	fmt.Println(intput)
	fmt.Println(output)
	fmt.Printf("\nThe value of i is: %v", i)
	p = &i
	fmt.Printf("\nThe value p points to is: %v", *p)
	// Printing **q *q &q and q will give you the addresses accordingly
	fmt.Printf("\nThe value q points to is: %v", q) // Dereferencing q twice to get the value of i
}
