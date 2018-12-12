package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"
)

func main() {

	input := flag.Int("serial", 8199, "super serial")

	grid := make([][]int, 300)
	integral := make([][]int, 300)
	tot := 0
	for x := range grid {
		grid[x] = make([]int, 300)
		integral[x] = make([]int, 300)
		for y := range grid[x] {
			p := powerLevel(x, y, *input)
			tot += p
			grid[x][y] = p
			s := p
			if y > 0 {
				s += integral[x][y-1]
			}
			if x > 0 {
				s += integral[x-1][y]
			}
			if y > 0 && x > 0 {
				s -= integral[x-1][y-1]
			}
			integral[x][y] = s
		}
	}
	t := time.Now()

	fmt.Printf("total score: %d\n", tot)
	var maxPowah *int
	var maxX, maxY, maxSize int
	downCt := 0
	for size := 4; size < 300; size++ {
		var thisMaxPowah *int
		for x := 0; x < len(grid)-size; x++ {
			for y := 0; y < len(grid[x])-size; y++ {
				right := x + size - 1
				top := y + size - 1
				a := integral[x][y]
				b := integral[right][y]
				c := integral[x][top]
				d := integral[right][top]
				totalPowah := a + d - b - c

				checkPowah := 0
				for i := 0; i < x+size; i++ {
					for j := 0; j < y+size; j++ {
						checkPowah += grid[i][j]
					}
				}

				if totalPowah != checkPowah {
					fmt.Printf("fuck: x: %d, y: %d, size: %d, tp: %d, cp: %d, a: %d b: %d c: %d d: %d\n",
						x, y, size, totalPowah, checkPowah, a, b, c, d)
				} else {
					fmt.Printf("match!\n")
				}

				if maxPowah == nil || totalPowah > *maxPowah {
					maxPowah = &totalPowah
					maxX = x + 2
					maxY = y + 2
					maxSize = size
				}
				if thisMaxPowah == nil || totalPowah > *thisMaxPowah {
					thisMaxPowah = &totalPowah
				}
			}
		}
		if *thisMaxPowah < *maxPowah {
			downCt++
		}
		if downCt == 3 {
			break
		}
	}
	elapsed := time.Since(t)
	fmt.Printf("part2: %d,%d,%d - %s\n", maxX, maxY, maxSize, elapsed.String())
}

func powerLevel(x, y, serial int) int {
	rackID := x + 11
	powah := rackID * (y + 1)
	powah += serial
	powah *= rackID
	powah = hundo(powah)
	powah -= 5
	return powah
}

func hundo(i int) int {
	if i < 100 {
		return 0
	}
	d := strconv.Itoa(i)
	hundo, _ := strconv.Atoi(string(d[len(d)-3]))
	return hundo
}
