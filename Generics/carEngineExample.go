package main

import "fmt"

// Define the gasEngine struct
type gasEngine struct {
	gallons float32 // Gallons of fuel
	mpg     float32 // Miles per gallon efficiency
}

// Define the electricEngine struct
type electricEngine struct {
	kwh   float32 // Kilowatt-hours of electricity
	mpkwh float32 // Miles per kilowatt-hour efficiency
}

// Define the car struct with an engine field of type Engine interface
type car[T gasEngine | electricEngine] struct {
	carMake  string // Make of the car
	carModel string // Model of the car
	engine   T      // Engine of the car, which implements the Engine interface
}

func main() {
	// Create a gas car
	gasCar := car[gasEngine]{
		carMake:  "Honda",
		carModel: "Civic",
		engine: gasEngine{
			gallons: 12.4,
			mpg:     40,
		},
	}

	// Create an electric car
	electricCar := car[electricEngine]{
		carMake:  "Tesla",
		carModel: "Model 3",
		engine: electricEngine{
			kwh:   57.5,
			mpkwh: 4.17,
		},
	}
	fmt.Println("Gas car:", gasCar)
	fmt.Println("Electric car:", electricCar)
}
