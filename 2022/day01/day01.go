package main

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/davecgh/go-spew/spew"
	"strconv"
)

func main() {
	fmt.Println("Day01 ..")
	Solve()
}

func Solve() {
	// Load test data; this t=year try out script
	l, err := script.File("testdata/sample.txt").Slice()
	if err != nil {
		panic(err)
	}
	spew.Dump(l)
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
