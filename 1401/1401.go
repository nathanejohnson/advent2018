package main

import (
	"flag"
	"fmt"
	"strconv"
)

func main() {

	var input int
	flag.IntVar(&input, "input", 157901, "puzzle input")
	flag.Parse()
	bigLine := make([]byte, 2, input+12)

	elf1 := 0
	elf2 := 1

	bigLine[elf1] = 3
	bigLine[elf2] = 7
	for len(bigLine) < input+10 {
		// combine
		newNum := int(bigLine[elf1] + bigLine[elf2])
		digs := strconv.Itoa(newNum)

		for _, d := range []byte(digs) {
			bigLine = append(bigLine, d-'0')
		}

		elf1 += int(bigLine[elf1]) + 1
		elf1 %= len(bigLine)
		elf2 += int(bigLine[elf2]) + 1
		elf2 %= len(bigLine)
	}
	fmt.Print("last 10: ")
	for _, d := range bigLine[input:] {
		fmt.Printf("%s", string(d+'0'))
	}
	fmt.Println("")
}
