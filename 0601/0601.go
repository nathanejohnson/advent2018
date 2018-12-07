package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type point struct {
	x int
	y int
}

func main() {
	var data []*point

	fh, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	reader := bufio.NewScanner(fh)

	maxX := 0
	maxY := 0
	for reader.Scan() {
		p := &point{}
		_, err := fmt.Sscanf(reader.Text(), "%d,%d", &p.x, &p.y)
		if err != nil {
			panic(err)
		}
		data = append(data, p)
		if p.x > maxX {
			maxX = p.x
		}
		if p.y > maxY {
			maxY = p.y
		}
	}

	grid := make([][]*point, maxX)
	areaSums := make(map[*point]int)
	theTenK := 0
	for x := range grid {
		grid[x] = make([]*point, maxY)
		for y := 0; y < maxY; y++ {
			me := &point{x, y}
			p := ClosestToMe(me, data)
			grid[x][y] = p
			areaSums[p]++
			theTenK += AccDistance10K(me, data)
		}

	}

	infinitePoints := make(map[*point]bool)
	for x := 0; x < maxX; x++ {
		infinitePoints[grid[x][0]] = true
		infinitePoints[grid[x][maxY-1]] = true
	}

	for y := 0; y < maxY; y++ {
		infinitePoints[grid[0][y]] = true
		infinitePoints[grid[maxX-1][y]] = true
	}

	maxArea := 0
	for i := range data {
		if infinitePoints[data[i]] {
			continue
		}
		if area := areaSums[data[i]]; area > maxArea {
			maxArea = area
		}
	}
	fmt.Printf("part 1: %d\n", maxArea)
	fmt.Printf("part 2: %d\n", theTenK)
}

func AccDistance10K(me *point, data []*point) int {
	agg := 0
	for i := range data {
		agg += Distance(me, data[i])
		if agg >= 10000 {
			return 0
		}
	}
	return 1
}

func ClosestToMe(me *point, points []*point) *point {
	distance := -1
	dct := 0
	var out *point
	for i := range points {
		d := Distance(me, points[i])
		switch {
		case d == 0:
			return points[i]
		case d == distance:
			dct++
			out = nil
		case distance == -1 || d < distance:
			distance = d
			dct = 0
			out = points[i]
		}
	}
	return out
}

func Exclude(i int, data []*point) []*point {
	notus := make([]*point, 0, len(data)-1)
	if i > 0 {
		notus = append(notus, data[:i]...)
	}
	return append(notus, data[i+1:]...)
}

func Distance(a, b *point) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func abs(i int) int {
	//mask := i >> (unsafe.Sizeof(i)*8 - 1)
	//return (i ^ mask) - mask
	if i < 0 {
		return -i
	}
	return i
}
