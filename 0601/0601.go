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

type direction int

const (
	right direction = iota
	up
	left
	down
)

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

	for x := range grid {
		grid[x] = make([]*point, maxY)
		for y := 0; y < maxY; y++ {
			p := ClosestToMe(&point{x, y}, data)
			grid[x][y] = p
			areaSums[p]++
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

func FindArea(grid [][]*point, us *point, notus []*point) int {
	var xPlus, yPlus, xMinus, yMinus *point

	if xPlus, notus = FindClosest(us, notus, right); xPlus == nil {
		fmt.Println("infinite right")
		return -1
	}
	if yPlus, notus = FindClosest(us, notus, up); yPlus == nil {
		fmt.Println("infinite up")
		return -1
	}
	if xMinus, notus = FindClosest(us, notus, left); xMinus == nil {
		fmt.Println("infinite left")
		return -1
	}
	if yMinus, notus = FindClosest(us, notus, down); yMinus == nil {
		fmt.Println("infinite down")
		return -1
	}
	area := 0
	for x := xMinus.x; x < xPlus.x; x++ {
		for y := yMinus.y; y < yPlus.y; y++ {
			// fmt.Printf("trying grid[%d][%d]\n", x, y)
			if grid[x][y] == us {
				area++
			}
		}
	}

	return area
}

func FindClosest(p *point, data []*point, d direction) (*point, []*point) {
	var closest *point
	var distance int
	closestIndex := -1
	for i := range data {
		matched := false
		switch d {
		case right:
			matched = (data[i].x > p.x && (closest == nil || Distance(data[i], p) <= distance))
		case up:
			matched = (data[i].y > p.y && (closest == nil || Distance(data[i], p) <= distance))
		case left:
			matched = (data[i].x < p.x && (closest == nil || Distance(data[i], p) <= distance))
		case down:
			matched = (data[i].y < p.y && (closest == nil || Distance(data[i], p) <= distance))
		}
		if matched {
			distance = Distance(data[i], p)
			closest = data[i]
			closestIndex = i
		}
	}

	if closest != nil {
		data = Exclude(closestIndex, data)
	}
	return closest, data
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
