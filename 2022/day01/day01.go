package main

import (
	"fmt"
	"github.com/bitfield/script"
	"sort"
	"strconv"
)

func main() {
	fmt.Println("Day01 ..")
	//SolvePart1()
	SolvePart2()
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
	// Do naive first; full in-mem + sort
	var totalPerElf []int
	var currentElfTotal int
	for _, v := range input {
		// Exit condition ..
		if v == "" {
			// Evaluate .. + exit condition ..
			// Add them all up
			// DEBUG
			//fmt.Println("CURRENT_TOTAL:", currentElfTotal)
			// When transitioning to new Elf; add to slice
			totalPerElf = append(totalPerElf, currentElfTotal)
			// reset for each Elf after append to slice
			currentElfTotal = 0
		} else {
			n, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			// DEBUG
			//fmt.Println("N:", n)
			// Sum up ..
			currentElfTotal += n
		}
	}
	// Edge case: Last item
	if currentElfTotal > 0 {
		totalPerElf = append(totalPerElf, currentElfTotal)
	}
	// Sort; pick up the highest
	sort.Sort(sort.Reverse(sort.IntSlice(totalPerElf)))

	for _, v := range totalPerElf[0:3] {
		// DEBUG
		//fmt.Println("TOP_3:", v)
		highestTotal += v
	}

	return highestTotal
}
