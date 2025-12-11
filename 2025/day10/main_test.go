package main

import (
	"os"
	"reflect"
	"testing"
)

func TestParseMachine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Machine
		wantErr  bool
	}{
		{
			name:  "simple pattern with multiple buttons",
			input: "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}",
			expected: Machine{
				Target: []bool{false, true, true, false},
				Buttons: [][]int{
					{3},
					{1, 3},
					{2},
					{2, 3},
					{0, 2},
					{0, 1},
				},
			},
			wantErr: false,
		},
		{
			name:  "pattern with dots only",
			input: "[....] (0,1) {1}",
			expected: Machine{
				Target:  []bool{false, false, false, false},
				Buttons: [][]int{{0, 1}},
			},
			wantErr: false,
		},
		{
			name:  "pattern with hashes only",
			input: "[####] (0,1,2,3) {1}",
			expected: Machine{
				Target:  []bool{true, true, true, true},
				Buttons: [][]int{{0, 1, 2, 3}},
			},
			wantErr: false,
		},
		{
			name:  "five light pattern",
			input: "[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}",
			expected: Machine{
				Target: []bool{false, false, false, true, false},
				Buttons: [][]int{
					{0, 2, 3, 4},
					{2, 3},
					{0, 4},
					{0, 1, 2},
					{1, 2, 3, 4},
				},
			},
			wantErr: false,
		},
		{
			name:  "six light pattern",
			input: "[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}",
			expected: Machine{
				Target: []bool{false, true, true, true, false, true},
				Buttons: [][]int{
					{0, 1, 2, 3, 4},
					{0, 3, 4},
					{0, 1, 2, 4, 5},
					{1, 2},
				},
			},
			wantErr: false,
		},
		{
			name:    "missing brackets",
			input:   ".##. (3) (1,3) {3,5,4,7}",
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMachine(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMachine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got.Target, tt.expected.Target) {
					t.Errorf("ParseMachine() Target = %v, want %v", got.Target, tt.expected.Target)
				}
				if !reflect.DeepEqual(got.Buttons, tt.expected.Buttons) {
					t.Errorf("ParseMachine() Buttons = %v, want %v", got.Buttons, tt.expected.Buttons)
				}
			}
		})
	}
}

func TestSimulatePresses(t *testing.T) {
	tests := []struct {
		name      string
		numLights int
		buttons   [][]int
		presses   []bool
		expected  []bool
	}{
		{
			name:      "no presses",
			numLights: 4,
			buttons:   [][]int{{0, 1}, {2, 3}},
			presses:   []bool{false, false},
			expected:  []bool{false, false, false, false},
		},
		{
			name:      "single button press",
			numLights: 4,
			buttons:   [][]int{{0, 1}, {2, 3}},
			presses:   []bool{true, false},
			expected:  []bool{true, true, false, false},
		},
		{
			name:      "both buttons pressed",
			numLights: 4,
			buttons:   [][]int{{0, 1}, {2, 3}},
			presses:   []bool{true, true},
			expected:  []bool{true, true, true, true},
		},
		{
			name:      "overlapping toggles",
			numLights: 4,
			buttons:   [][]int{{0, 1}, {1, 2}},
			presses:   []bool{true, true},
			expected:  []bool{true, false, true, false}, // light 1 toggled twice = off
		},
		{
			name:      "toggle same light multiple times",
			numLights: 3,
			buttons:   [][]int{{0}, {0}, {0}},
			presses:   []bool{true, true, true},
			expected:  []bool{true, false, false}, // toggled 3 times = on
		},
		{
			name:      "example from problem - [#.....] -> [...##.]",
			numLights: 6,
			buttons:   [][]int{{0, 3, 4}},
			presses:   []bool{true},
			expected:  []bool{true, false, false, true, true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SimulatePresses(tt.numLights, tt.buttons, tt.presses)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("SimulatePresses() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMatchesTarget(t *testing.T) {
	tests := []struct {
		name     string
		lights   []bool
		target   []bool
		expected bool
	}{
		{
			name:     "exact match",
			lights:   []bool{true, false, true},
			target:   []bool{true, false, true},
			expected: true,
		},
		{
			name:     "no match",
			lights:   []bool{true, false, true},
			target:   []bool{false, true, false},
			expected: false,
		},
		{
			name:     "partial match",
			lights:   []bool{true, false, true},
			target:   []bool{true, false, false},
			expected: false,
		},
		{
			name:     "different lengths",
			lights:   []bool{true, false},
			target:   []bool{true, false, true},
			expected: false,
		},
		{
			name:     "empty arrays",
			lights:   []bool{},
			target:   []bool{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchesTarget(tt.lights, tt.target)
			if got != tt.expected {
				t.Errorf("MatchesTarget() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCountPresses(t *testing.T) {
	tests := []struct {
		name     string
		presses  []bool
		expected int
	}{
		{
			name:     "no presses",
			presses:  []bool{false, false, false},
			expected: 0,
		},
		{
			name:     "all pressed",
			presses:  []bool{true, true, true},
			expected: 3,
		},
		{
			name:     "some pressed",
			presses:  []bool{true, false, true, false, true},
			expected: 3,
		},
		{
			name:     "empty",
			presses:  []bool{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountPresses(tt.presses)
			if got != tt.expected {
				t.Errorf("CountPresses() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFindMinPresses(t *testing.T) {
	tests := []struct {
		name     string
		machine  Machine
		expected int
	}{
		{
			name: "example 1 - [.##.] with 6 buttons, min 2",
			machine: Machine{
				Target: []bool{false, true, true, false},
				Buttons: [][]int{
					{3},
					{1, 3},
					{2},
					{2, 3},
					{0, 2},
					{0, 1},
				},
			},
			expected: 2,
		},
		{
			name: "example 2 - [...#.] with 5 buttons, min 3",
			machine: Machine{
				Target: []bool{false, false, false, true, false},
				Buttons: [][]int{
					{0, 2, 3, 4},
					{2, 3},
					{0, 4},
					{0, 1, 2},
					{1, 2, 3, 4},
				},
			},
			expected: 3,
		},
		{
			name: "example 3 - [.###.#] with 4 buttons, min 2",
			machine: Machine{
				Target: []bool{false, true, true, true, false, true},
				Buttons: [][]int{
					{0, 1, 2, 3, 4},
					{0, 3, 4},
					{0, 1, 2, 4, 5},
					{1, 2},
				},
			},
			expected: 2,
		},
		{
			name: "all lights off target - no presses needed",
			machine: Machine{
				Target:  []bool{false, false, false},
				Buttons: [][]int{{0, 1}, {1, 2}},
			},
			expected: 0,
		},
		{
			name: "single button toggles exactly target lights",
			machine: Machine{
				Target:  []bool{true, true, false},
				Buttons: [][]int{{0, 1}},
			},
			expected: 1,
		},
		{
			name: "no buttons available but target is all off",
			machine: Machine{
				Target:  []bool{false, false},
				Buttons: [][]int{},
			},
			expected: 0,
		},
		{
			name: "no buttons available but target requires on - impossible",
			machine: Machine{
				Target:  []bool{true, false},
				Buttons: [][]int{},
			},
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindMinPresses(tt.machine)
			if got != tt.expected {
				t.Errorf("FindMinPresses() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFindMinPresses_ExampleSum(t *testing.T) {
	// From the problem: 2 + 3 + 2 = 7
	machines := []Machine{
		{
			Target: []bool{false, true, true, false},
			Buttons: [][]int{
				{3},
				{1, 3},
				{2},
				{2, 3},
				{0, 2},
				{0, 1},
			},
		},
		{
			Target: []bool{false, false, false, true, false},
			Buttons: [][]int{
				{0, 2, 3, 4},
				{2, 3},
				{0, 4},
				{0, 1, 2},
				{1, 2, 3, 4},
			},
		},
		{
			Target: []bool{false, true, true, true, false, true},
			Buttons: [][]int{
				{0, 1, 2, 3, 4},
				{0, 3, 4},
				{0, 1, 2, 4, 5},
				{1, 2},
			},
		},
	}

	total := 0
	for _, m := range machines {
		total += FindMinPresses(m)
	}

	if total != 7 {
		t.Errorf("Total minimum presses = %v, want 7", total)
	}
}

func TestSolvePart1_ExampleFile(t *testing.T) {
	// Create a temporary test file with the example input
	content := `[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}
`
	tmpfile, err := os.CreateTemp("", "test_input_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	result, err := SolvePart1(tmpfile.Name())
	if err != nil {
		t.Fatalf("SolvePart1() error = %v", err)
	}

	if result != 7 {
		t.Errorf("SolvePart1() = %v, want 7", result)
	}
}

func TestSolvePart1_EmptyLines(t *testing.T) {
	content := `[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}

[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
`
	tmpfile, err := os.CreateTemp("", "test_input_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	result, err := SolvePart1(tmpfile.Name())
	if err != nil {
		t.Fatalf("SolvePart1() error = %v", err)
	}

	// 2 + 3 = 5
	if result != 5 {
		t.Errorf("SolvePart1() = %v, want 5", result)
	}
}

func TestSolvePart1_FileNotFound(t *testing.T) {
	_, err := SolvePart1("nonexistent_file.txt")
	if err == nil {
		t.Error("SolvePart1() expected error for non-existent file")
	}
}

func TestParseMachine_ComplexInput(t *testing.T) {
	// Test parsing from actual input format
	input := "[.#...#...#] (3,5,7,8) (0,3,4) (0,1,2,3,4,7,9) (0,1,3,4,6,7,9) (1,4,5,6,8) (0,1,6,9) (0,2,3,4,5,7,8,9) (1,2,5,6,7,9) (0,2,3,5,6,7,8,9) (0,2,3,4,5,7,8) {46,36,54,60,41,78,47,75,59,57}"

	machine, err := ParseMachine(input)
	if err != nil {
		t.Fatalf("ParseMachine() error = %v", err)
	}

	expectedTarget := []bool{false, true, false, false, false, true, false, false, false, true}
	if !reflect.DeepEqual(machine.Target, expectedTarget) {
		t.Errorf("Target = %v, want %v", machine.Target, expectedTarget)
	}

	if len(machine.Buttons) != 10 {
		t.Errorf("Number of buttons = %v, want 10", len(machine.Buttons))
	}

	// Check first button
	expectedFirstButton := []int{3, 5, 7, 8}
	if !reflect.DeepEqual(machine.Buttons[0], expectedFirstButton) {
		t.Errorf("First button = %v, want %v", machine.Buttons[0], expectedFirstButton)
	}
}

func TestSolvePart1_ActualInput(t *testing.T) {
	// Test with actual input file if it exists
	if _, err := os.Stat("input.txt"); os.IsNotExist(err) {
		t.Skip("input.txt not found")
	}

	result, err := SolvePart1("input.txt")
	if err != nil {
		t.Fatalf("SolvePart1() error = %v", err)
	}

	// Verify result is a positive number
	if result <= 0 {
		t.Errorf("SolvePart1() = %v, expected positive result", result)
	}

	t.Logf("Part 1 result: %d", result)
}

func TestSimulatePresses_EmptyButtons(t *testing.T) {
	result := SimulatePresses(3, [][]int{}, []bool{})
	expected := []bool{false, false, false}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SimulatePresses() = %v, want %v", result, expected)
	}
}

func TestSimulatePresses_ButtonOutOfRange(t *testing.T) {
	// Button toggles index 5 but only 3 lights exist - should be ignored
	result := SimulatePresses(3, [][]int{{0, 5}}, []bool{true})
	expected := []bool{true, false, false}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SimulatePresses() = %v, want %v", result, expected)
	}
}

func TestParseMachineJoltages(t *testing.T) {
	input := "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}"
	machine, err := ParseMachine(input)
	if err != nil {
		t.Fatalf("ParseMachine() error = %v", err)
	}

	expectedJoltages := []int{3, 5, 4, 7}
	if !reflect.DeepEqual(machine.Joltages, expectedJoltages) {
		t.Errorf("Joltages = %v, want %v", machine.Joltages, expectedJoltages)
	}
}

func TestFindMinPressesPart2(t *testing.T) {
	tests := []struct {
		name     string
		machine  Machine
		expected int
	}{
		{
			name: "example 1 - {3,5,4,7} min 10",
			machine: Machine{
				Buttons: [][]int{
					{3},
					{1, 3},
					{2},
					{2, 3},
					{0, 2},
					{0, 1},
				},
				Joltages: []int{3, 5, 4, 7},
			},
			expected: 10,
		},
		{
			name: "example 2 - {7,5,12,7,2} min 12",
			machine: Machine{
				Buttons: [][]int{
					{0, 2, 3, 4},
					{2, 3},
					{0, 4},
					{0, 1, 2},
					{1, 2, 3, 4},
				},
				Joltages: []int{7, 5, 12, 7, 2},
			},
			expected: 12,
		},
		{
			name: "example 3 - {10,11,11,5,10,5} min 11",
			machine: Machine{
				Buttons: [][]int{
					{0, 1, 2, 3, 4},
					{0, 3, 4},
					{0, 1, 2, 4, 5},
					{1, 2},
				},
				Joltages: []int{10, 11, 11, 5, 10, 5},
			},
			expected: 11,
		},
		{
			name: "all zeros - no presses needed",
			machine: Machine{
				Buttons:  [][]int{{0, 1}, {1, 2}},
				Joltages: []int{0, 0, 0},
			},
			expected: 0,
		},
		{
			name: "single button single counter",
			machine: Machine{
				Buttons:  [][]int{{0}},
				Joltages: []int{5},
			},
			expected: 5,
		},
		{
			name: "no buttons all zero targets",
			machine: Machine{
				Buttons:  [][]int{},
				Joltages: []int{0, 0},
			},
			expected: 0,
		},
		{
			name: "no buttons nonzero target - impossible",
			machine: Machine{
				Buttons:  [][]int{},
				Joltages: []int{1, 0},
			},
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindMinPressesPart2(tt.machine)
			if got != tt.expected {
				t.Errorf("FindMinPressesPart2() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFindMinPressesPart2_ExampleSum(t *testing.T) {
	// From the problem: 10 + 12 + 11 = 33
	machines := []Machine{
		{
			Buttons: [][]int{
				{3},
				{1, 3},
				{2},
				{2, 3},
				{0, 2},
				{0, 1},
			},
			Joltages: []int{3, 5, 4, 7},
		},
		{
			Buttons: [][]int{
				{0, 2, 3, 4},
				{2, 3},
				{0, 4},
				{0, 1, 2},
				{1, 2, 3, 4},
			},
			Joltages: []int{7, 5, 12, 7, 2},
		},
		{
			Buttons: [][]int{
				{0, 1, 2, 3, 4},
				{0, 3, 4},
				{0, 1, 2, 4, 5},
				{1, 2},
			},
			Joltages: []int{10, 11, 11, 5, 10, 5},
		},
	}

	total := 0
	for _, m := range machines {
		total += FindMinPressesPart2(m)
	}

	if total != 33 {
		t.Errorf("Total minimum presses = %v, want 33", total)
	}
}

func TestSolvePart2_ExampleFile(t *testing.T) {
	content := `[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}
`
	tmpfile, err := os.CreateTemp("", "test_input_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	result, err := SolvePart2(tmpfile.Name())
	if err != nil {
		t.Fatalf("SolvePart2() error = %v", err)
	}

	if result != 33 {
		t.Errorf("SolvePart2() = %v, want 33", result)
	}
}

// Benchmark test to ensure performance is acceptable
func BenchmarkFindMinPresses(b *testing.B) {
	machine := Machine{
		Target: []bool{false, true, false, false, false, true, false, false, false, true},
		Buttons: [][]int{
			{3, 5, 7, 8},
			{0, 3, 4},
			{0, 1, 2, 3, 4, 7, 9},
			{0, 1, 3, 4, 6, 7, 9},
			{1, 4, 5, 6, 8},
			{0, 1, 6, 9},
			{0, 2, 3, 4, 5, 7, 8, 9},
			{1, 2, 5, 6, 7, 9},
			{0, 2, 3, 5, 6, 7, 8, 9},
			{0, 2, 3, 4, 5, 7, 8},
		},
	}

	for i := 0; i < b.N; i++ {
		FindMinPresses(machine)
	}
}
