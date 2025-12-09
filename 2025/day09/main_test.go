package main

import (
	"testing"
)

func TestFindLargestRectangle(t *testing.T) {
	tests := []struct {
		name     string
		points   []Point
		expected int
	}{
		{
			name: "Example from problem - area 50 (corners 2,5 and 11,1)",
			points: []Point{
				{7, 1}, {11, 1}, {11, 7}, {9, 7}, {9, 5},
				{2, 5}, {2, 3}, {7, 3},
			},
			expected: 50,
		},
		{
			name:     "No points",
			points:   []Point{},
			expected: 0,
		},
		{
			name: "Single point",
			points: []Point{
				{5, 5},
			},
			expected: 0,
		},
		{
			name: "Two points - vertical line (same x, area 11)",
			points: []Point{
				{5, 0}, {5, 10},
			},
			expected: 11,
		},
		{
			name: "Two points - horizontal line (same y, area 11)",
			points: []Point{
				{0, 5}, {10, 5},
			},
			expected: 11,
		},
		{
			name: "Two points - thin rectangle (area 6, 1 tall line)",
			points: []Point{
				{2, 3}, {7, 3},
			},
			expected: 6,
		},
		{
			name: "Two points - rectangle (area 24)",
			points: []Point{
				{2, 5}, {9, 7},
			},
			expected: 24,
		},
		{
			name: "Three points - max area 121",
			points: []Point{
				{0, 0}, {5, 5}, {10, 10},
			},
			expected: 121,
		},
		{
			name: "Four points - square (area 36)",
			points: []Point{
				{0, 0}, {5, 0}, {0, 5}, {5, 5},
			},
			expected: 36,
		},
		{
			name: "Negative coordinates (area 121)",
			points: []Point{
				{-5, -5}, {5, 5},
			},
			expected: 121,
		},
		{
			name: "Mixed positive and negative (area 231)",
			points: []Point{
				{-10, -5}, {10, 5},
			},
			expected: 231,
		},
		{
			name: "All points on same x-axis (max area 16, from 0 to 15)",
			points: []Point{
				{0, 0}, {5, 0}, {10, 0}, {15, 0},
			},
			expected: 16,
		},
		{
			name: "All points on same y-axis (max area 16, from 0 to 15)",
			points: []Point{
				{0, 0}, {0, 5}, {0, 10}, {0, 15},
			},
			expected: 16,
		},
		{
			name: "Large coordinate values (same y, area 675)",
			points: []Point{
				{97615, 50359}, {98289, 50359},
			},
			expected: 675,
		},
		{
			name: "Two points with large area (1002001)",
			points: []Point{
				{0, 0}, {1000, 1000},
			},
			expected: 1002001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findLargestRectangle(tt.points)
			if result != tt.expected {
				t.Errorf("findLargestRectangle(%v) = %d, expected %d", tt.name, result, tt.expected)
			}
		})
	}
}

func TestAbsFunction(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
		{-1000, 1000},
		{1000, 1000},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := abs(tt.input)
			if result != tt.expected {
				t.Errorf("abs(%d) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseInput(t *testing.T) {
	t.Run("Parse valid input file", func(t *testing.T) {
		points := parseInput("input.txt")
		if len(points) == 0 {
			t.Error("Expected points to be parsed from input file")
		}
		// Verify first point is parsed correctly
		if points[0].x == 0 && points[0].y == 0 {
			t.Error("Expected valid coordinates, got 0,0")
		}
	})
}

func TestRectangleAreaCalculation(t *testing.T) {
	tests := []struct {
		name    string
		p1      Point
		p2      Point
		expArea int
	}{
		{
			name:    "Simple rectangle 3x4 (inclusive)",
			p1:      Point{0, 0},
			p2:      Point{2, 3},
			expArea: 12,
		},
		{
			name:    "Rectangle 6x11 (inclusive)",
			p1:      Point{0, 0},
			p2:      Point{5, 10},
			expArea: 66,
		},
		{
			name:    "Square 6x6 (inclusive)",
			p1:      Point{0, 0},
			p2:      Point{5, 5},
			expArea: 36,
		},
		{
			name:    "Reversed coordinates (same result)",
			p1:      Point{10, 10},
			p2:      Point{0, 0},
			expArea: 121,
		},
		{
			name:    "Same point (area 1)",
			p1:      Point{5, 5},
			p2:      Point{5, 5},
			expArea: 1,
		},
		{
			name:    "Same x - vertical line (width 1)",
			p1:      Point{5, 0},
			p2:      Point{5, 10},
			expArea: 11,
		},
		{
			name:    "Same y - horizontal line (height 1)",
			p1:      Point{0, 5},
			p2:      Point{10, 5},
			expArea: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := abs(tt.p2.x - tt.p1.x) + 1
			height := abs(tt.p2.y - tt.p1.y) + 1
			area := width * height

			if area != tt.expArea {
				t.Errorf("Area calculation failed: expected %d, got %d", tt.expArea, area)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("Very close coordinates", func(t *testing.T) {
		points := []Point{
			{0, 0}, {1, 1},
		}
		result := findLargestRectangle(points)
		if result != 4 {
			t.Errorf("Expected 4, got %d", result)
		}
	})

	t.Run("Duplicate points", func(t *testing.T) {
		points := []Point{
			{5, 5}, {5, 5}, {10, 10},
		}
		result := findLargestRectangle(points)
		expected := 36
		if result != expected {
			t.Errorf("Expected %d, got %d", expected, result)
		}
	})

	t.Run("Three collinear points", func(t *testing.T) {
		points := []Point{
			{0, 0}, {5, 5}, {10, 10},
		}
		result := findLargestRectangle(points)
		if result != 121 {
			t.Errorf("Expected 121, got %d", result)
		}
	})

	t.Run("Many points, largest is last pair", func(t *testing.T) {
		points := []Point{
			{1, 1}, {2, 2}, {3, 3}, {0, 0}, {100, 100},
		}
		result := findLargestRectangle(points)
		if result != 10201 {
			t.Errorf("Expected 10201, got %d", result)
		}
	})
}

func BenchmarkFindLargestRectangle(b *testing.B) {
	// Create a set of points for benchmarking
	var points []Point
	for i := 0; i < 100; i++ {
		points = append(points, Point{i, i})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findLargestRectangle(points)
	}
}

// Part 2 Tests

func TestIdentifyRedGreenTiles(t *testing.T) {
	tests := []struct {
		name          string
		redTiles      []Point
		expectedCount int // at least this many red/green tiles
	}{
		{
			name: "Simple square boundary",
			redTiles: []Point{
				{0, 0}, {0, 2}, {2, 2}, {2, 0},
			},
			expectedCount: 8, // boundary tiles on perimeter (no interior since we skip flood fill)
		},
		{
			name: "Vertical line",
			redTiles: []Point{
				{0, 0}, {0, 2},
			},
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redGreen := identifyRedGreenTiles(tt.redTiles)
			if len(redGreen) < tt.expectedCount {
				t.Errorf("identifyRedGreenTiles: expected at least %d tiles, got %d", tt.expectedCount, len(redGreen))
			}
		})
	}
}

func TestAddLinePoints(t *testing.T) {
	tests := []struct {
		name     string
		p1       Point
		p2       Point
		expected int // number of points on the line
	}{
		{
			name:     "Vertical line from (0,0) to (0,5)",
			p1:       Point{0, 0},
			p2:       Point{0, 5},
			expected: 6,
		},
		{
			name:     "Horizontal line from (0,0) to (5,0)",
			p1:       Point{0, 0},
			p2:       Point{5, 0},
			expected: 6,
		},
		{
			name:     "Single point (same coordinates)",
			p1:       Point{3, 3},
			p2:       Point{3, 3},
			expected: 1,
		},
		{
			name:     "Vertical line reversed",
			p1:       Point{0, 5},
			p2:       Point{0, 0},
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redGreen := make(map[Point]bool)
			addLinePoints(tt.p1, tt.p2, redGreen)
			if len(redGreen) != tt.expected {
				t.Errorf("addLinePoints: expected %d points, got %d", tt.expected, len(redGreen))
			}
		})
	}
}

func TestGetBounds(t *testing.T) {
	tests := []struct {
		name      string
		points    []Point
		expMinX   int
		expMaxX   int
		expMinY   int
		expMaxY   int
	}{
		{
			name:      "Simple rectangle",
			points:    []Point{{0, 0}, {5, 10}},
			expMinX:   0,
			expMaxX:   5,
			expMinY:   0,
			expMaxY:   10,
		},
		{
			name:      "Negative coordinates",
			points:    []Point{{-5, -10}, {5, 10}},
			expMinX:   -5,
			expMaxX:   5,
			expMinY:   -10,
			expMaxY:   10,
		},
		{
			name:      "Single point",
			points:    []Point{{3, 7}},
			expMinX:   3,
			expMaxX:   3,
			expMinY:   7,
			expMaxY:   7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minX, maxX, minY, maxY := getBounds(tt.points)
			if minX != tt.expMinX || maxX != tt.expMaxX || minY != tt.expMinY || maxY != tt.expMaxY {
				t.Errorf("getBounds: expected (%d,%d,%d,%d), got (%d,%d,%d,%d)",
					tt.expMinX, tt.expMaxX, tt.expMinY, tt.expMaxY,
					minX, maxX, minY, maxY)
			}
		})
	}
}

func TestIsInsidePolygon(t *testing.T) {
	tests := []struct {
		name     string
		point    Point
		polygon  []Point
		expected bool
	}{
		{
			name:  "Point inside square",
			point: Point{2, 2},
			polygon: []Point{
				{0, 0}, {4, 0}, {4, 4}, {0, 4},
			},
			expected: true,
		},
		{
			name:  "Point outside square",
			point: Point{5, 5},
			polygon: []Point{
				{0, 0}, {4, 0}, {4, 4}, {0, 4},
			},
			expected: false,
		},
		{
			name:  "Point inside at (1, 1)",
			point: Point{1, 1},
			polygon: []Point{
				{0, 0}, {4, 0}, {4, 4}, {0, 4},
			},
			expected: true,
		},
		{
			name:  "Point outside at (-1, -1)",
			point: Point{-1, -1},
			polygon: []Point{
				{0, 0}, {4, 0}, {4, 4}, {0, 4},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isInsidePolygon(tt.point, tt.polygon)
			if result != tt.expected {
				t.Errorf("isInsidePolygon: expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsValidRectangle(t *testing.T) {
	// Create a simple red/green map
	redGreen := make(map[Point]bool)
	// Mark a 3x3 square as red/green
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			redGreen[Point{x, y}] = true
		}
	}

	tests := []struct {
		name     string
		p1       Point
		p2       Point
		expected bool
	}{
		{
			name:     "Valid rectangle within red/green area",
			p1:       Point{0, 0},
			p2:       Point{1, 1},
			expected: true,
		},
		{
			name:     "Valid rectangle larger area",
			p1:       Point{0, 0},
			p2:       Point{2, 2},
			expected: true,
		},
		{
			name:     "Invalid rectangle outside area",
			p1:       Point{0, 0},
			p2:       Point{5, 5},
			expected: false,
		},
		{
			name:     "Valid single point",
			p1:       Point{1, 1},
			p2:       Point{1, 1},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidRectangle(tt.p1, tt.p2, redGreen)
			if result != tt.expected {
				t.Errorf("isValidRectangle: expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFindLargestRectanglePart2(t *testing.T) {
	tests := []struct {
		name     string
		redTiles []Point
		expected int
	}{
		{
			name: "Simple loop - area 9 (3x3 square)",
			redTiles: []Point{
				{0, 0}, {0, 2}, {2, 2}, {2, 0},
			},
			expected: 9,
		},
		{
			name: "Simple horizontal line",
			redTiles: []Point{
				{0, 0}, {2, 0},
			},
			expected: 3,
		},
		{
			name: "Single point",
			redTiles: []Point{
				{5, 5},
			},
			expected: 0,
		},
		{
			name: "Two points vertical",
			redTiles: []Point{
				{0, 0}, {0, 2},
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findLargestRectanglePart2(tt.redTiles)
			if result != tt.expected {
				t.Errorf("findLargestRectanglePart2: expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func BenchmarkFindLargestRectanglePart2(b *testing.B) {
	// Create a rectangular loop
	var points []Point
	for x := 0; x < 20; x++ {
		points = append(points, Point{x, 0})
	}
	for y := 1; y < 20; y++ {
		points = append(points, Point{19, y})
	}
	for x := 18; x >= 0; x-- {
		points = append(points, Point{x, 19})
	}
	for y := 18; y > 0; y-- {
		points = append(points, Point{0, y})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findLargestRectanglePart2(points)
	}
}
