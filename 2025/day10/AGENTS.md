# AGENTS.md - Day 10 Learnings

## Part 1: Light Toggle / XOR System

### Key Insight
When you have a system where:
- Elements have binary states (on/off)
- Actions toggle specific elements
- Toggling twice returns to original state

This is a **linear algebra problem over GF(2)** (binary field / XOR operations).

### Critical Realization
Each button only needs to be pressed **0 or 1 times** because:
- Pressing twice = no effect (XOR property)
- This reduces search space from infinite to 2^n combinations

### Solution Approach

1. **Brute Force (when n ≤ ~20 buttons)**
   - Enumerate all 2^n combinations
   - For each combination, simulate and check if target achieved
   - Track minimum number of presses
   - Time: O(2^n * m) where m = number of lights

2. **Gaussian Elimination over GF(2) (for larger n)**
   - Treat as system of linear equations: Ax = b (mod 2)
   - Use Gaussian elimination adapted for binary field
   - Find minimum weight solution in null space
   - More complex but polynomial time

### Implementation Notes

```go
// XOR toggle pattern
lights[idx] = !lights[idx]  // Toggle is simply NOT operation

// Enumerate all 2^n combinations using bitmask
for combo := 0; combo < (1 << numButtons); combo++ {
    for i := 0; i < numButtons; i++ {
        if (combo & (1 << i)) != 0 {
            // Button i is pressed
        }
    }
}
```

---

## Part 2: Joltage Counters / Integer Linear Programming

### Problem Transformation
Part 2 changes the problem fundamentally:
- **Part 1**: Toggle (XOR) - pressing twice cancels out → binary choice (0 or 1)
- **Part 2**: Increment (addition) - pressing multiple times accumulates → non-negative integers

This transforms from a **GF(2) problem** to an **Integer Linear Programming (ILP)** problem.

### Mathematical Formulation
Given:
- `n` buttons and `m` counters (joltages)
- Matrix `A[m][n]` where `A[i][j] = 1` if button `j` affects counter `i`
- Target vector `b[m]` (the joltage requirements)
- Variable vector `x[n]` (number of presses per button)

**Objective**: Minimize `sum(x)` subject to `A*x = b` and `x >= 0` (integers)

### Solution Approach: Gaussian Elimination + Search

1. **Gaussian Elimination** reduces the augmented matrix `[A|b]` to row echelon form
2. Identify **pivot variables** (determined by equations) and **free variables** (can be any value)
3. If no free variables → unique solution, compute directly
4. If free variables exist → search over their possible values to minimize total

### Why Gaussian Elimination Works
The system `A*x = b` may be:
- **Overdetermined** (more equations than variables) - may have no solution
- **Underdetermined** (more variables than equations) - infinite solutions, parameterized by free variables
- **Exactly determined** - unique solution

Gaussian elimination handles all cases by:
1. Reducing to row echelon form
2. Detecting inconsistencies (0 = nonzero rows)
3. Identifying the solution space dimension (number of free variables)

### Implementation Details

```go
// Gaussian elimination with partial pivoting
for col := 0; col < numButtons && row < numCounters; col++ {
    // Find best pivot (largest absolute value for numerical stability)
    maxRow := row
    for i := row + 1; i < numCounters; i++ {
        if abs(matrix[i][col]) > abs(matrix[maxRow][col]) {
            maxRow = i
        }
    }

    // Swap rows, scale pivot, eliminate column
    // ... (standard Gaussian elimination)
}

// Back-substitution with free variable search
// For each combination of free variable values, compute pivot variables
// and check for valid non-negative integer solution
```

### Performance Analysis
- Gaussian elimination: O(n²m) where n=buttons, m=counters
- Free variable search: O(maxVal^k) where k=number of free variables
- In practice, k is usually small (0-2) due to the problem structure
- Total runtime for 156 machines: < 100ms

---

## Input Parsing Pattern
Format: `[target] (button1) (button2) ... {joltage}`
- `[.##.]` - Target pattern (`.` = off, `#` = on) for Part 1
- `(0,2,3)` - Button affects indices 0, 2, 3
- `{3,5,4,7}` - Joltage requirements for Part 2

---

## Unit Test Design Rationale

### Parser Tests (`TestParseMachine`)
**Purpose**: Verify correct extraction of all input components

| Test Case | Why It Matters |
|-----------|----------------|
| `simple_pattern_with_multiple_buttons` | Main happy path - validates the full format |
| `pattern_with_dots_only` | Edge case: all lights off target |
| `pattern_with_hashes_only` | Edge case: all lights on target |
| `five_light_pattern`, `six_light_pattern` | Different sizes work correctly |
| `missing_brackets` | Error handling for malformed input |
| `empty_input` | Guard against empty string crash |

### Simulation Tests (`TestSimulatePresses`)
**Purpose**: Verify the core toggle/increment logic

| Test Case | Why It Matters |
|-----------|----------------|
| `no_presses` | Baseline: unpressed buttons have no effect |
| `single_button_press` | Basic functionality works |
| `overlapping_toggles` | **Critical**: Two buttons toggling same light cancel out |
| `toggle_same_light_multiple_times` | Odd presses = on, even presses = off |
| `example_from_problem` | Validates against known correct behavior |

The `overlapping_toggles` test is especially important because it verifies the XOR property that makes Part 1 solvable.

### Algorithm Tests (`TestFindMinPresses`, `TestFindMinPressesPart2`)
**Purpose**: Verify optimization finds true minimum

| Test Case | Why It Matters |
|-----------|----------------|
| Problem examples (1, 2, 3) | Known answers from problem statement |
| `all_lights_off_target` | Zero presses is valid solution |
| `single_button_single_counter` | Simplest non-trivial case |
| `no_buttons_*` | Edge cases for empty button list |
| `impossible` cases | Algorithm correctly reports no solution |

### Integration Tests (`TestSolvePart1_ExampleFile`, `TestSolvePart2_ExampleFile`)
**Purpose**: End-to-end validation

These tests create temporary files with example input and verify the complete pipeline (parse → solve → sum) produces the expected answer from the problem statement.

### Why These Specific Tests?
1. **Example-based**: Problem statements provide examples with known answers - these are gold-standard test cases
2. **Edge cases**: Empty inputs, impossible configurations, boundary conditions
3. **Property-based**: Tests that verify mathematical properties (e.g., toggling twice = identity)
4. **Regression prevention**: If optimization changes break something, tests catch it

---

## References for Further Reading

### Linear Algebra over Finite Fields (GF(2))
- [Lights Out Puzzle - Wikipedia](https://en.wikipedia.org/wiki/Lights_Out_(game)) - Classic example of GF(2) linear algebra
- [Gaussian Elimination over GF(2)](https://codeforces.com/blog/entry/68953) - Codeforces tutorial on XOR basis and Gaussian elimination

### Integer Linear Programming
- [Introduction to Linear Programming](https://www.math.ucla.edu/~tom/LP.pdf) - UCLA course notes
- [Simplex Method](https://en.wikipedia.org/wiki/Simplex_algorithm) - General LP solving
- [Integer Programming](https://en.wikipedia.org/wiki/Integer_programming) - When variables must be integers

### Gaussian Elimination
- [Gaussian Elimination - Wikipedia](https://en.wikipedia.org/wiki/Gaussian_elimination) - Standard algorithm
- [Numerical Stability in Gaussian Elimination](https://nhigham.com/2021/04/20/what-is-gaussian-elimination/) - Why partial pivoting matters

### Related Competitive Programming Problems
- [Codeforces - XOR Basis](https://codeforces.com/problemset/problem/895/C) - Practice with GF(2)
- [LeetCode - Bulb Switcher](https://leetcode.com/problems/bulb-switcher/) - Simpler toggle problem

---

## Gotchas
1. Button indices in parentheses may affect counters outside the target range - handle gracefully
2. Empty lines in input should be skipped
3. The `{joltage}` section comes after buttons - don't stop parsing buttons too early
4. **Part 2 pitfall**: Initial brute-force approach (iterating all totals) is O(maxTarget * n^maxTarget) - way too slow!
5. **Floating point**: Gaussian elimination with floats can have precision issues; use tolerance (1e-9) for comparisons
