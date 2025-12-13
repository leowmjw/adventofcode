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
	//part1("test.txt")
	//part1("part1.txt")
	part1("input.txt")
}

func findDupes(start, finish int) int {
	// Start from slice of 1 till .. max half ..
	// Start with pair of 1; 2; 3 ..
	fmt.Println("Finding dupes between", start, "and", finish)
	invalidSum := 0
	for i := start; i <= finish; i++ {
		// Split into half by len which wiuld be max
		// If len is not divisible by 2 .. skip ..
		strascii := strconv.Itoa(i)
		strlen := len(strascii)
		// DEBUG
		//fmt.Println(i, " LEN: ", strlen)
		if strlen%2 != 0 {
			//fmt.Println("Skipping ", i)
			continue
		}
		left := strascii[:strlen/2]
		right := strascii[strlen/2:]
		fmt.Println("LEFT:", left, " RIGHT: ", right)
		if left == right {
			fmt.Println("Found dupes: ", strascii)
			invalidSum += i
		}
	}
	fmt.Println("Invalid Sum: ", invalidSum)
	fmt.Println("Done ==========>>")

	return invalidSum
}

func splitCount(line string) {
	s := strings.Split(line, ",")
	invalidSum := 0
	for i := 0; i < len(s); i++ {
		ids := strings.Split(s[i], "-")
		// DEBUG
		//spew.Dump(ids)
		start, cerr := strconv.Atoi(ids[0])
		finish, cerr2 := strconv.Atoi(ids[1])
		if cerr != nil || cerr2 != nil {
			panic(cerr)
		}
		invalidSum += findDupes(start, finish)
		fmt.Println("SUM: ", invalidSum)
	}
	fmt.Println("TOTAL Invalid Sum: ", invalidSum)
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
