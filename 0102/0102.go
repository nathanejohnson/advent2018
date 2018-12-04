package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {

	totals := make(map[int]bool)
	total := 0
	fh, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
outside:
	for outer := 0; ; outer++ {
		reader := bufio.NewScanner(fh)
		for i := 0; reader.Scan(); i++ {
			t := reader.Text()
			nval, err := strconv.Atoi(t)
			if err != nil {
				log.Print("%d: error %s", i, err)
				continue
			}
			total += nval
			if totals[total] {
				fmt.Printf("found total %d on outer %d inner %d\n", total, outer, i)
				break outside
			}
			totals[total] = true
		}
		_, err := fh.Seek(0, 0)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("total: %d\n", total)
}
