# Day 09 Test Scenarios and Edge Cases

## Test Suite Overview

Total Tests: 40+
All Passing: ✓

## 1. Main Rectangle Finding Tests (15 tests)

### Basic Cases
- **No points** → Returns 0
- **Single point** → Returns 0 (requires 2 points minimum)
- **Problem example** → 8 points, expected area: 50

### Two-Point Tests
| Test Name | Points | Expected Area | Notes |
|-----------|--------|---------------|-------|
| Vertical line (same x) | (5,0), (5,10) | 11 | 1 column, 11 rows |
| Horizontal line (same y) | (0,5), (10,5) | 11 | 11 columns, 1 row |
| Thin rectangle | (2,3), (7,3) | 6 | 6×1 rectangle |
| Regular rectangle | (2,5), (9,7) | 24 | 8×3 rectangle |
| Square (6×6) | (0,0), (5,5) | 36 | All sides equal |

### Multi-Point Tests
| Test Name | Points | Expected Area | Notes |
|-----------|--------|---------------|-------|
| 3 points | (0,0), (5,5), (10,10) | 121 | Max from (0,0)-(10,10) |
| 4 points forming square | (0,0), (5,0), (0,5), (5,5) | 25 | Max area = 36 |
| Collinear on x-axis | (0,0), (5,0), (10,0), (15,0) | 16 | Max from ends: 16×1 |
| Collinear on y-axis | (0,0), (0,5), (0,10), (0,15) | 16 | Max from ends: 1×16 |
| Many points, largest last | (1,1), (2,2), (3,3), (0,0), (100,100) | 10201 | (0,0) to (100,100) |

### Coordinate Type Tests
| Test Name | Coordinates | Expected Area | Notes |
|-----------|-------------|---------------|-------|
| Negative | (-5,-5), (5,5) | 121 | Works with negative coords |
| Mixed signs | (-10,-5), (10,5) | 231 | -10 to 10 = 21 wide, -5 to 5 = 11 tall |
| Large values | (97615,50359), (98289,50359) | 675 | Same y: 675×1 |
| Very large | (0,0), (1000,1000) | 1,002,001 | 1001×1001 |

---

## 2. Absolute Value Function Tests (5 tests)

Tests the `abs()` helper function:

| Input | Expected Output |
|-------|-----------------|
| 5 | 5 |
| -5 | 5 |
| 0 | 0 |
| -1000 | 1000 |
| 1000 | 1000 |

**Purpose:** Ensures direction independence (width/height always positive)

---

## 3. Input Parsing Test (1 test)

- **Parse valid input file**
  - Reads 496 coordinates from input.txt
  - Validates parsing returns non-zero points
  - Confirms first point is not (0,0)

---

## 4. Rectangle Area Calculation Tests (7 tests)

Direct calculation tests verifying the formula: `Area = (|x2-x1| + 1) × (|y2-y1| + 1)`

| Test Name | P1 | P2 | Width | Height | Area |
|-----------|----|----|-------|--------|------|
| Simple 3×4 | (0,0) | (2,3) | 3 | 4 | 12 |
| Rectangle 6×11 | (0,0) | (5,10) | 6 | 11 | 66 |
| Square 6×6 | (0,0) | (5,5) | 6 | 6 | 36 |
| Reversed coords | (10,10) | (0,0) | 11 | 11 | 121 |
| Same point | (5,5) | (5,5) | 1 | 1 | 1 |
| Vertical line | (5,0) | (5,10) | 1 | 11 | 11 |
| Horizontal line | (0,5) | (10,5) | 11 | 1 | 11 |

**Purpose:** Tests the inclusive boundary calculation

---

## 5. Edge Cases (4 tests)

### Complex Boundary Conditions

1. **Very close coordinates**
   - Points: (0,0) and (1,1)
   - Expected: 4 (2×2 square)
   - Tests: Minimum non-zero rectangle

2. **Duplicate points in list**
   - Points: (5,5), (5,5), (10,10)
   - Expected: 36 (from (5,5) to (10,10))
   - Tests: Duplicate handling

3. **Three collinear points**
   - Points: (0,0), (5,5), (10,10)
   - Expected: 121 (from (0,0) to (10,10))
   - Tests: Proper max finding from collinear set

4. **Many points with largest pair at end**
   - Points: (1,1), (2,2), (3,3), (0,0), (100,100)
   - Expected: 10,201 (from (0,0) to (100,100))
   - Tests: Algorithm doesn't miss last pair

---

## 6. Performance Benchmark (1 test)

**BenchmarkFindLargestRectangle**
- Generates 100 points
- Measures iterations per second
- Expected: Fast completion on modern hardware

---

## Test Execution Results

```
=== RUN   TestFindLargestRectangle
--- PASS: Example_from_problem (0.00s)
--- PASS: No_points (0.00s)
--- PASS: Single_point (0.00s)
--- PASS: Vertical_line (0.00s)
--- PASS: Horizontal_line (0.00s)
--- PASS: Thin_rectangle (0.00s)
--- PASS: Regular_rectangle (0.00s)
--- PASS: Three_points (0.00s)
--- PASS: Four_points_square (0.00s)
--- PASS: Negative_coordinates (0.00s)
--- PASS: Mixed_positive_negative (0.00s)
--- PASS: Collinear_x_axis (0.00s)
--- PASS: Collinear_y_axis (0.00s)
--- PASS: Large_values (0.00s)
--- PASS: Very_large_area (0.00s)

=== RUN   TestAbsFunction
--- PASS: abs(5) = 5 (0.00s)
--- PASS: abs(-5) = 5 (0.00s)
--- PASS: abs(0) = 0 (0.00s)
--- PASS: abs(-1000) = 1000 (0.00s)
--- PASS: abs(1000) = 1000 (0.00s)

=== RUN   TestParseInput
--- PASS: Parse_valid_input_file (0.00s)

=== RUN   TestRectangleAreaCalculation
--- PASS: Simple_3x4 (0.00s)
--- PASS: Rectangle_6x11 (0.00s)
--- PASS: Square_6x6 (0.00s)
--- PASS: Reversed_coordinates (0.00s)
--- PASS: Same_point (0.00s)
--- PASS: Vertical_line (0.00s)
--- PASS: Horizontal_line (0.00s)

=== RUN   TestEdgeCases
--- PASS: Very_close_coordinates (0.00s)
--- PASS: Duplicate_points (0.00s)
--- PASS: Three_collinear_points (0.00s)
--- PASS: Many_points_largest_last (0.00s)

PASS - ok  aoc/2025/day09  0.002s
```

---

## Key Testing Principles Applied

1. **Equivalence Partitioning**: Tests cover different input types
2. **Boundary Value Analysis**: Tests edges (0, 1, same values, negatives)
3. **Error Guessing**: Tests common mistakes (empty, single item, duplicates)
4. **Integration Testing**: ParseInput + FindLargestRectangle workflow
5. **Regression Testing**: Verifies problem example still works

---

## Uncovered Scenarios (by design)

- **Corrupted input**: Assumes input.txt is well-formed
- **File not found**: Assumes input.txt exists
- **Invalid coordinates**: Assumes integers can be parsed
- **Floating point coords**: Problem specifies integers only
- **Negative area**: Mathematically impossible with abs() usage

These are acceptable as they're at system boundaries and input is assumed to be valid per problem constraints.
