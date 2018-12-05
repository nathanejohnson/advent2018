package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func main() {
	fh, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	reader := bufio.NewScanner(fh)
	reader.Scan()
	input := reader.Text()
	reduced := reduce(input)
	fmt.Printf("part 1: input: %d %d\n", len(input), len(reduced))
	shortest := len(reduced)
	for c := 'a'; c <= 'z'; c++ {
		re := regexp.MustCompile(fmt.Sprintf("[%s%s]", string(c), string(c-32)))
		newInput := re.ReplaceAllString(input, "")
		fmt.Printf("newInput length: %d , old length: %d\n", len(newInput), len(input))
		l := len(reduce(newInput))
		if l < shortest {
			shortest = l
		}
	}
	fmt.Printf("shortest: %d\n", shortest)
}

func reduce(input string) string {
	i := 0
	j := 1
	for j < len(input) {
		if abs(int(input[i])-int(input[j])) == 32 {
			return reduce(input[:i] + input[j+1:])
		}
		i++
		j++
	}
	return input
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
