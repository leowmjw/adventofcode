package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Machine represents a factory machine with indicator lights and buttons
type Machine struct {
	Target   []bool  // desired light states (true = on, false = off)
	Buttons  [][]int // each button is a list of light indices it toggles
	Joltages []int   // joltage requirements for Part 2
}

// ParseMachine parses a single line into a Machine struct
func ParseMachine(line string) (Machine, error) {
	var m Machine

	// Extract target pattern from [...]
	startBracket := strings.Index(line, "[")
	endBracket := strings.Index(line, "]")
	if startBracket == -1 || endBracket == -1 || endBracket <= startBracket {
		return m, fmt.Errorf("invalid target pattern in line: %s", line)
	}

	pattern := line[startBracket+1 : endBracket]
	m.Target = make([]bool, len(pattern))
	for i, ch := range pattern {
		m.Target[i] = ch == '#'
	}

	// Extract button configurations from (...)
	remaining := line[endBracket+1:]
	for {
		startParen := strings.Index(remaining, "(")
		if startParen == -1 {
			break
		}
		endParen := strings.Index(remaining, ")")
		if endParen == -1 || endParen <= startParen {
			return m, fmt.Errorf("invalid button configuration in line: %s", line)
		}

		buttonStr := remaining[startParen+1 : endParen]
		// Check if this is the joltage section (starts with {)
		if strings.Index(remaining, "{") != -1 && strings.Index(remaining, "{") < startParen {
			break
		}

		var button []int
		if buttonStr != "" {
			parts := strings.Split(buttonStr, ",")
			for _, p := range parts {
				idx, err := strconv.Atoi(strings.TrimSpace(p))
				if err != nil {
					return m, fmt.Errorf("invalid button index %q: %v", p, err)
				}
				button = append(button, idx)
			}
		}
		m.Buttons = append(m.Buttons, button)
		remaining = remaining[endParen+1:]
	}

	// Extract joltage requirements from {...}
	startBrace := strings.Index(line, "{")
	endBrace := strings.Index(line, "}")
	if startBrace != -1 && endBrace != -1 && endBrace > startBrace {
		joltageStr := line[startBrace+1 : endBrace]
		parts := strings.Split(joltageStr, ",")
		for _, p := range parts {
			val, err := strconv.Atoi(strings.TrimSpace(p))
			if err != nil {
				return m, fmt.Errorf("invalid joltage value %q: %v", p, err)
			}
			m.Joltages = append(m.Joltages, val)
		}
	}

	return m, nil
}

// SimulatePresses simulates pressing buttons and returns the resulting light states
func SimulatePresses(numLights int, buttons [][]int, presses []bool) []bool {
	lights := make([]bool, numLights)
	for i, pressed := range presses {
		if pressed {
			for _, lightIdx := range buttons[i] {
				if lightIdx < numLights {
					lights[lightIdx] = !lights[lightIdx]
				}
			}
		}
	}
	return lights
}

// MatchesTarget checks if the light state matches the target
func MatchesTarget(lights, target []bool) bool {
	if len(lights) != len(target) {
		return false
	}
	for i := range lights {
		if lights[i] != target[i] {
			return false
		}
	}
	return true
}

// CountPresses counts the number of true values in the presses slice
func CountPresses(presses []bool) int {
	count := 0
	for _, p := range presses {
		if p {
			count++
		}
	}
	return count
}

// FindMinPresses finds the minimum number of button presses to achieve the target
// Uses brute force enumeration over all 2^n combinations
func FindMinPresses(m Machine) int {
	numButtons := len(m.Buttons)
	numLights := len(m.Target)

	if numButtons == 0 {
		// No buttons - check if target is already achieved (all off)
		allOff := true
		for _, t := range m.Target {
			if t {
				allOff = false
				break
			}
		}
		if allOff {
			return 0
		}
		return -1 // Impossible
	}

	minPresses := numButtons + 1 // Start with impossible value
	totalCombinations := 1 << numButtons

	for combo := 0; combo < totalCombinations; combo++ {
		// Convert combo to button press pattern
		presses := make([]bool, numButtons)
		for i := 0; i < numButtons; i++ {
			presses[i] = (combo & (1 << i)) != 0
		}

		// Simulate and check
		result := SimulatePresses(numLights, m.Buttons, presses)
		if MatchesTarget(result, m.Target) {
			pressCount := CountPresses(presses)
			if pressCount < minPresses {
				minPresses = pressCount
			}
		}
	}

	if minPresses > numButtons {
		return -1 // No solution found
	}
	return minPresses
}

// FindMinPressesPart2 finds the minimum number of button presses to achieve joltage targets
// Each button can be pressed multiple times, and each press increments the affected counters by 1
// This is an integer linear programming problem: minimize sum(x_i) subject to A*x = b, x >= 0
func FindMinPressesPart2(m Machine) int {
	numButtons := len(m.Buttons)
	numCounters := len(m.Joltages)

	if numButtons == 0 {
		// No buttons - check if all joltages are already 0
		for _, j := range m.Joltages {
			if j != 0 {
				return -1 // Impossible
			}
		}
		return 0
	}

	// Build the augmented matrix [A | b] for Gaussian elimination
	// A[counter][button] = 1 if button affects counter, else 0
	// Working with rationals represented as fractions to avoid floating point issues
	// But since we only have 0s and 1s, we can use integers with careful division

	// Use float64 for simplicity in Gaussian elimination
	matrix := make([][]float64, numCounters)
	for i := range matrix {
		matrix[i] = make([]float64, numButtons+1) // +1 for the target column
		matrix[i][numButtons] = float64(m.Joltages[i])
	}
	for btnIdx, btn := range m.Buttons {
		for _, counterIdx := range btn {
			if counterIdx < numCounters {
				matrix[counterIdx][btnIdx] = 1
			}
		}
	}

	// Gaussian elimination with partial pivoting
	pivotCol := make([]int, numCounters) // Which column is the pivot for each row
	for i := range pivotCol {
		pivotCol[i] = -1
	}

	row := 0
	for col := 0; col < numButtons && row < numCounters; col++ {
		// Find pivot
		maxRow := row
		for i := row + 1; i < numCounters; i++ {
			if abs(matrix[i][col]) > abs(matrix[maxRow][col]) {
				maxRow = i
			}
		}

		if abs(matrix[maxRow][col]) < 1e-9 {
			continue // No pivot in this column
		}

		// Swap rows
		matrix[row], matrix[maxRow] = matrix[maxRow], matrix[row]
		pivotCol[row] = col

		// Scale pivot row
		scale := matrix[row][col]
		for j := col; j <= numButtons; j++ {
			matrix[row][j] /= scale
		}

		// Eliminate column
		for i := 0; i < numCounters; i++ {
			if i != row && abs(matrix[i][col]) > 1e-9 {
				factor := matrix[i][col]
				for j := col; j <= numButtons; j++ {
					matrix[i][j] -= factor * matrix[row][j]
				}
			}
		}
		row++
	}

	numPivots := row

	// Check for inconsistency (0 = nonzero)
	for i := numPivots; i < numCounters; i++ {
		if abs(matrix[i][numButtons]) > 1e-9 {
			return -1 // No solution
		}
	}

	// Identify free variables (columns without pivots)
	isPivotCol := make([]bool, numButtons)
	for i := 0; i < numPivots; i++ {
		if pivotCol[i] >= 0 {
			isPivotCol[pivotCol[i]] = true
		}
	}

	freeVars := []int{}
	for col := 0; col < numButtons; col++ {
		if !isPivotCol[col] {
			freeVars = append(freeVars, col)
		}
	}

	// If no free variables, we have a unique solution
	if len(freeVars) == 0 {
		solution := make([]float64, numButtons)
		for i := numPivots - 1; i >= 0; i-- {
			col := pivotCol[i]
			solution[col] = matrix[i][numButtons]
			for j := col + 1; j < numButtons; j++ {
				solution[col] -= matrix[i][j] * solution[j]
			}
		}

		// Check if solution is non-negative integers
		total := 0
		for _, v := range solution {
			if v < -1e-9 || abs(v-float64(int(v+0.5))) > 1e-9 {
				return -1 // Not a valid non-negative integer solution
			}
			total += int(v + 0.5)
		}
		return total
	}

	// With free variables, search over non-negative integer values for them
	// and compute the rest, looking for minimum total
	return searchFreeVariables(matrix, pivotCol, numPivots, freeVars, numButtons, m.Joltages)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// searchFreeVariables searches over possible values of free variables to find minimum total presses
func searchFreeVariables(matrix [][]float64, pivotCol []int, numPivots int, freeVars []int, numButtons int, targets []int) int {
	// Upper bound for free variables
	maxVal := 0
	for _, t := range targets {
		if t > maxVal {
			maxVal = t
		}
	}

	minTotal := -1
	freeValues := make([]int, len(freeVars))

	// DFS over free variable values
	var search func(idx int, currentFreeSum int)
	search = func(idx int, currentFreeSum int) {
		if idx == len(freeVars) {
			// Compute pivot variable values
			solution := make([]float64, numButtons)
			for i, fv := range freeVars {
				solution[fv] = float64(freeValues[i])
			}

			valid := true
			total := currentFreeSum
			for i := numPivots - 1; i >= 0 && valid; i-- {
				col := pivotCol[i]
				val := matrix[i][numButtons]
				for j := col + 1; j < numButtons; j++ {
					val -= matrix[i][j] * solution[j]
				}
				solution[col] = val

				// Check non-negative integer
				rounded := int(val + 0.5)
				if val < -1e-9 || abs(val-float64(rounded)) > 1e-9 {
					valid = false
				} else {
					total += rounded
				}
			}

			if valid && (minTotal < 0 || total < minTotal) {
				minTotal = total
			}
			return
		}

		// Prune: if current free sum already exceeds best, stop
		if minTotal >= 0 && currentFreeSum >= minTotal {
			return
		}

		for v := 0; v <= maxVal; v++ {
			freeValues[idx] = v
			search(idx+1, currentFreeSum+v)
		}
	}

	search(0, 0)
	return minTotal
}

// SolvePart1 reads the input file and returns the sum of minimum presses
func SolvePart1(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	total := 0
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNum++
		if line == "" {
			continue
		}

		machine, err := ParseMachine(line)
		if err != nil {
			return 0, fmt.Errorf("line %d: %v", lineNum, err)
		}

		minPresses := FindMinPresses(machine)
		if minPresses < 0 {
			return 0, fmt.Errorf("line %d: no solution found", lineNum)
		}
		total += minPresses
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return total, nil
}

// SolvePart2 reads the input file and returns the sum of minimum presses for joltage targets
func SolvePart2(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	total := 0
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNum++
		if line == "" {
			continue
		}

		machine, err := ParseMachine(line)
		if err != nil {
			return 0, fmt.Errorf("line %d: %v", lineNum, err)
		}

		minPresses := FindMinPressesPart2(machine)
		if minPresses < 0 {
			return 0, fmt.Errorf("line %d: no solution found", lineNum)
		}
		total += minPresses
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return total, nil
}

func main() {
	result1, err := SolvePart1("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Part 1 Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Part 1:", result1)

	result2, err := SolvePart2("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Part 2 Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Part 2:", result2)
}
