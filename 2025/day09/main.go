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
	fmt.Println("Day 09: Movie Theater - Finding Largest Rectangle")
	result := run("input.txt")
	fmt.Printf("Answer: %d\n", result)
}

func run(inputFile string) int {
	points := parseInput(inputFile)
	return findLargestRectangle(points)
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
