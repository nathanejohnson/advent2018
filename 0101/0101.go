package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	fh, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	reader := bufio.NewScanner(fh)
	total := 0
	for i := 0; reader.Scan(); i++ {
		t := reader.Text()
		nval, err := strconv.Atoi(t)
		if err != nil {
			log.Print("%d: error %s", i, err)
			continue
		}
		total += nval
	}

	fmt.Printf("total: %d\n", total)
}
