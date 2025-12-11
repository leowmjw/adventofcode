# AGENTS.md - AoC 2025 Day 9 Part 2 Solution

## Project Overview

This Go solution solves Advent of Code 2025 Day 9 Part 2: "Movie Theater" puzzle.

### Problem Summary

- **Part 1**: Find the largest rectangle using any two red tiles as opposite corners. Answer: 4,771,508,457
- **Part 2**: Red tiles form vertices of a rectilinear polygon. Green tiles are on edges between consecutive red vertices and inside the polygon. Find the largest rectangle with red tile corners where ALL tiles in the rectangle are red or green (fully contained in the polygon). **Answer: 1,539,809,693**

## Solution Architecture

### Algorithm

1. **Polygon Construction** - O(n log n)
   - Parse red tile coordinates (polygon vertices in order)
   - Extract unique Y coordinates and sort them
   - Build vertical edges from consecutive vertex pairs
   - Create horizontal slabs between Y coordinates
   - For each slab, determine "inside" X intervals using scanline algorithm

2. **Rectangle Validation** - O(log n + k) per query
   - Binary search to find relevant slabs
   - Check if rectangle's X range is contained in slab's inside intervals
   - Rectangle is valid iff it passes all slab checks

3. **Concurrent Search** - O(n²/P) where P = number of CPUs
   - Distribute vertex pairs across goroutines
   - Each worker maintains local maximum
   - Aggressive pruning: skip pairs that can't beat current best
   - Atomic updates to global maximum

### Key Insight: Area Calculation

The rectangle area includes boundary tiles:
- Width = |x2 - x1| + 1
- Height = |y2 - y1| + 1
- Area = Width × Height

For example, rectangle from (2,3) to (9,5):
- Width = 9 - 2 + 1 = 8 tiles
- Height = 5 - 3 + 1 = 3 tiles
- Area = 24 tiles

### Data Structures

| Structure | Purpose | Complexity |
|-----------|---------|------------|
| `Point` | 2D coordinate (X, Y) | O(1) |
| `Interval` | Horizontal range [Left, Right] | O(1) |
| `HorizontalSlab` | Y range + inside X intervals | O(k) where k = intervals |
| `Polygon` | Vertices + sorted slabs + Y coords | O(n) |

### Key Files

```
aoc-day9/
├── main.go          # Main solution implementation
├── main_test.go     # Comprehensive unit tests
├── input.txt        # Puzzle input (496 red tiles)
├── go.mod           # Go module definition
├── AGENTS.md        # This file
└── solver           # Compiled binary
```

## Algorithm Deep Dive

### Slab Decomposition

For a rectilinear polygon, we decompose the plane into horizontal slabs:

```
Y coordinates: [1, 3, 5, 7]
Slabs:
  - Slab 0: Y ∈ [1, 3), inside X intervals: [{7, 11}]
  - Slab 1: Y ∈ [3, 5), inside X intervals: [{2, 11}]
  - Slab 2: Y ∈ [5, 7), inside X intervals: [{9, 11}]
```

Each slab's "inside" intervals are determined by counting vertical edge crossings using a scanline algorithm. An odd number of crossings from -∞ means we're inside.

### Rectangle Containment Check

A rectangle R with corners (x1,y1) and (x2,y2) is fully inside the polygon iff:
- For every slab S that R intersects vertically
- R's horizontal extent [min(x1,x2), max(x1,x2)] is contained in one of S's inside intervals

### Concurrency Model

```
┌─────────────────────────────────────────┐
│              Main Thread                │
│  - Creates work channel                 │
│  - Spawns 16 worker goroutines          │
│  - Distributes vertex indices           │
└─────────────────────────────────────────┘
                    │
        ┌───────────┼───────────┐
        ▼           ▼           ▼
   ┌─────────┐ ┌─────────┐ ┌─────────┐
   │Worker 1 │ │Worker 2 │ │Worker N │
   │         │ │         │ │   ...   │
   │Local max│ │Local max│ │Local max│
   └────┬────┘ └────┬────┘ └────┬────┘
        │           │           │
        └───────────┼───────────┘
                    ▼
        ┌─────────────────────┐
        │  Atomic Global Max  │
        │  (CompareAndSwap)   │
        └─────────────────────┘
```

## Performance Characteristics

### Time Complexity

| Operation | Complexity |
|-----------|------------|
| Parse input | O(n) |
| Build polygon | O(n log n) |
| Single rectangle check | O(log n + k) |
| All pairs search | O(n² × (log n + k) / P) |

Where:
- n = number of red tiles (vertices)
- k = average intervals per slab
- P = number of CPU cores (16)

### Space Complexity

| Component | Space |
|-----------|-------|
| Vertices | O(n) |
| Slabs | O(n) |
| Y coordinates | O(n) |
| Worker goroutines | O(P) |
| **Total** | **O(n)** |

### Actual Performance (496 vertices, 16 CPUs)

```
Polygon build time: ~0.4ms
Search time: ~0.8ms
Total time: ~1.3ms
```

## Testing

### Run Tests

```bash
go test -v ./...
```

### Run Benchmarks

```bash
go test -bench=. -benchmem
```

### Test Coverage

```bash
go test -cover ./...
```

### Test Categories

1. **Unit Tests**
   - `TestParseInput` - Input parsing with normal and whitespace input
   - `TestBuildPolygon` - Polygon construction and slab generation
   - `TestIsRectangleValid` - Rectangle validation for various cases
   - `TestAreaCalculation` - Verify area includes boundary tiles
   - `TestAbs` - Utility functions

2. **Integration Tests**
   - `TestFindLargestRectangleConcurrent` - Full algorithm with example data
   - `TestFindLargestRectangleSimpleSquare` - Simple square polygon
   - `TestRealInput` - Actual puzzle input

3. **Edge Case Tests**
   - Empty input
   - Single point
   - Collinear points (no valid rectangle)

4. **Concurrency Tests**
   - `TestConcurrencyCorrectness` - Race condition detection (runs 10x)

5. **Benchmarks**
   - `BenchmarkBuildPolygon` - Polygon construction speed
   - `BenchmarkIsRectangleValid` - Rectangle validation speed
   - `BenchmarkFindLargestRectangle*` - Full algorithm with varying workers

## Usage

### Build

```bash
go build -o solver main.go
```

### Run

```bash
./solver
```

### With Custom Input

Modify `input.txt` or change the path in `main()`.

## Solution Results

| Part | Answer |
|------|--------|
| Part 1 | 4,771,508,457 |
| **Part 2** | **1,539,809,693** |

## Key Insights

1. **Rectilinear Polygon Property**: All edges are axis-aligned, enabling efficient slab decomposition.

2. **Scanline Algorithm**: Classic computational geometry technique for inside/outside determination. Vertical edges toggle the inside/outside state.

3. **Coordinate Compression**: Only Y coordinates of vertices matter for slab boundaries.

4. **Embarrassingly Parallel**: Pair checking has no dependencies between pairs, perfect for goroutines with work-stealing via channels.

5. **Early Pruning**: Checking area before validation and updating global max immediately provides significant speedup.

6. **Area Includes Boundaries**: Rectangle area = (|x2-x1|+1) × (|y2-y1|+1), not just |x2-x1| × |y2-y1|.

## References

- [Advent of Code 2025](https://adventofcode.com/2025)
- [Rectilinear Polygon](https://en.wikipedia.org/wiki/Rectilinear_polygon)
- [Scanline Algorithm](https://en.wikipedia.org/wiki/Scanline_rendering)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)

## Author

Solution developed with Claude (Anthropic) for AoC 2025 Day 9.
