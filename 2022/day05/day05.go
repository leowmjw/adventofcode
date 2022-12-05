package main

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/davecgh/go-spew/spew"
	"strings"
)

func main() {

	fmt.Println("Day 05 ..")
	Part1()

}

func Part1() {
	data, err := script.File("testdata/sample.txt").Slice()
	if err != nil {
		panic(err)
	}
	for _, l := range data {
		cols := strings.ReplaceAll(l, " ", "@")
		//cols := strings.Split(l, " ")
		spew.Dump(cols)
		//fmt.Println("COLS:", len(cols))
		// If see '1'; is a border; last number is x * n
		// extract + skip 4 everytime ..

		// Now in the instructions mode
	}

	// From raw data; reverse = push into the stack
	// Or from top now; from multiple, pick up the data
	// and push there and then ..
}
