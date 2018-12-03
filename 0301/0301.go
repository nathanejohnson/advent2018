package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	fh, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer fh.Close()
	reader := bufio.NewScanner(fh)
	re := regexp.MustCompile(`^#(\d+)\s+@\s+(\d+),(\d+):\s+(\d+)x(\d+)$`)
	var grid [1000][1000]int
	for lineNo := 0; reader.Scan(); lineNo++ {
		matches := re.FindStringSubmatch(reader.Text())
		if len(matches) == 6 {
			ox := cvt(matches[2])
			oy := cvt(matches[3])
			dx := cvt(matches[4])
			dy := cvt(matches[5])
			for x := ox; x < ox+dx; x++ {
				for y := oy; y < oy+dy; y++ {
					grid[x][y]++
				}
			}
		}
	}
	overlapInches := 0
	for x := 0; x < 1000; x++ {
		for y := 0; y < 1000; y++ {
			if grid[x][y] > 1 {
				overlapInches++
			}
		}
	}

	fmt.Printf("overlap inches: %d\n", overlapInches)
}

func cvt(s string) int {
	r, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return r
}
