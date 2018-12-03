package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	fh, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	reader := bufio.NewScanner(fh)
	twos := 0
	threes := 0
	for reader.Scan() {
		rmap := make(map[rune]int)

		t := reader.Text()
		for _, r := range t {
			rmap[r]++
		}
		fmt.Printf("%s:\t", t)
		var gotTwo, gotThree bool
		for r, cnt := range rmap {
			switch {
			case cnt == 2 && !gotTwo:
				fmt.Printf("2 %s,", string(r))
				twos++
				gotTwo = true
			case cnt == 3 && !gotThree:
				fmt.Printf("3 %s,", string(r))
				threes++
				gotThree = true
			}
		}
		fmt.Println("")
	}

	fmt.Printf("checksum: %d\n", twos*threes)
}
