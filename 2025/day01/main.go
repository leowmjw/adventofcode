package main

import (
	"fmt"
	"strconv"

	"github.com/bitfield/script"
	"github.com/davecgh/go-spew/spew"
)

func main() {

	fmt.Println("Welcome to AOC 2025 Day01!!")
	// Always starts at 50; with dial from 0 - 99 (100 items)
	run("input.txt")
}

func run(input string) {

	// If action is 'L'; then compare if eq, less, more ..
	s, serr := script.File(input).FilterLine(func(s string) string {
		// Split command; Left is negative
		if s[0] == 'L' {
			s = "-" + s[1:]
		} else {
			s = "+" + s[1:]
		}
		return s
	}).Slice()
	if serr != nil {
		panic(serr)
	}

	// Find out how far from the edge; start of 50; it is 44 from edge both side
	// howFarFromHundred ...
	spew.Dump(s)

	count := 0
	current := 50
	// Rule is simpler .. no need modulus
	// If it reaches exactly 0; the increase the counter
	for _, action := range s {
		fmt.Println(action)

		step, cerr := strconv.Atoi(action)
		if cerr != nil {
			panic(cerr)
		}
		current += step
		current = current % 100
		fmt.Println("Current after step: ", current)
		if current == 0 {
			count++
		}
	}
	fmt.Println("Count: ", count)
	// Below might be a more sophisitcated one .,
	// Do nothing; continue
	// If current + action < 100 and > 0
	//      continue ..
	// Special cases (increase counter):
	// If current + action == 100
	// If current + action == 0
	// If current + action > 100
	//		left over modules 100?
	//		left over divide 100 += counter
	// If current + action < 0
	//		left over modules 0?
	//		left over divide 100 += counter

}

func howFarFromHundred(direction string, current int) int {
	if direction == "R" {
		return 100 - current
	}
	// else is current ..
	return current
}
