package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type point struct {
	x, y, vx, vy int
}

func atoi(s string) int {
	d, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return d
}
func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func main() {
	fh, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	reader := bufio.NewScanner(fh)

	re := regexp.MustCompile(`^position=<\s*(-?\d+),\s*(-?\d+)> velocity=<\s*(-?\d+),\s*(-?\d+)>$`)

	var points []*point
	minXSteps := -1
	maxXSteps := -1
	minYSteps := -1
	maxYSteps := -1
	xTot := 0
	yTot := 0
	for reader.Scan() {
		t := reader.Text()
		matches := re.FindStringSubmatch(reader.Text())
		if len(matches) == 0 {
			log.Fatalf("regex failed on %s\n", t)
		}
		p := &point{
			x:  atoi(matches[1]),
			y:  atoi(matches[2]),
			vx: atoi(matches[3]),
			vy: atoi(matches[4]),
		}

		xTot += p.x
		yTot += p.y
		xsteps := abs(p.x / p.vx)
		ysteps := abs(p.y / p.vy)

		points = append(points, p)
		if minXSteps == -1 {
			// initialize
			minXSteps = xsteps
			maxXSteps = xsteps
			minYSteps = ysteps
			minYSteps = ysteps
		} else {
			switch {
			case xsteps < minXSteps:
				minXSteps = xsteps
			case xsteps > maxXSteps:
				maxXSteps = xsteps
			}
			switch {
			case ysteps < minYSteps:
				minYSteps = ysteps
			case ysteps > maxYSteps:
				maxYSteps = ysteps
			}
		}

	}

	fmt.Printf("minXSteps: %d, maxXSteps: %d, minYSteps: %d, maxYSteps:, xTot: %d, yTot: %d %d\n",
		minXSteps, maxXSteps, minYSteps, maxYSteps, xTot, yTot)

	//for i := 0; i < 9000; i++ {
	//	for _, p := range points {
	//		p.x += p.vx
	//		p.y += p.vy
	//
	//	}
	//}

	var lastScore int = -1
	secs := 0
	for {
		s := score(points)
		if lastScore > -1 && s > lastScore {
			for _, p := range points {
				p.x -= p.vx
				p.y -= p.vy

			}
			secs--
			break
		}
		lastScore = s
		for _, p := range points {
			p.x += p.vx
			p.y += p.vy

		}
		secs++

	}
	draw(edges(points))
	fmt.Printf("part2: %d\n", secs)
}

func score(points []*point) int {
	var minX, minY, maxX, maxY *int
	for _, p := range points {
		switch {
		case minX == nil:
			minX = &p.x
			minY = &p.y
			maxX = &p.x
			maxY = &p.y
		case p.x < *minX:
			minX = &p.x
		case p.y < *minY:
			minY = &p.y
		case p.x > *maxX:
			maxX = &p.x
		case p.y > *maxY:
			maxY = &p.y
		}
	}

	return abs(*maxX-*minX) + abs(*maxY-*minY)
}

func edges(points []*point) (grid [][]bool) {
	var minX, maxX, minY, maxY *int
	var xMin, xMax, yMin, yMax int
	hash := make(map[int]map[int]bool)
	for _, p := range points {
		switch {
		case minX == nil:
			minX = &p.x
			minY = &p.y
			maxX = &p.x
			maxY = &p.y
		case p.x < *minX:
			minX = &p.x
		case p.y < *minY:
			minY = &p.y
		case p.x > *maxX:
			maxX = &p.x
		case p.y > *maxY:
			maxY = &p.y
		}
		if _, ok := hash[p.x]; !ok {
			hash[p.x] = make(map[int]bool)
		}
		hash[p.x][p.y] = true
	}

	xMin, xMax, yMin, yMax = *minX, *maxX, *minY, *maxY

	grid = make([][]bool, xMax-xMin+1)
	for x := range grid {
		grid[x] = make([]bool, yMax-yMin+1)
	}
	for x := range hash {
		for y := range hash[x] {
			grid[x-xMin][y-yMin] = true
		}
	}
	return
}

func draw(grid [][]bool) {

	for y := 0; y < len(grid[0]); y++ {
		for x := range grid {
			if grid[x][y] {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print(" |\n")
	}

	for x := 0; x < len(grid); x++ {
		fmt.Print("-")
	}
	fmt.Println()
}
