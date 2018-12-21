package main

import (
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var reader *bufio.Scanner
	var filePath string
	var minutes int
	flag.StringVar(&filePath, "input", "", "input file")
	flag.IntVar(&minutes, "mins", 10, "specify minutes")
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
	var grid [][]byte
	for reader.Scan() {
		buf := reader.Bytes()
		line := make([]byte, len(buf))
		copy(line, buf)
		grid = append(grid, line)
	}

	sums := make(map[string]int)
	repeatCt := 0
	for minute := 0; minute < minutes; minute++ {
		lastGrid := clone(grid)
		for y := range lastGrid {
			for x := range lastGrid[y] {
				around := adjacents(x, y, lastGrid)
				switch lastGrid[y][x] {
				case '.':
					openCt := 0
					for _, a := range around {
						if a == '|' {
							openCt++
						}
					}
					if openCt >= 3 {
						grid[y][x] = '|'
					}
				case '|':
					lumberCt := 0
					for _, a := range around {
						if a == '#' {
							lumberCt++
						}
					}
					if lumberCt >= 3 {
						grid[y][x] = '#'
					}
				case '#':
					treeCt := 0
					lumberCt := 0
					for _, a := range around {
						switch a {
						case '#':
							treeCt++
						case '|':
							lumberCt++
						}
					}
					if treeCt < 1 || lumberCt < 1 {
						grid[y][x] = '.'
					}
				}
			}
		}
		h := hash(grid)
		if offset, ok := sums[h]; ok {
			fmt.Printf("got a repeat at %d for %d\n", minute, offset)
			dt := minute - offset
			repeats := (minutes-offset)/dt - 1
			minute += repeats * dt
			fmt.Printf("dt %d offset %d repeats %d newminute %d\n", dt, offset, repeats, minute)
			repeatCt++
			sums = make(map[string]int)
		} else {
			sums[h] = minute
		}
	}
	yards := 0
	trees := 0
	for y := range grid {
		fmt.Printf("%s\n", grid[y])
		for x := range grid[y] {
			switch grid[y][x] {
			case '|':
				trees++
			case '#':
				yards++
			}
		}
	}
	fmt.Printf("trees %d yards %d combined %d\n", trees, yards, trees*yards)
}

func hash(grid [][]byte) string {
	s := md5.New()
	for y := range grid {
		s.Write(grid[y])
	}
	return base64.RawStdEncoding.EncodeToString(s.Sum(nil))
}

func adjacents(x, y int, grid [][]byte) []byte {
	ret := make([]byte, 0, 8)
	ly := len(grid)
	lx := len(grid[0])

	left := x - 1
	if left < 0 {
		left = 0
	}
	right := x + 1
	if right >= lx {
		right = lx - 1
	}
	top := y - 1
	if top < 0 {
		top = 0
	}
	bottom := y + 1
	if bottom >= ly {
		bottom = ly - 1
	}
	for dy := top; dy <= bottom; dy++ {
		for dx := left; dx <= right; dx++ {
			if dy == y && dx == x {
				continue
			}
			ret = append(ret, grid[dy][dx])
		}
	}
	return ret
}

func clone(in [][]byte) [][]byte {
	out := make([][]byte, len(in))
	for y := range in {
		out[y] = make([]byte, len(in[y]))
		copy(out[y], in[y])
	}
	return out
}
