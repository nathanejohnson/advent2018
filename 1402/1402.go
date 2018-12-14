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
	flag.StringVar(&input, "input", "157901", "puzzle input")
	flag.Parse()
	tmp, err := strconv.ParseUint(input, 10, 32)
	if err != nil {
		log.Fatalf("not numeric input")
	}
	bigLine := make([]byte, 2, tmp*150)

	elf1 := 0
	elf2 := 1

	bigLine[elf1] = 3
	bigLine[elf2] = 7

	sequence := []byte(input)
	for i, d := range sequence {
		sequence[i] = d - '0'
	}
	t := time.Now()
	for !HasSuffix(bigLine, sequence) {
		bigLine = append(bigLine, Itobs(int(bigLine[elf1]+bigLine[elf2]))...)
		elf1 += int(bigLine[elf1]) + 1
		elf1 %= len(bigLine)
		elf2 += int(bigLine[elf2]) + 1
		elf2 %= len(bigLine)
	}
	elapsed := time.Since(t)
	fmt.Printf("%s first appears after %d recipes : %s\n", input, bytes.LastIndex(bigLine, sequence), elapsed)
}

func Itobs(i int) []byte {
	if i < 10 {
		return []byte{byte(i)}
	}
	return []byte{byte(i / 10), byte(i % 10)}
}

func HasSuffix(s, suffix []byte) bool {
	n := len(s)
	sl := len(suffix)
	switch {
	case n > sl+1:
		return bytes.Contains(s[n-sl-1:], suffix)
	default:
		return bytes.HasSuffix(s, suffix)
	}
}
