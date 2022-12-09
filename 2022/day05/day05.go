package main

import (
	"fmt"
	"github.com/bitfield/script"
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
	numofStack := 0
	executeAction := false
	for _, l := range data {
		if executeAction {
			var bob string
			fmt.Println("NOW START ACZtION ..")
			return
			// Now parse the actions .. instead of split to chars
			// Now in the instructions mode
			n, err := fmt.Scanf("boo %s", bob)
			if err != nil {
				panic(err)
			}
			fmt.Println("NO chsrs:", n)
		} else {
			cols := strings.Split(l, "")
			if numofStack > 0 {
				// Check exit condition; a '1'
				// If see '1'; is a border; last number is x * n
				if cols[1] == "1" {
					executeAction = true
					continue
				}
				// Extract out the block based on numofStack
				for i := 0; i < numofStack; i++ {
					cellContent := strings.TrimSpace(cols[i*4+1])
					// DEBUG
					// fmt.Print("BLOCK:", cellContent)
					// Push into your own stack .. if NOT empty
					if cellContent != "" {
						fmt.Println("PUSH", cellContent, "into Stack", i)
					}
				}
				fmt.Println("=================>")
			} else {
				// extract + skip 4 everytime ..
				numofStack = len(cols) % 4
				// DEBUG
				//spew.Dump(cols)
				fmt.Println("COLS:", numofStack)
				// Extract out the block based on numofStack
				for i := 0; i < numofStack; i++ {
					cellContent := strings.TrimSpace(cols[i*4+1])
					// DEBUG
					// fmt.Print("BLOCK:", cellContent)
					// Push into your own stack .. if NOT empty
					if cellContent != "" {
						fmt.Println("PUSH", cellContent, "into Stack", i)
					}
				}
				fmt.Println("=================>")
			}
		}
	}

	// From raw data; reverse = push into the stack
	// Or from top now; from multiple, pick up the data
	// and push there and then ..
}
