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
	for outer := 0; func() bool {
		fh, err := os.Open("./input.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer fh.Close()
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
				return false
			}
			totals[total] = true
		}
		return true
	}(); outer++ {
	}

	fmt.Printf("total: %d\n", total)
}
