package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type point struct {
	x, y int
}

func main() {
	var reader *bufio.Scanner
	var filePath string
	flag.StringVar(&filePath, "input", "", "input file")
	flag.Parse()

	if filePath != "" {
		fh, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		reader = bufio.NewScanner(fh)
		defer fh.Close()
	} else {
		reader = bufio.NewScanner(os.Stdin)
	}
	re := regexp.MustCompile(`^([xy])=(\d+), ([xy])=(\d+)..(\d+)$`)

	xpoints := make(map[int][]int)
	ypoints := make(map[int][]int)

	var minX, minY, maxX, maxY int
	first := true
	for reader.Scan() {
		t := reader.Text()
		matches := re.FindStringSubmatch(t)
		if len(matches) != 6 {
			log.Fatalf("bad input: %s", t)
		}
		digAnchor, _ := strconv.Atoi(matches[2])
		digB, _ := strconv.Atoi(matches[4])
		digE, _ := strconv.Atoi(matches[5])

		if digB > digE {
			log.Fatalf("weird inverted range input: %s", t)
		}

		rng := make([]int, 0, digE-digB+1)

		for i := digB; i <= digE; i++ {
			rng = append(rng, i)
		}
		n := len(rng) - 1
		switch matches[1] {
		case "x":
			xpoints[digAnchor] = append(xpoints[digAnchor], rng...)

			switch {
			case first:
				first = false
				minX = digAnchor
				maxX = digAnchor
				minY = rng[0]
				maxY = rng[n]
			case digAnchor < minX:
				minX = digAnchor
			case digAnchor > maxX:
				maxX = digAnchor
			}
			switch {
			case rng[0] < minY:
				minY = rng[0]
			case rng[n] > maxY:
				maxY = rng[n]
			}
		case "y":
			ypoints[digAnchor] = append(ypoints[digAnchor], rng...)
			switch {
			case first:
				first = false
				minY = digAnchor
				maxY = digAnchor
				minX = rng[0]
				maxX = rng[n]
			case digAnchor < minY:
				minY = digAnchor
			case digAnchor > maxY:
				maxY = digAnchor
			}
			switch {
			case rng[0] < minX:
				minX = rng[0]
			case rng[n] > maxX:
				maxX = rng[n]
			}
		}
	}

	fmt.Printf("minX: %d, maxX: %d, minY: %d, maxY: %d\n", minX, maxX, minY, maxY)
	// dy := maxY - minY
	dx := maxX - minX

	xOffset := minX - 2
	fmt.Printf("xOffset: %d, yOffset: %d\n", xOffset, 0)
	grid := make([][]byte, maxY+3)
	cb := []byte{'.'}
	for g := range grid {
		grid[g] = bytes.Repeat(cb, dx+5)
	}

	grid[0][500-xOffset] = '+'

	for x := range xpoints {
		for _, y := range xpoints[x] {
			xp, yp := x-xOffset, y
			grid[yp][xp] = '#'
		}
	}
	for y := range ypoints {
		for _, x := range ypoints[y] {
			xp, yp := x-xOffset, y
			grid[yp][xp] = '#'
		}
	}
	for y := range grid {
		fmt.Printf("%s\n", grid[y])
	}
	fmt.Println("-----------------------------")
	x := 500 - xOffset
	y := 1
	grid[y][x] = '|'
	dive(x, y, grid)
	fmt.Println("-----------------------------")
	ct := 0
	ctp2 := 0
	for y := range grid {
		fmt.Printf("%s\n", grid[y])
		if y < minY || y > maxY {
			fmt.Printf("skipping %d\n", y)
			continue
		}
		for x := range grid[y] {

			switch grid[y][x] {
			case '~', '|':
				ct++
			}
			if grid[y][x] == '~' {
				ctp2++
			}
		}
	}
	fmt.Printf("part 1 count is %d, part 2: %d\n", ct, ctp2)
}

func dive(x, y int, grid [][]byte) {
	// fmt.Printf("diving with %d, %d\n", x, y)
	maxX := len(grid[0]) - 1
	maxY := len(grid) - 1
	for y < maxY && y >= 0 && x > 0 && x < maxX {
		switch grid[y+1][x] {
		case '.', '|':
			y++
			grid[y][x] = '|'
		case '#', '~':
			bound := 0
			// seek left
			left := x
		LEFTLOOP:
			for ; left > 0; left-- {
				switch grid[y+1][left] {
				case '#', '~':
				default:
					// fmt.Printf("diving left\n")
					if grid[y][left] == '.' {
						grid[y][left] = '|'
					}

					dive(left, y, grid)
					break LEFTLOOP
				}
				switch grid[y][left-1] {
				case '#', '~':
					if grid[y][left] == '.' {
						grid[y][left] = '|'
					}

					bound++
					break LEFTLOOP
				default:
					if grid[y][left] == '.' {
						grid[y][left] = '|'
					}
				}
			}
			// seek right
			right := x
		RIGHTLOOP:
			for ; right < maxX; right++ {
				switch grid[y+1][right] {
				case '#', '~':
				default:
					// fmt.Printf("diving right\n")
					if grid[y][right] == '.' {
						grid[y][right] = '|'
					}
					dive(right, y, grid)
					break RIGHTLOOP
				}
				switch grid[y][right+1] {
				case '#', '~':
					if grid[y][right] == '.' {
						grid[y][right] = '|'
					}
					bound++
					break RIGHTLOOP
				default:
					if grid[y][right] == '.' {
						grid[y][right] = '|'
					}
				}
			}
			if bound == 2 {
				// fmt.Printf("bound is 2!\n")
				for xpr := left; xpr <= right; xpr++ {
					grid[y][xpr] = '~'
				}
				dive(x, y-1, grid)
			}

			return
		}
	}

}
