package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bitfield/script"
)

func main() {
	fmt.Println("Welcome to Day02!!")

	// Look for double numbers , hash
	// single digit; single digit ..
	// Start length and final length .
	// Samples: 11-22
	part1("test.txt")
}

func findDupes(start, finish int) {
	// Start from slice of 1 till .. max half ..
	// Start with pair of 1; 2; 3 ..
	fmt.Println("Finding dupes between", start, "and", finish)
	for i := start; i <= finish; i++ {
		// Split into half by len which wiuld be max
		fmt.Println(i)
	}
	fmt.Println("Done ==========>>")
}

func splitCount(line string) {
	s := strings.Split(line, ",")
	for i := 0; i < len(s); i++ {
		ids := strings.Split(s[i], "-")
		// DEBUG
		//spew.Dump(ids)
		start, cerr := strconv.Atoi(ids[0])
		finish, cerr2 := strconv.Atoi(ids[1])
		if cerr != nil || cerr2 != nil {
			panic(cerr)
		}
		findDupes(start, finish)
	}
}

func part1(input string) {
	_, serr := script.File(input).FilterLine(func(line string) string {
		splitCount(line)
		// take action here ..
		return ""
	}).Slice()

	if serr != nil {
		panic(serr)
	}
	//splitCount(s[0])
}
