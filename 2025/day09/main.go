package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

func main() {
	fmt.Println("Day 09: Movie Theater - Finding Largest Rectangle (Part 2)")
	result := run("input.txt")
	fmt.Printf("Answer: %d\n", result)
}

func run(inputFile string) int {
	points := parseInput(inputFile)
	return findLargestRectanglePart2(points)
}

func parseInput(filename string) []Point {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var points []Point
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		points = append(points, Point{x, y})
	}
	return points
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Identify all tiles that are red or green
// Red tiles are the input points
// Green tiles are those on the boundary path and inside the loop
func identifyRedGreenTiles(redTiles []Point) map[Point]bool {
	redGreen := make(map[Point]bool)

	// Add all red tiles
	for _, p := range redTiles {
		redGreen[p] = true
	}

	// Connect red tiles with green tiles
	// Red tiles are connected by straight lines of green tiles
	for i := 0; i < len(redTiles); i++ {
		from := redTiles[i]
		to := redTiles[(i+1)%len(redTiles)] // wrap around to first

		// Add all points on the line between from and to
		addLinePoints(from, to, redGreen)
	}

	// Flood fill to find interior tiles
	floodFillInterior(redGreen, redTiles, struct{ minX, maxX, minY, maxY int }{0, 0, 0, 0})

	return redGreen
}

// Add all points on a straight line between p1 and p2
func addLinePoints(p1, p2 Point, redGreen map[Point]bool) {
	if p1.x == p2.x {
		// Vertical line
		minY := p1.y
		maxY := p2.y
		if minY > maxY {
			minY, maxY = maxY, minY
		}
		for y := minY; y <= maxY; y++ {
			redGreen[Point{p1.x, y}] = true
		}
	} else if p1.y == p2.y {
		// Horizontal line
		minX := p1.x
		maxX := p2.x
		if minX > maxX {
			minX, maxX = maxX, minX
		}
		for x := minX; x <= maxX; x++ {
			redGreen[Point{x, p1.y}] = true
		}
	}
}

// Get bounding box of all red tiles
func getBounds(points []Point) (minX, maxX, minY, maxY int) {
	if len(points) == 0 {
		return 0, 0, 0, 0
	}
	minX, maxX = points[0].x, points[0].x
	minY, maxY = points[0].y, points[0].y

	for _, p := range points {
		if p.x < minX {
			minX = p.x
		}
		if p.x > maxX {
			maxX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
		if p.y > maxY {
			maxY = p.y
		}
	}
	return
}

// Flood fill from boundary to mark interior points as green
func floodFillInterior(redGreen map[Point]bool, redTiles []Point, bounds struct {
	minX, maxX, minY, maxY int
}) {
	// Skip flood filling for now - rectangles will only include boundary tiles
	// The interior is implicitly green where the rectangle is valid
}

// Ray casting algorithm to check if a point is inside a polygon
func isInsidePolygon(p Point, polygon []Point) bool {
	if len(polygon) < 3 {
		return false
	}

	count := 0
	for i := 0; i < len(polygon); i++ {
		p1 := polygon[i]
		p2 := polygon[(i+1)%len(polygon)]

		// Check if horizontal ray from p to the right crosses edge p1->p2
		if (p1.y <= p.y && p.y < p2.y) || (p2.y <= p.y && p.y < p1.y) {
			// Calculate x-coordinate of intersection
			xinters := float64(p1.x) + float64(p.y-p1.y)*float64(p2.x-p1.x)/float64(p2.y-p1.y)
			if float64(p.x) < xinters {
				count++
			}
		}
	}

	return count%2 == 1
}

// Check if rectangle contains only red/green tiles
func isValidRectangle(p1, p2 Point, redGreen map[Point]bool) bool {
	minX := p1.x
	maxX := p2.x
	if minX > maxX {
		minX, maxX = maxX, minX
	}

	minY := p1.y
	maxY := p2.y
	if minY > maxY {
		minY, maxY = maxY, minY
	}

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			if !redGreen[Point{x, y}] {
				return false
			}
		}
	}
	return true
}

// Part 2: Find largest rectangle with red corners that contains only red/green tiles
func findLargestRectanglePart2(redTiles []Point) int {
	if len(redTiles) < 2 {
		return 0
	}

	redGreen := identifyRedGreenTiles(redTiles)

	// Pre-compute all valid red tile pairs that form axis-aligned rectangles
	// Only consider rectangles where width and height are both positive
	maxArea := 0

	// Try all pairs of red tiles as opposite corners
	for i := 0; i < len(redTiles); i++ {
		for j := i + 1; j < len(redTiles); j++ {
			p1 := redTiles[i]
			p2 := redTiles[j]

			// Skip if not a valid rectangle (same x or same y won't be interesting for part 2)
			// But we still need to check them
			width := abs(p2.x - p1.x) + 1
			height := abs(p2.y - p1.y) + 1
			area := width * height

			// Quick bounds check - skip obviously invalid large rectangles
			if area > maxArea {
				// Check if rectangle contains only red/green tiles
				if isValidRectangleFast(p1, p2, redGreen, redTiles) {
					maxArea = area
				}
			}
		}
	}

	return maxArea
}

// Optimized validation that fails fast
// Uses ray casting to check if all rectangle interior points are inside the polygon
func isValidRectangleFast(p1, p2 Point, redGreen map[Point]bool, redTiles []Point) bool {
	minX := p1.x
	maxX := p2.x
	if minX > maxX {
		minX, maxX = maxX, minX
	}

	minY := p1.y
	maxY := p2.y
	if minY > maxY {
		minY, maxY = maxY, minY
	}

	// Check all points on boundary path are in redGreen
	for x := minX; x <= maxX; x++ {
		// Check top and bottom edges
		if !redGreen[Point{x, minY}] || !redGreen[Point{x, maxY}] {
			return false
		}
	}
	for y := minY; y <= maxY; y++ {
		// Check left and right edges
		if !redGreen[Point{minX, y}] || !redGreen[Point{maxX, y}] {
			return false
		}
	}

	// Check interior points are inside the polygon (green)
	for x := minX + 1; x < maxX; x++ {
		for y := minY + 1; y < maxY; y++ {
			if !isInsidePolygon(Point{x, y}, redTiles) {
				return false
			}
		}
	}
	return true
}

// Part 1: Find largest rectangle without constraints
func findLargestRectangle(points []Point) int {
	if len(points) < 2 {
		return 0
	}

	maxArea := 0

	// Try all pairs of points as opposite corners
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			p1 := points[i]
			p2 := points[j]

			// Calculate rectangle area with opposite corners at p1 and p2
			// Area includes both boundary cells, so we add 1 to each dimension
			width := abs(p2.x - p1.x) + 1
			height := abs(p2.y - p1.y) + 1
			area := width * height

			if area > maxArea {
				maxArea = area
			}
		}
	}

	return maxArea
}
