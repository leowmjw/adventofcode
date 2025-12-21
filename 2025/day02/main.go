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
	//part1("input.txt")

	//part2("test.txt")
	part2("part1.txt")
}

func findRepeatedAtLeastTwice(start, finish int) int {
	// Start from slice of 1 till .. max half ..
	// Start with pair of 1; 2; 3 ..
	fmt.Println("Finding dupes between", start, "and", finish)
	invalidSum := 0
	for i := start; i <= finish; i++ {
		// Split into half by len which wiuld be max
		// If len is not divisible by 2 .. skip ..
		strascii := strconv.Itoa(i)
		strlen := len(strascii)
		fmt.Println("LEN: ", strlen)

		var left, right string

		// same as part1; but we only extract out left + right ..
		if strlen%2 == 0 {
			fmt.Println("MULTIPLE of 2: ", i)
			left = strascii[:strlen/2]
			right = strascii[strlen/2:]
		} else if strlen%3 == 0 {
			fmt.Println("MULTIPLE of 3: ", i)
			// left is the first 3rd ..
			left = strascii[:strlen/3]
			if strlen != 3 {
				right = strascii[strlen/3 : strlen/3+3]
			}
		} else if strlen%5 == 0 {
			fmt.Println("MULTIPLE of 5: ", i)
			left = strascii[:strlen/5]
			if strlen != 5 {
				right = strascii[strlen/5 : strlen/5+5]
			}
		} else if strlen%7 == 0 {
			fmt.Println("MULTIPLE of 7: ", i)
			left = strascii[:strlen/7]
			if strlen != 7 {
				right = strascii[strlen/7 : strlen/7+7]
			}
		} else {
			continue
		}
		fmt.Println("LEFT:", left, " RIGHT: ", right)
		// If get here can now check if left == right; meaning at least 2 match!
		if left == right {
			fmt.Println("Found dupes: ", strascii)
			invalidSum += i
		} else if strlen > 2 {
			// Check all the rest of the numbers in the
			for j := 1; j < strlen; j++ {
				if strascii[0] != strascii[j] {
					// quick break
					fmt.Println("NOT SPECIAL CASE: Skip!!", strascii)
					goto skipdupe
				}
			}
			// got this far means they are all the same!!
			fmt.Println("SPECIAL SINGLE CASE: Found dupes: ", strascii)
			invalidSum += i
		skipdupe:
		}

	}
	fmt.Println("Invalid Sum: ", invalidSum)
	fmt.Println("Done ==========>>")

	return invalidSum
}

func part2(input string) {
	_, serr := script.File(input).FilterLine(func(line string) string {

		if line[0] == '#' {
			return ""
		}

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
			invalidSum += findRepeatedAtLeastTwice(start, finish)
			fmt.Println("SUM: ", invalidSum)
		}
		fmt.Println("TOTAL Invalid Sum: ", invalidSum)

		// take action here ..
		return ""
	}).Slice()

	if serr != nil {
		panic(serr)
	}

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
}
