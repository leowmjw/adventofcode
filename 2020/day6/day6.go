package day6

import (
	"bufio"
	"os"
	"strings"
)

func splitCustomDeclarationForms(filename string) []string {
	var allContent []string

	fd, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(fd)
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
				fullLine = fullLine + line
			}
		}
	}
	// Don;t forget last line ..
	allContent = append(allContent, fullLine)
	// DEBUG
	//spew.Dump(allContent)
	return allContent
}

func countUniqueYesPerForm(content string) int {
	//var countUniqueYes int
	//countUniqueYes = len(strings.Split(content, ""))
	mapQuestion := make(map[string]bool, 100)
	for _, q := range strings.Split(content, "") {
		mapQuestion[q] = true
	}
	return len(mapQuestion)
}

// Part1 calculates Yes answers
func Part1(filename string) int {
	var countAnswerYes int
	customDeclarationForms := splitCustomDeclarationForms(filename)
	for _, singleCustomDeclarationForm := range customDeclarationForms {
		countAnswerYes += countUniqueYesPerForm(singleCustomDeclarationForm)
	}
	// DEBUG
	//fmt.Println("COUNT: ", countAnswerYes)
	return countAnswerYes
}

func splitCustomDeclarationFormsByIndividual(filename string) []string {
	var allContent []string

	fd, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(fd)
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
	//spew.Dump(allContent)
	return allContent
}

func countEveryoneYesPerForm(content string) int {
	var countEveryoneYes int
	individualsPerForm := len(strings.Split(content, " "))
	mapQuestion := make(map[string]int, 100)
	// For every Individual
	for _, individualAnswers := range strings.Split(content, " ") {
		// Increase count for particular questions
		for _, q := range strings.Split(individualAnswers, "") {
			mapQuestion[q]++
		}
	}

	// For every unique questions
	for _, countYes := range mapQuestion {
		// DEBUG
		//fmt.Println("Q: ", q, " COUNTY: ", countYes)
		// If equal number of people for the question, countEveryoneYes++
		if individualsPerForm == countYes {
			countEveryoneYes++
		}
	}
	return countEveryoneYes
}

// Part2 calculates everyone Yes answers
func Part2(filename string) int {
	var countEveryOneAnswerYes int
	customDeclarationForms := splitCustomDeclarationFormsByIndividual(filename)
	for _, singleCustomDeclarationForm := range customDeclarationForms {
		countEveryOneAnswerYes += countEveryoneYesPerForm(singleCustomDeclarationForm)
	}
	return countEveryOneAnswerYes
}
