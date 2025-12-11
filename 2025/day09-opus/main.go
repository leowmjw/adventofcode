package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Point represents a 2D coordinate
type Point struct {
	X, Y int
}

// Interval represents a horizontal interval [Left, Right]
type Interval struct {
	Left, Right int
}

// HorizontalSlab represents a horizontal strip with y in [YMin, YMax)
type HorizontalSlab struct {
	YMin, YMax    int
	InsideXRanges []Interval // sorted, non-overlapping intervals
}

// Polygon holds the rectilinear polygon data
type Polygon struct {
	Vertices []Point
	Slabs    []HorizontalSlab // sorted by YMin
	YCoords  []int            // sorted unique Y coordinates
}

func main() {
	runtime.GOMAXPROCS(16)
	start := time.Now()

	// Parse input
	redTiles := parseInput("/tmp/aoc-day9/input.txt")
	fmt.Printf("Parsed %d red tiles (polygon vertices)\n", len(redTiles))

	// Build polygon structure
	polygon := buildPolygon(redTiles)
	fmt.Printf("Built polygon with %d horizontal slabs\n", len(polygon.Slabs))
	fmt.Printf("Polygon build time: %v\n", time.Since(start))

	// Find largest valid rectangle using concurrent search
	searchStart := time.Now()
	answer := findLargestRectangleConcurrent(redTiles, polygon, 16)
	fmt.Printf("Search time: %v\n", time.Since(searchStart))
	fmt.Printf("Total time: %v\n", time.Since(start))
	fmt.Printf("\n=== Answer (Part 2): %d ===\n", answer)
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
		x, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		y, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		points = append(points, Point{X: x, Y: y})
	}
	return points
}

// buildPolygon constructs the slab decomposition for efficient containment queries
func buildPolygon(vertices []Point) *Polygon {
	n := len(vertices)
	if n == 0 {
		return &Polygon{}
	}

	// Collect all unique Y coordinates
	ySet := make(map[int]bool)
	for _, v := range vertices {
		ySet[v.Y] = true
	}

	yCoords := make([]int, 0, len(ySet))
	for y := range ySet {
		yCoords = append(yCoords, y)
	}
	sort.Ints(yCoords)

	// Build vertical edges of the polygon
	type VerticalEdge struct {
		X          int
		YMin, YMax int
	}

	var verticalEdges []VerticalEdge
	for i := 0; i < n; i++ {
		p1 := vertices[i]
		p2 := vertices[(i+1)%n]

		if p1.X == p2.X {
			// Vertical edge
			yMin, yMax := p1.Y, p2.Y
			if yMin > yMax {
				yMin, yMax = yMax, yMin
			}
			verticalEdges = append(verticalEdges, VerticalEdge{X: p1.X, YMin: yMin, YMax: yMax})
		}
	}

	// Build slabs using scanline
	slabs := make([]HorizontalSlab, 0, len(yCoords)-1)

	for i := 0; i < len(yCoords)-1; i++ {
		yMin := yCoords[i]
		yMax := yCoords[i+1]

		// Use midpoint for crossing test (avoids boundary issues)
		yMid := float64(yMin+yMax) / 2.0

		// Find all vertical edges that cross through this slab's interior
		var crossingX []int
		for _, edge := range verticalEdges {
			// Edge crosses the slab if its Y range strictly contains the midpoint
			if float64(edge.YMin) < yMid && float64(edge.YMax) > yMid {
				crossingX = append(crossingX, edge.X)
			}
		}
		sort.Ints(crossingX)

		// Pair up crossings to form inside intervals (scanline fill)
		// Going left-to-right, we toggle in/out at each crossing
		var insideRanges []Interval
		for j := 0; j+1 < len(crossingX); j += 2 {
			insideRanges = append(insideRanges, Interval{
				Left:  crossingX[j],
				Right: crossingX[j+1],
			})
		}

		slabs = append(slabs, HorizontalSlab{
			YMin:          yMin,
			YMax:          yMax,
			InsideXRanges: insideRanges,
		})
	}

	return &Polygon{
		Vertices: vertices,
		Slabs:    slabs,
		YCoords:  yCoords,
	}
}

// isRectangleValid checks if a rectangle with opposite corners (x1,y1) and (x2,y2)
// is fully contained within the polygon
func (p *Polygon) isRectangleValid(x1, y1, x2, y2 int) bool {
	// Normalize coordinates
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	// Find first slab that could overlap (where YMax > y1)
	startIdx := sort.Search(len(p.Slabs), func(i int) bool {
		return p.Slabs[i].YMax > y1
	})

	// Check each relevant slab
	for i := startIdx; i < len(p.Slabs); i++ {
		slab := &p.Slabs[i]

		// Stop if slab is entirely above our rectangle
		if slab.YMin >= y2 {
			break
		}

		// Check if rectangle's x-range [x1, x2] is fully contained in some inside interval
		contained := false
		for _, interval := range slab.InsideXRanges {
			if interval.Left <= x1 && x2 <= interval.Right {
				contained = true
				break
			}
		}

		if !contained {
			return false
		}
	}

	return true
}

// findLargestRectangleConcurrent uses goroutines to parallelize the search
func findLargestRectangleConcurrent(redTiles []Point, polygon *Polygon, numWorkers int) int {
	n := len(redTiles)
	if n < 2 {
		return 0
	}

	var globalMax atomic.Int64
	var wg sync.WaitGroup

	// Create work channel
	workChan := make(chan int, n)

	// Progress tracking
	var processed atomic.Int64

	// Spawn workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			localMax := int64(0)

			for i := range workChan {
				p1 := redTiles[i]

				for j := i + 1; j < n; j++ {
					p2 := redTiles[j]

					// Skip if same row or column
					if p1.X == p2.X || p1.Y == p2.Y {
						continue
					}

					// Calculate area (including boundary tiles)
					// Rectangle from (x1,y1) to (x2,y2) has width |x2-x1|+1 and height |y2-y1|+1
					width := int64(abs(p2.X-p1.X) + 1)
					height := int64(abs(p2.Y-p1.Y) + 1)
					area := width * height

					// Skip if can't beat current max
					if area <= localMax {
						continue
					}

					// Also skip if can't beat global max
					if area <= globalMax.Load() {
						continue
					}

					// Check if rectangle is valid
					if polygon.isRectangleValid(p1.X, p1.Y, p2.X, p2.Y) {
						localMax = area
						// Update global max immediately for better pruning
						for {
							old := globalMax.Load()
							if localMax <= old {
								break
							}
							if globalMax.CompareAndSwap(old, localMax) {
								fmt.Printf("New max found: %d\n", localMax)
								break
							}
						}
					}
				}

				count := processed.Add(1)
				if count%50 == 0 {
					fmt.Printf("Progress: %d/%d vertices (%.1f%%)\n", count, n, float64(count)*100/float64(n))
				}
			}
		}(w)
	}

	// Distribute work
	for i := 0; i < n; i++ {
		workChan <- i
	}
	close(workChan)

	wg.Wait()
	return int(globalMax.Load())
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
