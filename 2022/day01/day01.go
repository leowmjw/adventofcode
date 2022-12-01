package main

import (
	"fmt"
	"github.com/bitfield/script"
	"strconv"
)

func main() {
	fmt.Println("Day01 ..")
	//SolvePart1()
	//SolvePart2()
}

func SolvePart1() {
	// Load test data; this t=year try out script
	l, err := script.File("./testdata/input.txt").Slice()
	if err != nil {
		panic(err)
	}
	fmt.Println(maxCalBySingleElf(l))
}

func maxCalBySingleElf(input []string) (highestTotal int) {
	currentCaloric := 0
	// Track highest calory - highestTotal
	for _, v := range input {
		// Exit condition ..
		if v == "" {
			// Evaluate .. + exit condition ..
			if currentCaloric > highestTotal {
				highestTotal = currentCaloric
			}
			// reset for each Elf
			currentCaloric = 0
		} else {
			n, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			currentCaloric += n
		}
	}
	return highestTotal
}

func SolvePart2() {
	// Load test data; this t=year try out script
	l, err := script.File("./testdata/input.txt").Slice()
	if err != nil {
		panic(err)
	}
	fmt.Println(maxCalByTop3Elves(l))
}

func maxCalByTop3Elves(input []string) (highestTotal int) {
	highestTotal = 4000
	return highestTotal
}

// Return new slice with the lowest removed
// slot in place? mark the one location smallest; slice it out
