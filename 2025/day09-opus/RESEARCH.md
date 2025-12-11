# RESEARCH by GPT 5.1

Here’s how I’d attack this as an actual production-ish Go solution with 16 cores, plus the reasoning for *why* this is a good approach.

I’ll go in this order:

1. Model the problem mathematically (Part 1 vs Part 2).
2. Part 1: O(N²) algorithm + parallelization in Go.
3. Part 2: how the red+green tiles form a rectilinear polygon.
4. Data structures to represent the “valid region” efficiently.
5. Parallel max-rectangle search for Part 2 (using those structures).
6. Complexity and performance justification.

---

## 1. Problem model

You are given N red tiles with integer coordinates:

```text
R = { (x_i, y_i) | i = 0..N-1 }
```

### Part 1

You may pick any two red tiles as opposite corners of an axis-aligned rectangle. AoC’s geometry is inclusive: the number of tiles horizontally between x1 and x2 (including both ends) is `|x1 - x2| + 1`; similarly for y.

So for two red tiles `p=(x1,y1)` and `q=(x2,y2)`:

```text
width  = |x1 - x2| + 1
height = |y1 - y2| + 1
area   = width * height
```

Even “thin” rectangles (height or width = 1) are allowed.

There’s no restriction on the interior for Part 1; every pair is valid. So:

> Part 1: maximize area over all unordered pairs of red tiles.

This is symmetric: `(i,j)` and `(j,i)` are the same rectangle, so we only need `i < j`.

### Part 2

Now there is a mask of allowed tiles (red + green):

* Red tiles are your given inputs.
* Between `R[i]` and `R[i+1]` (and `R[N-1]` to `R[0]`), there is a straight horizontal or vertical chain of **green** tiles.
* Those edges form a **simple rectilinear polygon** (axis-aligned, no self-crossing).
* All tiles *inside* that loop are also green.
* Everything outside this polygon is forbidden.

So:

> Part 2: same set of candidate rectangles (two red corners),
> but we only accept those whose entire area is contained in the polygon
> (every tile inside is red or green).

We need an efficient way to answer:

> Is the axis-aligned closed rectangle `[x_min..x_max] × [y_min..y_max]`
> fully contained in this rectilinear polygon?

We must do this for many rectangles (O(N²) candidates), so membership checks must be **fast**.

---

## 2. Part 1: O(N²) + concurrency in Go

### Core algorithm

Brute force is straightforward:

```text
maxArea = 0
for i in [0..N-1]:
  for j in [i+1..N-1]:
    dx = abs(x[i] - x[j]) + 1
    dy = abs(y[i] - y[j]) + 1
    area = dx * dy
    if area > maxArea: maxArea = area
```

Time complexity: `O(N²)` arithmetic operations.

Is that acceptable? For AoC-like inputs:

* N is typically a few thousand at most.
* For N = 5000, pairs ≈ 12.5M.
* At ~5–10 ns per simple pair evaluation in optimized Go on a modern CPU,
  we’re in the low tens of milliseconds. Even with overhead, we’re well below 1 second.
* With 16 cores we can divide the work and drive it down further.

There isn’t a “smarter” algorithm that beats O(N²) *and* is simple and safe here. You’re trying to maximize a function over all pairs; with no monotonicity or convexity guarantees beyond the trivial box bounds, any sub-quadratic algorithm is complex and brittle compared to the N² brute, and unnecessary at AoC scales.

### Data types

The sample full input answer is `4,771,508,457`, which is > 2³¹−1. So:

* Coordinates should be `int64`.
* Areas should be `int64`.
* Any `int` use must be carefully limited to indices.

### Parallelization strategy

We want to minimize overhead and false sharing.

Good pattern:

* Precompute `points []Point`.
* Determine number of workers = `min(16, N)` (no point in more than N).
* Split the outer loop index `i` into contiguous chunks:

  * Worker `w` handles `i ∈ [start_w, end_w)`.
  * For each `i`, inner loop still runs `j ∈ [i+1..N-1]`.
* Each worker keeps its own `localMax int64` in **stack** (no shared writes).
* After all workers are done, we reduce `max(all localMax)` in a single-threaded pass.

This avoids:

* Channels inside the hot loop.
* Atomic operations inside the hot loop.
* Contention on shared variables.

#### Go code sketch (Part 1)

```go
package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
)

type Point struct {
	X, Y int64
}

func abs64(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func maxRectAreaPart1(points []Point, maxWorkers int) int64 {
	n := len(points)
	if n < 2 {
		return 0
	}

	// Limit workers reasonably.
	workers := runtime.NumCPU()
	if maxWorkers > 0 && maxWorkers < workers {
		workers = maxWorkers
	}
	if workers > n {
		workers = n
	}

	chunk := (n + workers - 1) / workers
	maxes := make([]int64, workers)

	done := make(chan struct{}, workers)

	for w := 0; w < workers; w++ {
		start := w * chunk
		end := start + chunk
		if start >= n {
			break
		}
		if end > n {
			end = n
		}

		go func(idx, start, end int) {
			var localMax int64
			for i := start; i < end; i++ {
				pi := points[i]
				for j := i + 1; j < n; j++ {
					pj := points[j]
					dx := abs64(pi.X-pj.X) + 1
					dy := abs64(pi.Y-pj.Y) + 1
					area := dx * dy
					if area > localMax {
						localMax = area
					}
				}
			}
			maxes[idx] = localMax
			done <- struct{}{}
		}(w, start, end)
	}

	var globalMax int64
	for w := 0; w < workers; w++ {
		<-done
		if maxes[w] > globalMax {
			globalMax = maxes[w]
		}
	}
	return globalMax
}

func main() {
	// Very barebones input parse, for completeness.
	in := bufio.NewScanner(os.Stdin)
	points := make([]Point, 0, 4096)
	for in.Scan() {
		var x, y int64
		_, err := fmt.Sscanf(in.Text(), "%d,%d", &x, &y)
		if err != nil {
			continue
		}
		points = append(points, Point{X: x, Y: y})
	}
	fmt.Println("Part 1:", maxRectAreaPart1(points, 16))
}
```

This is essentially optimal for Part 1 at AoC scales:

* O(N²) work, but with 16-way parallelization.
* No lock contention in the inner loop.
* Straight-line arithmetic.

---

## 3. Understanding Part 2 geometry: polygon of valid tiles

From the Part 2 description:

* The red points, in input order, form vertices of a rectilinear polygon.
* Between each `R[i]` and `R[i+1]` (and last to first) there is a straight horizontal or vertical segment of green tiles.
* All tiles *inside* this closed loop are also green.
* Outside is invalid.

So we can treat:

* The union “red ∪ green” tiles as the **interior + boundary** of a simple rectilinear polygon `P`.
* What we need: find the largest area rectangle whose tile grid is fully inside `P`, and whose two opposite corners are in `R`.

So the core Part 2 problem:

> Given a simple rectilinear polygon P and a set of its vertices R, find a maximum area axis-aligned rectangle with two corners in R that is fully contained in P.

We must precompute some representation of P that lets us answer:
“Is rectangle R fully inside P?” in O(log N) or O(1).

---

## 4. Best data structures for Part 2

### High-level idea

We cannot afford to enumerate all tiles (coordinates may be very large). Instead, we exploit:

* The polygon is axis-aligned.
* Its boundary only changes at the y-coordinates of vertices.
* Between two consecutive distinct y-coordinates, the vertical cross-section of the polygon is **topologically constant**.

This suggests a **scanline / stripe** representation:

1. Extract all unique y-coordinates from red points: `Y = sorted unique { y_i }`.

2. Consider the bands (strips) between them:

   ```text
   band k: (Y[k], Y[k+1])   for k = 0..len(Y)-2
   ```

   For all tile rows in that open interval, the polygon’s horizontal valid intervals are the same (because you don’t cross a horizontal edge until you hit another Y[k]).

3. For each band, we compute the set of **x-intervals** that are inside P on that band:

   * Each vertical boundary segment contributes an x coordinate where “inside/outside” flips.
   * Sort those x’s and pair them up: `[x0,x1], [x2,x3], ...` are the filled intervals.

We also need to treat **rows that sit exactly on horizontal edges** (y equal to some vertex). Tiles along the edges are green by definition, so they are valid. Practically:

* You can treat top/bottom bands as including the boundary, or define special bands at these exact y-levels.
* As long as we are consistent and never mix tiles on opposing sides of a boundary, the algorithm is sound.

For AoC scale, the simplest robust choice is:

* Work with **scanlines at y + 0.5** (half-integer positions between integer rows).
* That is standard in point-in-polygon parity tests.
* A tile whose center is at (x+0.5, y+0.5) lies in the band between y and y+1.
* So if we know the intervals for band `(Y[k], Y[k+1])`, we know exactly which tiles in those rows are valid.

### Step 1: Build vertical segments

Given the red points in order:

```go
type Seg struct {
	X        int64
	Y1, Y2   int64 // Y1 < Y2
}
```

For each consecutive pair `p = R[i]` and `q = R[(i+1)%N]`:

* If `p.X == q.X` → vertical segment:

  ```go
  y1, y2 := p.Y, q.Y
  if y1 > y2 {
      y1, y2 = y2, y1
  }
  segs = append(segs, Seg{X: p.X, Y1: y1, Y2: y2})
  ```

* If `p.Y == q.Y` → horizontal edge; it doesn’t affect scanline intersections directly but defines bands – we already account via unique Y.

We’ll have O(N) segments.

### Step 2: compute Y bands

```go
ys := make([]int64, 0, len(points))
for _, p := range points {
    ys = append(ys, p.Y)
}
sort.Slice(ys, func(i, j int) bool { return ys[i] < ys[j] })
ys = slices.Compact(ys) // Go 1.21+, or manual dedup
// We'll have len(ys) = M
```

Bands are between `ys[k]` and `ys[k+1]`.

A scanline for band k is, e.g.:

```go
scanY := (ys[k] + ys[k+1]) / 2.0
```

But we never need float64; for each segment we just determine whether it intersects the open interval `(ys[k], ys[k+1])`, which is equivalent to:

```text
y1 < ys[k+1] AND y2 > ys[k]
```

If true, the segment crosses every scanline in that band.

### Step 3: per-band x-intervals

For each band k:

1. Collect all `X` of vertical segments crossing that band:

   ```go
   xs := []int64{}
   for each seg in segs:
       if seg.Y1 < ys[k+1] && seg.Y2 > ys[k]:
           xs = append(xs, seg.X)
   ```

2. Sort xs ascending.

3. Because polygon is simple, the intersection with a horizontal line is an even number of points; inside/outside alternates. So:

   * For `xs[0], xs[1]` → inside interval `[xs[0], xs[1]]`
   * For `xs[2], xs[3]` → next inside interval, etc.

4. Store for band k:

```go
type Interval struct {
    Lo, Hi int64 // inclusive, original coords
}

type Stripe struct {
    YLo, YHi  int64      // the band (ys[k], ys[k+1])
    Inside    []Interval // sorted by Lo, disjoint
}
```

We can build `[]Stripe` once, sequentially (O(N²) worst-case if you naively test all segments per band). But:

* N is small (AoC scale).
* There are at most O(N) segments and O(N) bands, so N² ~ few million ops, which is absolutely fine.
* For extra polish, you can optimize using event lists per band (sweep-line), but not necessary.

Now we have a compact representation of “for every vertical band, which horizontal intervals are inside the polygon”.

### Step 4: fast membership check for a rectangle

Given a candidate rectangle with corners `(x1,y1)` and `(x2,y2)`:

1. Normalize:

   ```go
   xLo, xHi := min(x1,x2), max(x1,x2)
   yLo, yHi := min(y1,y2), max(y1,y2)
   ```

2. The rectangle covers all tile rows between yLo and yHi inclusive. That corresponds to all bands whose `(YLo,YHi)` intersects `[yLo, yHi]` at tile-centre level.

   In practice:

   * Find band index of the lowest scan band whose upper bound > yLo.
   * Continue until band’s lower bound >= yHi.

   Roughly:

   ```go
   // Prebuild slice of stripes sorted by YLo
   // Use binary search to find starting band
   idx := lowerBoundOnStripe( stripes, yLo )
   for s := idx; s < len(stripes) && stripes[s].YLo < yHi; s++ {
       // check stripe s
   }
   ```

3. For each stripe s in this range, we need to ensure that every tile in the rectangle width `[xLo..xHi]` is inside P. Since we precomputed `Inside []Interval` for that band, we just need to know:

   > Does there exist some interval [Lo,Hi] in stripe s such that Lo ≤ xLo and Hi ≥ xHi?

   That is a simple binary search on intervals (they’re sorted by `Lo`):

   ```go
   func stripeCovers(stripe Stripe, xLo, xHi int64) bool {
       intervals := stripe.Inside
       // binary search to find rightmost interval with Lo <= xLo
       i := sort.Search(len(intervals), func(i int) bool {
           return intervals[i].Lo > xLo
       }) - 1
       if i < 0 {
           return false
       }
       return intervals[i].Hi >= xHi
   }
   ```

   If any stripe fails this `stripeCovers`, the rectangle is invalid.

So the per-rectangle check is:

* O(log(#intervals_per_stripe)) per stripe for the binary search (tiny).
* Times O(#stripes_intersected) which is at most O(#unique Y) for worst-case tall rectangles.

Given N is small, this is perfectly manageable.

At AoC scale, average number of stripes per rectangle is usually much smaller, because the polygon is not tall and skinny across all coordinates.

---

## 5. Parallel max-rectangle search for Part 2

The outer shape is identical to Part 1, but with an extra validity check.

### Steps

1. Build `points []Point` from input.

2. Build:

   * `ys []int64` (unique y’s).
   * `segs []Seg` (vertical segments).
   * `stripes []Stripe` (as above).

3. Implement `isRectInside(x1,y1,x2,y2)` using the stripe structure.

4. Parallel outer loop over `i` exactly as Part 1.

5. For each `i`, inner loop:

   ```go
   for j := i+1; j < n; j++ {
       pi, pj := points[i], points[j]
       if !isRectInside(pi.X, pi.Y, pj.X, pj.Y) {
           continue
       }
       dx := abs64(pi.X-pj.X) + 1
       dy := abs64(pi.Y-pj.Y) + 1
       area := dx * dy
       if area > localMax {
           localMax = area
       }
   }
   ```

Because `stripes` are read-only once constructed, the membership function is safe to call from multiple goroutines without any locks.

### Go skeleton for Part 2 “worker”

```go
func isRectInside(stripes []Stripe, ys []int64, x1, y1, x2, y2 int64) bool {
    if x1 == x2 && y1 == y2 {
        // single tile; if it's red it must be inside
        return true
    }
    if y1 == y2 || x1 == x2 {
        // thin rectangle is still fine, handled like generic case
    }

    xLo, xHi := x1, x2
    if xLo > xHi {
        xLo, xHi = xHi, xLo
    }
    yLo, yHi := y1, y2
    if yLo > yHi {
        yLo, yHi = yHi, yLo
    }

    // find first stripe that potentially intersects [yLo, yHi)
    // e.g. stripe.YHi > yLo
    // stripes are sorted by YLo
    idx := sort.Search(len(stripes), func(i int) bool {
        return stripes[i].YHi > yLo
    })

    for s := idx; s < len(stripes); s++ {
        stripe := stripes[s]
        if stripe.YLo >= yHi {
            break
        }
        if !stripeCovers(stripe, xLo, xHi) {
            return false
        }
    }
    return true
}

func maxRectAreaPart2(points []Point, stripes []Stripe, maxWorkers int) int64 {
    n := len(points)
    if n < 2 {
        return 0
    }

    workers := runtime.NumCPU()
    if maxWorkers > 0 && maxWorkers < workers {
        workers = maxWorkers
    }
    if workers > n {
        workers = n
    }

    chunk := (n + workers - 1) / workers
    maxes := make([]int64, workers)
    done := make(chan struct{}, workers)

    for w := 0; w < workers; w++ {
        start := w * chunk
        end := start + chunk
        if start >= n {
            break
        }
        if end > n {
            end = n
        }
        go func(idx, start, end int) {
            var localMax int64
            for i := start; i < end; i++ {
                pi := points[i]
                for j := i + 1; j < n; j++ {
                    pj := points[j]
                    if !isRectInside(stripes, nil, pi.X, pi.Y, pj.X, pj.Y) {
                        continue
                    }
                    dx := abs64(pi.X-pj.X) + 1
                    dy := abs64(pi.Y-pj.Y) + 1
                    area := dx * dy
                    if area > localMax {
                        localMax = area
                    }
                }
            }
            maxes[idx] = localMax
            done <- struct{}{}
        }(w, start, end)
    }

    var globalMax int64
    for w := 0; w < workers; w++ {
        <-done
        if maxes[w] > globalMax {
            globalMax = maxes[w]
        }
    }
    return globalMax
}
```

I’ve omitted the `stripeCovers` and stripe-building code for brevity, but the core structure is clear:

* Precomputation of polygon → stripes done once.
* Hot loop is pair enumeration with a pure read-only predicate.

---

## 6. Why this is a good algorithm (evidence / justification)

### Correctness reasoning

* Part 1: trivial; we directly compute the intended formula over all candidate pairs.
* Part 2:

  * The red vertices and connecting segments define the boundary of a simple rectilinear polygon P.
  * On any fixed horizontal line y (or equivalently, any band between vertex y’s), the intersection of P with that line is a union of closed segments. The “even-odd” rule for polygon filling tells us that between each pair of vertical crossings we alternate outside/inside.
  * Our stripe construction exactly captures these intervals: sorted intersection x’s, paired up.
  * By checking that `[xLo,xHi]` is contained in one of those intervals for every band intersecting our rectangle’s y-range, we guarantee that every tile center in the rectangle is inside or on the boundary of P.
  * Since boundary tiles are green and interior tiles are green, the rectangle is valid iff all stripes covering its rows validate.

This is standard scanline / even-odd fill logic for rectilinear polygons.

### Complexity evidence

Let:

* N = number of red tiles (vertices).
* S = number of vertical segments (≤ N).
* B = number of bands (≤ N-1).
* I = average number of intervals per band (typically small, often 1–3).

Costs:

1. Stripe building (naive version):

   ```text
   For each band: inspect each segment
   O(B * S) = O(N²)
   ```

   For N in the low thousands, this is at most a few million operations.

   * Example: N=4000 → ~16M segment checks.
   * Each check is a couple of integer comparisons; order of tens of ms.

   If you really want, you can reduce this to O(N log N) with a sweep-line, but it’s not necessary.

2. Pair enumeration:

   * Pairs ≈ N(N-1)/2 ≈ O(N²).
   * For each pair, you do:

     * A handful of arithmetic ops.
     * A few binary searches over tiny slices (intervals per stripe).
     * Iterate over the bands between yLo & yHi.

   For random polygons with moderate complexity, average stripes per rectangle is much smaller than B. Worst case, B ≈ N, so O(N) per rectangle → O(N³) worst-case. But AoC inputs are not adversarial, and the polygon loop tends to be “fat”, not a deep comb of narrow fingers. In practice, this is much closer to O(N²) with a small constant.

3. Concurrency:

   * With 16 CPUs, pair enumeration distributes almost linearly in the outer loop.
   * If single-threaded runtime is `T`, 16-way parallelization gives near `T/16` in practice, minus scheduling overhead.
   * With N in the 2–4K range, you’re very comfortably inside a second even with the extra Part 2 checks.

### Memory footprint

* `points []Point` is O(N).
* `segs []Seg` is O(N).
* `stripes []Stripe`: each stripe holds a small number of intervals. Total intervals across all stripes are O(N).
* No huge 2D arrays; we never materialize the polygon area as a matrix.

So memory is O(N), dominated by the input itself.

### Why not a 2D grid + prefix sums?

Alternate idea:

* Coordinate compress x and y.
* Flood fill the polygon interior in a compressed grid.
* Build 2D prefix sums for “allowed” cells.
* Check rectangles via area equality.

This works, but:

* In worst case, compressed grid size is O(N²) cells.
* 2D prefix sums require O(N²) memory and time.
* It scales less gracefully than the stripe representation.
* You also lose the very nice “only O(N) intervals across all stripes” guarantee.

Given that your polygon is rectilinear and you only need axis-aligned rectangles, the stripe-interval representation is strictly better in asymptotic memory and at least as good in speed.

---

If you’d like, I can next:

* Flesh out the exact stripe-building code, including careful handling of band boundaries.
* Or: adapt this into a fully compilable AoC-ready `main.go` with CLI flags to run part 1/2 separately and a simple benchmark harness to check scaling across cores.

