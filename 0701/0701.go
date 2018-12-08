package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

func main() {

	fh, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	reader := bufio.NewScanner(fh)

	depList := make(map[string][]string)
	var empty []string
	for reader.Scan() {
		var a, b string
		_, err := fmt.Sscanf(reader.Text(), "Step %s must be finished before step %s can begin.", &a, &b)
		if err != nil {
			panic(err)
		}
		depList[b] = append(depList[b], a)
		depList[a] = append(depList[a], empty...)
	}
	fmt.Printf("data len: %d\n", len(depList))
	var output string
	for i := 0; ; i++ {
		dones := readies(depList)
		sort.Strings(dones)
		if len(dones) == 0 {
			break
		}

		v := dones[0]
		output += v
		dones = dones[1:]
		remove(v, depList)
	}

	fmt.Printf("part1: %s , len: %d\n", output, len(output))

}

func readies(data map[string][]string) []string {
	var out []string
	for k, v := range data {
		if len(v) == 0 {
			out = append(out, k)
		}
	}
	return out
}

func remove(key string, data map[string][]string) {
	for k, v := range data {
		switch k {
		case key:
			delete(data, k)
		default:
			var newV []string
			for _, s := range v {
				if s != key {
					newV = append(newV, s)
				}
			}
			data[k] = newV
		}
	}
}
