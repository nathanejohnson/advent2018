package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"
)

func main() {

	var input string
	var capSize uint
	flag.StringVar(&input, "input", "157901", "puzzle input")
	flag.UintVar(&capSize, "capsize", 25000000, "tweak the slice cap")
	flag.Parse()
	_, err := strconv.ParseUint(input, 10, 32)
	if err != nil {
		log.Fatalf("input is not positive integer: %s", err)
	}
	scores := make([]byte, 2, capSize)

	elf1, elf2 := 0, 1
	scores[elf1], scores[elf2] = 3, 7

	sequence := []byte(input)
	for i := range sequence {
		sequence[i] -= '0'
	}
	t := time.Now()
OUTER:
	for {
		for _, d := range digitSlice(uint(scores[elf1] + scores[elf2])) {
			scores = append(scores, d)
			if bytes.HasSuffix(scores, sequence) {
				break OUTER
			}
		}
		elf1 += int(scores[elf1]) + 1
		elf1 %= len(scores)
		elf2 += int(scores[elf2]) + 1
		elf2 %= len(scores)
	}
	elapsed := time.Since(t)
	fmt.Printf("%s first appears after %d recipes : %s\n", input, len(scores)-len(sequence), elapsed)

}

func digitSlice(ui uint) []byte {
	if ui < 10 {
		return []byte{byte(ui)}
	}
	return []byte{byte(ui / 10), byte(ui % 10)}
}
