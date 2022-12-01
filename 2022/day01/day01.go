package main

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/davecgh/go-spew/spew"
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
	// Track highest calory
	highestTotal = 24000
	return highestTotal
}
