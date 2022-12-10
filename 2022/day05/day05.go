package main

import (
	"fmt"
	"github.com/bitfield/script"
	"github.com/emirpasic/gods/stacks/arraystack"
	"strings"
)

func main() {

	fmt.Println("Day 05 ..")
	//Part1("testdata/sample.txt")
	Part1("testdata/input.txt")

}

func Part1(filePath string) {
	data, err := script.File(filePath).Slice()
	if err != nil {
		panic(err)
	}
	numOfStack := 0
	executeAction := false
	var queueSetup [][]string
	var stacks []*arraystack.Stack
	for _, l := range data {
		if executeAction {
			var numBlocks, fromStack, toStack int
			//var SnumBlocks, SfromStack, StoStack string
			//fmt.Println("NOW START ACTZION ..")
			// DEBUG
			//spew.Dump(stacks)
			// Now in the instructions mode
			_, err := fmt.Sscanf(l, "move %d from %d to %d", &numBlocks, &fromStack, &toStack)
			if err != nil {
				//panic(err)
				continue
			}
			fid := fromStack - 1
			tid := toStack - 1
			// DEBUG
			//fmt.Println("No chars:", n)
			for i := 0; i < numBlocks; i++ {
				// DEBUG
				//fmt.Println("POP_FROM:", fid, "TO:", tid)
				if v, ok := stacks[fid].Pop(); ok {
					stacks[tid].Push(v)
				}
			}

		} else {
			cols := strings.Split(l, "")
			// Check exit condition; a '1'
			// If see '1'; is a border; last number is x * n
			if cols[1] == "1" {
				// DEBUG
				//spew.Dump(queueSetup)
				executeAction = true
				stacks = make([]*arraystack.Stack, numOfStack)
				// Load from Slice into Stacks ..
				for sid, blocks := range queueSetup {
					// DEBUG
					//spew.Dump(blocks)
					stacks[sid] = arraystack.New() // empty
					// Reverse the order
					for i := len(blocks) - 1; i >= 0; i-- {
						// Push it into the stack
						// DEBUG
						//fmt.Println("PUSH:", blocks[i], "into Stack:", sid)
						stacks[sid].Push(blocks[i])
					}
				}
				// DEBUG - FULL STACK
				//spew.Dump(stacks)
				continue
			}
			// Setup one time at the start of the cycle
			if numOfStack == 0 {
				numOfStack = len(cols)/4 + 1
				//fmt.Println("NO_STACK:", numOfStack)
				queueSetup = make([][]string, numOfStack)
			}
			// extract + skip 4 everytime ..
			//numOfStack = len(cols) % 4
			//// DEBUG
			////spew.Dump(cols)
			//fmt.Println("COLS:", numOfStack)
			// Extract out the block based on numOfStack
			for i := 0; i < numOfStack; i++ {
				cellContent := strings.TrimSpace(cols[i*4+1])
				// DEBUG
				// fmt.Print("BLOCK:", cellContent)
				// Push into your own stack .. if NOT empty
				if cellContent != "" {
					// DEBUG
					//fmt.Println("APPEND", cellContent, "into Slice", i)
					//spew.Dump(queueSetup)
					//if queueSetup[i] == nil {
					//	// init make
					//	queueSetup[i] = make([]string, 0)
					//} else {
					queueSetup[i] = append(queueSetup[i], cellContent)
					//}
				}
			}
			// DEBUG
			//fmt.Println("=================>")
		}
	}

	// From raw data; reverse = push into the stack
	// Or from top now; from multiple, pick up the data
	// and push there and then ..
	// DEBUG
	//spew.Dump(stacks)

	for i := 0; i < len(stacks); i++ {
		if v, ok := stacks[i].Peek(); ok {
			fmt.Print(v)
		}
	}
}
