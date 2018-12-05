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
	reader.Scan()
	input := reader.Bytes()
	origL := len(input)
	scratch := make([]byte, origL)
	copy(scratch, input)
	reduced := reduce(scratch)
	fmt.Printf("part 1: input: %d %d\n", len(input), len(reduced))
	shortest := len(reduced)
	for c := 'a'; c <= 'z'; c++ {
		copy(scratch, input)
		reduced = removeChar(scratch, c)
		reduced = reduce(reduced)
		if len(reduced) < shortest {
			shortest = len(reduced)
		}
	}
	fmt.Printf("shortest: %d\n", shortest)
}

// modifies input bytes
func removeChar(input []byte, lower rune) []byte {
	lc := byte(lower)
	if lc <= 'Z' {
		lc -= 32
	}
	uc := lc - 32
	i := 0
	for i < len(input) {
		switch input[i] {
		case lc, uc:
			// squanch it
			copy(input[i:], input[i+1:])
			input = input[:len(input)-1]
			continue
		}
		i++
	}
	return input
}

// modifies input bytes
func reduce(input []byte) []byte {
	i := 0
	for i < len(input)-1 {
		if abs(int(input[i])-int(input[i+1])) == 32 {
			copy(input[i:], input[i+2:])
			input = input[:len(input)-2]
			if i > 0 {
				i--
			}
			continue
		}
		i++
	}
	return input
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
