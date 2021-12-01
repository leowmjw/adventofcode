package day7

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

var BagHasShinyGold map[string]bool

// splitRules takes each line to extract BAGCOLOR and coded INT:BAGCOLOR,INT:BAGCOLOR
func splitRules(filename string) []string {
	var allRules []string

	fd, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(fd)
	fullLine := ""
	for scanner.Scan() {
		line := scanner.Text()
		// Do transformation; remove superflous
		line = strings.ReplaceAll(line, "contain", "")
		line = strings.ReplaceAll(line, "bags", "")
		line = strings.ReplaceAll(line, "bag", "")
		line = strings.ReplaceAll(line, ".", "")
		// Case: No other bags
		if strings.Contains(line, "no other") {
			line = strings.ReplaceAll(line, "no other", "")
			fullLine = strings.TrimSpace(line) + ";0:none"
			// DEBUG
			//fmt.Println("NO OTHER: ", fullLine)
			allRules = append(allRules, fullLine)
			continue
		}
		reRule := regexp.MustCompile(`^(.+?)(\d.*)$`)
		subMatches := reRule.FindStringSubmatch(line)
		// Split Color bag with ;
		//spew.Dump(subMatches)
		var allBags []string
		for _, bag := range strings.Split(subMatches[2], ",") {
			bag = strings.TrimSpace(bag)
			bag = strings.Replace(bag, " ", ":", 1)
			allBags = append(allBags, bag)
		}
		fullLine = strings.TrimSpace(subMatches[1]) + ";" + strings.Join(allBags, ",")
		// Create a new line to be appended
		allRules = append(allRules, fullLine)
	}
	// DEBUG
	//spew.Dump(allContent)

	// edge --> set(a,b)
	// replace edge with set(c,d)

	return allRules
}

// extractRulesMap boils it down to map[BAGCOLOR] = set(a b z) BAGCOLORS
func extractRulesMap(allRules []string) {
	// Init once
	BagHasShinyGold = make(map[string]bool)

	for _, singleRule := range allRules {
		// Extract BagColor
		rule := strings.Split(singleRule, ";")
		// Test if has shiny gold
		if strings.Contains(rule[1], "shiny gold") {
			BagHasShinyGold[rule[0]] = true
		}
	}
	spew.Dump(BagHasShinyGold)

	// If already existing, flag if is NOT new
	// so we know ..
	// map['shiny brown'] = "posh gold,bright magenta,pale bronze,light brown
}

func hasShinyGold(bags string) bool {
	// DEBUG
	fmt.Println("BAGS: ", bags)
	for _, bag := range strings.Split(bags, ",") {
		if BagHasShinyGold[strings.Split(bag, ":")[1]] {
			return true
		}
		// If it does not have
	}
	return false
}

// countBagContainsShinyGold takes allRulesMap and replace until only ShinyGold remains
func countBagContainsShinyGold(allRules []string) int {
	var countBagWithShinyGold int
	colorBagsWithShinyGold := make(map[string]int)
	// All rules, which one contains shiny gold? These key are mapped true
	// If contains has the mapped shiny; otherwise is recursive?
	extractRulesMap(allRules)
	// Needs to reevaluate the rules again and again ..
a:
	for _, singleRule := range allRules {
		// Extract Bags contained
		content := strings.Split(singleRule, ";")
		if colorBagsWithShinyGold[content[0]] == 0 {
			//spew.Dump(content[0])
			// If is the leave level; count and continue
			if strings.Contains(content[1], "shiny gold") {
				//fmt.Println(content[0])
				//fmt.Println("JACKPOT!! ", content[0], " contains: ", content[1])
				colorBagsWithShinyGold[content[0]]++
				countBagWithShinyGold++
				goto a
			}
			// otherwise, analyze its contents, need to find out, can break out
			//fmt.Println("LOOK IN BAG: ", content[0], " contains: ", content[1])
			for _, bag := range strings.Split(content[1], ",") {
				if BagHasShinyGold[strings.Split(bag, ":")[1]] {
					//fmt.Println(content[0])
					colorBagsWithShinyGold[content[0]]++
					//fmt.Println("JACKPOT!! ", content[0], " contains: ", content[1])
					countBagWithShinyGold++
					goto a
				}
			}
		}
	}
	spew.Dump(colorBagsWithShinyGold)
	return len(colorBagsWithShinyGold)
}

// Part1 covers at least one shiny gold bag
func Part1(filename string) int {
	allRules := splitRules(filename)
	countBagContainsShinyGold(allRules)

	return countBagContainsShinyGold(allRules)
}
