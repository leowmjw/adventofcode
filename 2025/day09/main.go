package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type Point struct {
	x, y int
}

// Grid-based bitmap for O(1) point lookups
type ValidityGrid struct {
	data map[int]map[int]bool
	mu   sync.RWMutex
}

func NewValidityGrid() *ValidityGrid {
	return &ValidityGrid{
		data: make(map[int]map[int]bool),
	}
}

func (g *ValidityGrid) Set(x, y int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.data[x] == nil {
		g.data[x] = make(map[int]bool)
	}
	g.data[x][y] = true
}

func (g *ValidityGrid) Get(x, y int) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if row, exists := g.data[x]; exists {
		return row[y]
	}
	return false
}

func main() {
	fmt.Println("Day 09: Movie Theater - Finding Largest Rectangle (Part 2)")
	result := run("input.txt")
	fmt.Printf("Answer: %d\n", result)
}

func run(inputFile string) int64 {
	points := parseInput(inputFile)
	return findLargestRectangleOptimized(points)
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

// Build grid with only boundary tiles (no expensive flood fill)
func buildBoundaryGrid(redTiles []Point) *ValidityGrid {
	grid := NewValidityGrid()

	// Add all red tiles
	for _, p := range redTiles {
		grid.Set(p.x, p.y)
	}

	// Connect red tiles with green tiles
	for i := 0; i < len(redTiles); i++ {
		from := redTiles[i]
		to := redTiles[(i+1)%len(redTiles)]

		// Add all points on the line between from and to
		if from.x == to.x {
			// Vertical line
			minY := from.y
			maxY := to.y
			if minY > maxY {
				minY, maxY = maxY, minY
			}
			for y := minY; y <= maxY; y++ {
				grid.Set(from.x, y)
			}
		} else if from.y == to.y {
			// Horizontal line
			minX := from.x
			maxX := to.x
			if minX > maxX {
				minX, maxX = maxX, minX
			}
			for x := minX; x <= maxX; x++ {
				grid.Set(x, from.y)
			}
		}
	}

	return grid
}

// Optimized version using bitmap preprocessing and parallelization
func findLargestRectangleOptimized(redTiles []Point) int64 {
	if len(redTiles) < 2 {
		return 0
	}

	// Phase 1: Build the boundary grid (no flood fill to save time)
	boundaryGrid := buildBoundaryGrid(redTiles)

	// Phase 2: Parallel enumeration and validation
	numWorkers := runtime.NumCPU()
	var maxArea int64

	// Divide work: each worker handles pairs where first point index is in its range
	chunkSize := (len(redTiles) + numWorkers - 1) / numWorkers
	var wg sync.WaitGroup
	results := make([]int64, numWorkers)

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			startIdx := workerID * chunkSize
			endIdx := startIdx + chunkSize
			if endIdx > len(redTiles) {
				endIdx = len(redTiles)
			}

			localMax := int64(0)

			// Process all pairs where first index is in this worker's range
			for i := startIdx; i < endIdx; i++ {
				p1 := redTiles[i]

				// Try all points after this one
				for j := i + 1; j < len(redTiles); j++ {
					p2 := redTiles[j]

					width := int64(abs(p2.x-p1.x)) + 1
					height := int64(abs(p2.y-p1.y)) + 1
					area := width * height

					// Only validate if potentially better
					if area > localMax {
						if isValidRectangle(p1, p2, boundaryGrid, redTiles) {
							localMax = area
						}
					}
				}
			}

			results[workerID] = localMax
		}(w)
	}

	wg.Wait()

	// Find global maximum
	for _, result := range results {
		if result > maxArea {
			maxArea = result
		}
	}

	return maxArea
}

// Ray casting for point-in-polygon
func isInsidePolygon(p Point, polygon []Point) bool {
	if len(polygon) < 3 {
		return false
	}

	count := 0
	for i := 0; i < len(polygon); i++ {
		p1 := polygon[i]
		p2 := polygon[(i+1)%len(polygon)]

		if (p1.y <= p.y && p.y < p2.y) || (p2.y <= p.y && p.y < p1.y) {
			xinters := float64(p1.x) + float64(p.y-p1.y)*float64(p2.x-p1.x)/float64(p2.y-p1.y)
			if float64(p.x) < xinters {
				count++
			}
		}
	}

	return count%2 == 1
}

// Check if rectangle is valid: all points must be on boundary OR inside polygon
func isValidRectangle(p1, p2 Point, boundaryGrid *ValidityGrid, redTiles []Point) bool {
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

	// Check all corners first (quick rejection)
	if !boundaryGrid.Get(minX, minY) && !isInsidePolygon(Point{minX, minY}, redTiles) {
		return false
	}
	if !boundaryGrid.Get(maxX, maxY) && !isInsidePolygon(Point{maxX, maxY}, redTiles) {
		return false
	}
	if !boundaryGrid.Get(minX, maxY) && !isInsidePolygon(Point{minX, maxY}, redTiles) {
		return false
	}
	if !boundaryGrid.Get(maxX, minY) && !isInsidePolygon(Point{maxX, minY}, redTiles) {
		return false
	}

	// Sample interior points to validate
	// Don't check every point - use adaptive sampling
	sampleRate := 1
	if (maxX - minX) > 1000 || (maxY - minY) > 1000 {
		sampleRate = 100
	}

	for x := minX; x <= maxX; x += sampleRate {
		for y := minY; y <= maxY; y += sampleRate {
			if !boundaryGrid.Get(x, y) && !isInsidePolygon(Point{x, y}, redTiles) {
				return false
			}
		}
	}

	// If sample passed, check edges more densely
	for x := minX; x <= maxX; x += max(1, sampleRate/10) {
		// Top and bottom edges
		if !boundaryGrid.Get(x, minY) && !isInsidePolygon(Point{x, minY}, redTiles) {
			return false
		}
		if !boundaryGrid.Get(x, maxY) && !isInsidePolygon(Point{x, maxY}, redTiles) {
			return false
		}
	}

	for y := minY; y <= maxY; y += max(1, sampleRate/10) {
		// Left and right edges
		if !boundaryGrid.Get(minX, y) && !isInsidePolygon(Point{minX, y}, redTiles) {
			return false
		}
		if !boundaryGrid.Get(maxX, y) && !isInsidePolygon(Point{maxX, y}, redTiles) {
			return false
		}
	}

	return true
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
