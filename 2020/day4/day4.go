package day4

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func splitPassports(in io.ReadCloser) []string {

	allContent := []string{}

	scanner := bufio.NewScanner(in)
	fullLine := ""
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			// Create a new line to be appended
			allContent = append(allContent, fullLine)
			// reset
			fullLine = ""
		} else {
			// Needs space unless fullLine new
			if fullLine == "" {
				fullLine = line
			} else {
				fullLine = fullLine + " " + line
			}
		}
	}
	// Don;t forget last line ..
	allContent = append(allContent, fullLine)
	// DEBUG
	// spew.Dump(allContent)

	return allContent
}

func countPassportValidFields(content []string) int {
	var y, numValidPassportValidFields int

	for i, singleLine := range content {
		y++
		// Reset
		var numValidFields int
		if strings.Contains(singleLine, "byr:") {
			numValidFields++
		}
		if strings.Contains(singleLine, "iyr:") {
			numValidFields++
		}
		if strings.Contains(singleLine, "eyr:") {
			numValidFields++
		}
		if strings.Contains(singleLine, "hgt:") {
			numValidFields++
		}
		// Only hcl with # allowed!
		if strings.Contains(singleLine, "hcl:") {
			numValidFields++
		}
		if strings.Contains(singleLine, "ecl:") {
			numValidFields++
		}
		if strings.Contains(singleLine, "pid:") {
			numValidFields++
		}
		// Condition check
		if numValidFields < 7 {
			fmt.Println("INPUT #", i, " ", singleLine, " INVALID! FIELDS: ", numValidFields)
			continue
		}
		// Got here is valid!
		numValidPassportValidFields++
	}
	fmt.Println("TOTAL: ", y)
	return numValidPassportValidFields
}

func countValidPassport(content []string) int {
	// Split by word
	// For each, split by ':' check against valid values
	// Invalid, ignore; mark it out ..
	var numValidPassports int
	// If required field >=7; up the counter
	/*
		byr string
		iyr string
		eyr string
		hgt string
		hcl string
		ecl string
		pid string

	*/
	for i, singleLine := range content {
		// Reset
		var numValidFields int

		// Year (4 digits)
		reBYear, reBYerr := regexp.Compile(`byr:(\d{4})`)
		if reBYerr != nil {
			panic(reBYerr)
		}
		reIYear, reIYerr := regexp.Compile(`iyr:(\d{4})`)
		if reIYerr != nil {
			panic(reIYerr)
		}
		reEYear, reEYerr := regexp.Compile(`eyr:(\d{4})`)
		if reEYerr != nil {
			panic(reEYerr)
		}
		// PassportID (9 digits)
		// pid (Passport ID) - a nine-digit number, including leading zeroes.
		// Match those at edge [99\n] and at start/middle [99 ]
		rePid, rePerr := regexp.Compile(`pid:(\d{9}\s+|\d{9}$)`)
		if rePerr != nil {
			panic(rePerr)
		}
		// Eye Color (3 alpha code)
		// ecl (Eye Color) - exactly one of: amb blu brn gry grn hzl oth
		reEcl, reEerr := regexp.Compile(`ecl:(amb|blu|brn|gry|grn|hzl|oth)`)
		if reEerr != nil {
			panic(reEerr)
		}
		// Height: digit[cm,in]
		// hgt (Height) - a number followed by either cm or in:
		reHgt, reHGerr := regexp.Compile(`hgt:(\d+)(cm|in)`)
		if reHGerr != nil {
			panic(reHGerr)
		}
		// Hair color CSS 6 Hexcode
		// hcl (Hair Color) - a # followed by exactly six characters 0-9 or a-f.
		reHcl, reHerr := regexp.Compile(`hcl:#([a-f0-9]){6}`)
		if reHerr != nil {
			panic(reHerr)
		}

		// For years sanity check
		var byr, iyr, eyr int

		if strings.Contains(singleLine, "byr:") {
			// Reject those with #
			//if strings.Contains(singleLine, "byr:#") {
			//	fmt.Println("REJECTED #", i, " ", singleLine, "due to #!")
			//	continue
			//}
			if !reBYear.MatchString(singleLine) {
				fmt.Println("REJECTED #", i, " ", singleLine, "due to NOT YEAR!")
				continue
			}
			byr, _ = strconv.Atoi(reBYear.FindStringSubmatch(singleLine)[1])
			// byr (Birth Year) - four digits; at least 1920 and at most 2002.
			if !(byr >= 1920 && byr <= 2002) {
				fmt.Println("REJECTED #", i, " ", singleLine, "due to BYR least 1920 and at most 2002!")
				continue
			}
			//fmt.Sscanf(singleLine, "byr:%d", &byr)
			// DEBUG
			//fmt.Println("LINE: ", singleLine, " BYR: ", byr)
			numValidFields++
		}
		if strings.Contains(singleLine, "iyr:") {
			//if strings.Contains(singleLine, "iyr:#") {
			//	//fmt.Println("REJECTED #", i, " ", singleLine, "due to #!")
			//	continue
			//}
			if !reIYear.MatchString(singleLine) {
				fmt.Println("REJECTED #", i, " ", singleLine, "due to NOT YEAR!")
				continue
			}
			iyr, _ = strconv.Atoi(reIYear.FindStringSubmatch(singleLine)[1])
			// iyr (Issue Year) - four digits; at least 2010 and at most 2020.
			if !(iyr >= 2010 && iyr <= 2020) {
				fmt.Println("REJECTED #", i, " ", singleLine, "due to IYR least 2010 and at most 2020!")
				continue
			}
			//fmt.Sscanf(singleLine, "iyr:%d", &iyr)
			// DEBUG
			//fmt.Println("LINE: ", singleLine, " IYR: ", iyr)
			numValidFields++
		}
		if strings.Contains(singleLine, "eyr:") {
			//if strings.Contains(singleLine, "eyr:#") {
			//	//fmt.Println("REJECTED #", i, " ", singleLine, "due to #!")
			//	continue
			//}
			if !reEYear.MatchString(singleLine) {
				fmt.Println("REJECTED #", i, " ", singleLine, "due to NOT YEAR!")
				continue
			}
			eyr, _ = strconv.Atoi(reEYear.FindStringSubmatch(singleLine)[1])
			// eyr (Expiration Year) - four digits; at least 2020 and at most 2030.
			if !(eyr >= 2020 && eyr <= 2030) {
				fmt.Println("REJECTED #", i, " ", singleLine, "due to EYR least 2020 and at most 2030!")
				continue
			}
			//fmt.Sscanf(singleLine, "eyr:%d", &eyr)
			// DEBUG
			//fmt.Println("LINE: ", singleLine, " EYR: ", eyr)
			numValidFields++
		}
		// Should it have logic? byr < iyr < eyr .. probably .. nope, misunderstood :(
		//if byr == 0 || iyr == 0 || eyr == 0 {
		//	fmt.Println("REJECTED #", i, " ", singleLine, "due to NOT positive INT!")
		//	// DEBUG
		//	//fmt.Println(fmt.Sprintf("REJECTED ==> byr: %d iyr: %d eyr: %d", byr, iyr, eyr))
		//	continue
		//}

		if strings.Contains(singleLine, "hgt:") {
			if !reHgt.MatchString(singleLine) {
				fmt.Println("REJECTED #", i, " ", singleLine, "due to NOT HGT!")
				continue
			}
			// hgt (Height) - a number followed by either cm or in:
			hgt, _ := strconv.Atoi(reHgt.FindStringSubmatch(singleLine)[1])
			unit := reHgt.FindStringSubmatch(singleLine)[2]
			switch unit {
			case "cm":
				//If cm, the number must be at least 150 and at most 193.
				if !(hgt >= 150 && hgt <= 193) {
					fmt.Println("REJECTED #", i, " ", singleLine, "due to NOT the number must be at least 150 and at most 193 (cm)!")
					continue
				}
			case "in":
				//If in, the number must be at least 59 and at most 76.
				if !(hgt >= 59 && hgt <= 76) {
					fmt.Println("REJECTED #", i, " ", singleLine, "due to NOT the number must be at least 59 and at most 76 (in)!")
					continue
				}
			default:
				// Should NOT reach here!
				fmt.Println("REJECTED .. should NOT be here!")
				continue
			}

			numValidFields++
		}
		// Only hcl with # allowed!
		if strings.Contains(singleLine, "hcl:#") {
			if !reHcl.MatchString(singleLine) {
				fmt.Println("REJECTED #", i, " ", singleLine, "due to NOT HCL!")
				continue
			}
			numValidFields++
		}
		if strings.Contains(singleLine, "ecl:") {
			// Enum: grn, hzl; 3 word?
			// Assume: Open, take anything; can it be CSS?
			//if strings.Contains(singleLine, "ecl:#") {
			//	//fmt.Println("REJECTED #", i, " ", singleLine, "due to #!")
			//	continue
			//}
			if !reEcl.MatchString(singleLine) {
				fmt.Println("REJECTED #", i, " ", singleLine, "due to NOT Ecl!")
				continue
			}
			numValidFields++
		}
		if strings.Contains(singleLine, "pid:") {
			// Type: Int
			//if strings.Contains(singleLine, "pid:#") {
			//	//fmt.Println("REJECTED #", i, " ", singleLine, "due to #!")
			//	continue
			//}
			if !rePid.MatchString(singleLine) {
				fmt.Println("REJECTED #", i, " ", singleLine, "due to NOT PID!")
				continue
			}
			//pid, _ := strconv.Atoi(rePid.FindStringSubmatch(singleLine)[1])
			//// 9 digits
			//// 100000000 < x < 999999999
			//if !(pid >= 1 && pid <= 999999999) {
			//	fmt.Println("REJECTED #", i, " ", singleLine, "due to more than 9 digits!")
			//	continue
			//}
			//
			numValidFields++
		}
		// If got this far, is valid
		if numValidFields >= 7 {
			numValidPassports++
		} else {
			fmt.Println("don't have enough fields")
		}
	}

	return numValidPassports
}

// Process is the main call ..
func Process() error {

	full, err := os.Open("day4/testdata/full.txt")
	if err != nil {
		panic(err)
	}
	defer full.Close()

	//spew.Dump(splitPassports(full))
	//fmt.Println("Correct answer is 154 < x < 198 --> VALID PASSPORTS: ", countPassportValidFields(splitPassports(full)))
	//panic("err!")
	numValidPassports := countValidPassport(splitPassports(full))
	fmt.Println("Correct answer is x < 122 --> RESULTS: ", numValidPassports)
	return nil
}
