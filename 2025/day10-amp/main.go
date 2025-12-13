package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Welcome to Day 10!!")
	run("input.txt")
}

func run(input string) {
	fmt.Println("run!!")
	case1()
}

// formatIntsAsBinary takes a slice of integers and returns a space-separated
// string of their 4-bit zero-padded binary representations.
func formatIntsAsBinary(inputs []int) string {
	binaryParts := []string{}
	for _, num := range inputs {
		// %04b formats the integer as a binary string (b), with a minimum
		// width of 4, zero-padded (0 prefix).
		binaryParts = append(binaryParts, fmt.Sprintf("%04b", num))
	}
	return strings.Join(binaryParts, " ")
}

func showButton() {
	// Example input: an integer slice representing the values 1, 5, 2, 3, 10, 12
	// Note: these must be within the range representable by 4 bits (0 to 15)
	inputData := []int{1, 5, 2, 3, 10, 12}

	for _, i := range inputData {
		out := fmt.Sprintf("%04b", i)
		fmt.Println("INT:", i, "OUT:", out)
	}
	// Get the formatted binary string
	binaryOutput := formatIntsAsBinary(inputData)

	// Print the final output
	fmt.Printf("Button: %s\n", binaryOutput)
}

func case1() {
	fmt.Println("case1!!")
	showButton()
	// [.##.] (3) (1,3) (2) (2,3) (0,2) (0,1)
	// Start with target and work backwards - 0110 == 6
	// Button: 0001 0101 0010 0011 1010 1100

	// Ends with 6, always start with 0
	// Depth search 4; call each button combo; recursive?
}
