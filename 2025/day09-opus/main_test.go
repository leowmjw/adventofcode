package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// Test data from the problem description
var exampleRedTiles = []Point{
	{7, 1},
	{11, 1},
	{11, 7},
	{9, 7},
	{9, 5},
	{2, 5},
	{2, 3},
	{7, 3},
}

// TestParseInput tests the input parsing function
func TestParseInput(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_input.txt")

	content := `7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3`

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	points := parseInput(testFile)

	if len(points) != 8 {
		t.Errorf("Expected 8 points, got %d", len(points))
	}

	if points[0].X != 7 || points[0].Y != 1 {
		t.Errorf("Expected first point (7,1), got (%d,%d)", points[0].X, points[0].Y)
	}

	if points[7].X != 7 || points[7].Y != 3 {
		t.Errorf("Expected last point (7,3), got (%d,%d)", points[7].X, points[7].Y)
	}
}

// TestParseInputWithWhitespace tests parsing with extra whitespace
func TestParseInputWithWhitespace(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_input.txt")

	content := `  7, 1  
11 , 1

9,5
`

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	points := parseInput(testFile)

	if len(points) != 3 {
		t.Errorf("Expected 3 points, got %d", len(points))
	}
}

// TestBuildPolygon tests polygon construction
func TestBuildPolygon(t *testing.T) {
	polygon := buildPolygon(exampleRedTiles)

	if polygon == nil {
		t.Fatal("Expected non-nil polygon")
	}

	if len(polygon.Vertices) != 8 {
		t.Errorf("Expected 8 vertices, got %d", len(polygon.Vertices))
	}

	// Should have slabs between unique Y coordinates
	// Y values: 1, 3, 5, 7 -> 3 slabs
	if len(polygon.Slabs) != 3 {
		t.Errorf("Expected 3 slabs, got %d", len(polygon.Slabs))
	}

	// Check Y coordinates are sorted
	for i := 1; i < len(polygon.YCoords); i++ {
		if polygon.YCoords[i] <= polygon.YCoords[i-1] {
			t.Errorf("Y coordinates not sorted: %v", polygon.YCoords)
			break
		}
	}
}

// TestBuildPolygonEmpty tests building polygon with no vertices
func TestBuildPolygonEmpty(t *testing.T) {
	polygon := buildPolygon([]Point{})

	if polygon == nil {
		t.Fatal("Expected non-nil polygon even for empty input")
	}

	if len(polygon.Vertices) != 0 {
		t.Errorf("Expected 0 vertices, got %d", len(polygon.Vertices))
	}
}

// TestIsRectangleValid tests rectangle validation
func TestIsRectangleValid(t *testing.T) {
	polygon := buildPolygon(exampleRedTiles)

	tests := []struct {
		name     string
		x1, y1   int
		x2, y2   int
		expected bool
	}{
		{
			name:     "Valid rectangle (2,3) to (9,5) - the 24 area answer",
			x1:       2, y1: 3,
			x2:       9, y2: 5,
			expected: true,
		},
		{
			name:     "Valid rectangle (9,5) to (11,7)",
			x1:       9, y1: 5,
			x2:       11, y2: 7,
			expected: true,
		},
		{
			name:     "Invalid rectangle outside polygon - far left",
			x1:       0, y1: 0,
			x2:       1, y2: 2,
			expected: false,
		},
		{
			name:     "Invalid rectangle crossing outside boundary",
			x1:       1, y1: 1,
			x2:       5, y2: 5,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := polygon.isRectangleValid(tt.x1, tt.y1, tt.x2, tt.y2)
			if result != tt.expected {
				t.Errorf("isRectangleValid(%d,%d,%d,%d) = %v, expected %v",
					tt.x1, tt.y1, tt.x2, tt.y2, result, tt.expected)
			}
		})
	}
}

// TestIsRectangleValidNormalization tests that coordinate order doesn't matter
func TestIsRectangleValidNormalization(t *testing.T) {
	polygon := buildPolygon(exampleRedTiles)

	// These should all give the same result regardless of corner order
	coords := []struct{ x1, y1, x2, y2 int }{
		{2, 3, 9, 5},
		{9, 3, 2, 5},
		{2, 5, 9, 3},
		{9, 5, 2, 3},
	}

	expected := polygon.isRectangleValid(2, 3, 9, 5)

	for _, c := range coords {
		result := polygon.isRectangleValid(c.x1, c.y1, c.x2, c.y2)
		if result != expected {
			t.Errorf("isRectangleValid(%d,%d,%d,%d) = %v, expected %v (normalization issue)",
				c.x1, c.y1, c.x2, c.y2, result, expected)
		}
	}
}

// TestAreaCalculation tests that area includes boundary tiles
func TestAreaCalculation(t *testing.T) {
	// Rectangle from (2,3) to (9,5):
	// Width = 9-2+1 = 8
	// Height = 5-3+1 = 3
	// Area = 24
	x1, y1 := 2, 3
	x2, y2 := 9, 5
	
	width := abs(x2-x1) + 1
	height := abs(y2-y1) + 1
	area := width * height
	
	if area != 24 {
		t.Errorf("Expected area 24 for (2,3)-(9,5), got %d (width=%d, height=%d)", area, width, height)
	}
}

// TestFindLargestRectangleConcurrent tests the main algorithm with example data
func TestFindLargestRectangleConcurrent(t *testing.T) {
	polygon := buildPolygon(exampleRedTiles)

	// Test with different worker counts
	workerCounts := []int{1, 2, 4, 8, 16}

	for _, workers := range workerCounts {
		t.Run(fmt.Sprintf("workers=%d", workers), func(t *testing.T) {
			result := findLargestRectangleConcurrent(exampleRedTiles, polygon, workers)

			// From the problem: largest rectangle using only red and green tiles has area 24
			expected := 24

			if result != expected {
				t.Errorf("Expected largest rectangle area %d, got %d", expected, result)
			}
		})
	}
}

// TestFindLargestRectangleSimpleSquare tests with a simple square polygon
func TestFindLargestRectangleSimpleSquare(t *testing.T) {
	// Simple square polygon
	points := []Point{
		{0, 0},
		{10, 0},
		{10, 10},
		{0, 10},
	}

	polygon := buildPolygon(points)
	result := findLargestRectangleConcurrent(points, polygon, 4)

	// Largest rectangle between opposite corners: (0,0) to (10,10)
	// Width = 11, Height = 11, Area = 121
	expected := 121

	if result != expected {
		t.Errorf("Expected area %d, got %d", expected, result)
	}
}

// TestFindLargestRectangleSinglePoint tests edge case
func TestFindLargestRectangleSinglePoint(t *testing.T) {
	points := []Point{{5, 5}}
	polygon := buildPolygon(points)
	result := findLargestRectangleConcurrent(points, polygon, 4)

	if result != 0 {
		t.Errorf("Expected 0 for single point, got %d", result)
	}
}

// TestFindLargestRectangleEmpty tests edge case
func TestFindLargestRectangleEmpty(t *testing.T) {
	points := []Point{}
	polygon := buildPolygon(points)
	result := findLargestRectangleConcurrent(points, polygon, 4)

	if result != 0 {
		t.Errorf("Expected 0 for empty input, got %d", result)
	}
}

// TestFindLargestRectangleCollinearPoints tests points on same line
func TestFindLargestRectangleCollinearPoints(t *testing.T) {
	// All points on same X - no valid rectangle possible
	points := []Point{
		{5, 0},
		{5, 5},
		{5, 10},
	}
	polygon := buildPolygon(points)
	result := findLargestRectangleConcurrent(points, polygon, 4)

	if result != 0 {
		t.Errorf("Expected 0 for collinear points, got %d", result)
	}
}

// TestAbs tests the absolute value function
func TestAbs(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
		{-1000000, 1000000},
		{1000000, 1000000},
	}

	for _, tt := range tests {
		result := abs(tt.input)
		if result != tt.expected {
			t.Errorf("abs(%d) = %d, expected %d", tt.input, result, tt.expected)
		}
	}
}

// TestSlabStructure tests the internal slab structure
func TestSlabStructure(t *testing.T) {
	polygon := buildPolygon(exampleRedTiles)

	if len(polygon.Slabs) == 0 {
		t.Fatal("No slabs created")
	}

	// First slab should start at minimum Y
	minY := polygon.YCoords[0]
	if polygon.Slabs[0].YMin != minY {
		t.Errorf("First slab YMin = %d, expected %d", polygon.Slabs[0].YMin, minY)
	}

	// Last slab should end at maximum Y
	maxY := polygon.YCoords[len(polygon.YCoords)-1]
	lastSlab := polygon.Slabs[len(polygon.Slabs)-1]
	if lastSlab.YMax != maxY {
		t.Errorf("Last slab YMax = %d, expected %d", lastSlab.YMax, maxY)
	}

	// Slabs should be contiguous
	for i := 1; i < len(polygon.Slabs); i++ {
		if polygon.Slabs[i].YMin != polygon.Slabs[i-1].YMax {
			t.Errorf("Gap between slabs %d and %d", i-1, i)
		}
	}
}

// TestConcurrencyCorrectness runs the algorithm multiple times to check for race conditions
func TestConcurrencyCorrectness(t *testing.T) {
	polygon := buildPolygon(exampleRedTiles)

	var results []int
	for i := 0; i < 10; i++ {
		result := findLargestRectangleConcurrent(exampleRedTiles, polygon, 16)
		results = append(results, result)
	}

	for i := 1; i < len(results); i++ {
		if results[i] != results[0] {
			t.Errorf("Inconsistent results: got %d and %d (possible race condition)",
				results[0], results[i])
		}
	}
}

// TestRealInput tests with the actual puzzle input
func TestRealInput(t *testing.T) {
	if _, err := os.Stat("/tmp/aoc-day9/input.txt"); os.IsNotExist(err) {
		t.Skip("Skipping real input test - input.txt not found")
	}

	redTiles := parseInput("/tmp/aoc-day9/input.txt")
	polygon := buildPolygon(redTiles)
	result := findLargestRectangleConcurrent(redTiles, polygon, 16)

	// We'll update this once we know the correct answer
	t.Logf("Real input result: %d", result)
	if result <= 0 {
		t.Errorf("Expected positive result for real input, got %d", result)
	}
}

// Benchmark tests
func BenchmarkBuildPolygon(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buildPolygon(exampleRedTiles)
	}
}

func BenchmarkIsRectangleValid(b *testing.B) {
	polygon := buildPolygon(exampleRedTiles)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		polygon.isRectangleValid(2, 3, 9, 5)
	}
}

func BenchmarkFindLargestRectangle1Worker(b *testing.B) {
	polygon := buildPolygon(exampleRedTiles)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		findLargestRectangleConcurrent(exampleRedTiles, polygon, 1)
	}
}

func BenchmarkFindLargestRectangle16Workers(b *testing.B) {
	polygon := buildPolygon(exampleRedTiles)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		findLargestRectangleConcurrent(exampleRedTiles, polygon, 16)
	}
}
