package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type rec struct {
	id uint16
	ox uint16
	oy uint16
	dx uint16
	dy uint16
}

func main() {
	fh, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer fh.Close()
	reader := bufio.NewScanner(fh)
	var grid [1000][1000]uint8
	overlapInches := 0

	for lineNo := 0; reader.Scan(); lineNo++ {
		r := parse(reader.Text())
		for x := r.ox; x < r.ox+r.dx; x++ {
			for y := r.oy; y < r.oy+r.dy; y++ {
				grid[x][y]++
				if grid[x][y] == 2 {
					overlapInches++
				}
			}
		}
	}

	fmt.Printf("overlap inches: %d\n", overlapInches)
}

func parse(line string) rec {
	r := rec{}
	_, err := fmt.Sscanf(line, "#%d @ %d,%d: %dx%d",
		&r.id,
		&r.ox,
		&r.oy,
		&r.dx,
		&r.dy,
	)
	if err != nil {
		panic(fmt.Sprintf("%s %s", line, err))
	}
	return r
}

func cvt(s string) int {
	r, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return r
}
