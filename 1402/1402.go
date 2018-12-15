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
	bigLine := make([]byte, 2, capSize)

	elf1, elf2 := 0, 1
	bigLine[elf1], bigLine[elf2] = 3, 7

	sequence := []byte(input)
	for i := range sequence {
		sequence[i] -= '0'
	}
	t := time.Now()
OUTER:
	for {
		for _, d := range digitSlice(uint(bigLine[elf1] + bigLine[elf2])) {
			bigLine = append(bigLine, d)
			if bytes.HasSuffix(bigLine, sequence) {
				break OUTER
			}
		}
		elf1 += int(bigLine[elf1]) + 1
		elf1 %= len(bigLine)
		elf2 += int(bigLine[elf2]) + 1
		elf2 %= len(bigLine)
	}
	elapsed := time.Since(t)
	fmt.Printf("%s first appears after %d recipes : %s\n", input, len(bigLine)-len(sequence), elapsed)

}

func digitSlice(ui uint) []byte {
	if ui < 10 {
		return []byte{byte(ui)}
	}
	return []byte{byte(ui / 10), byte(ui % 10)}
}
