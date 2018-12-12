package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"time"
)

type rule struct {
	matches []byte
	state   byte
}

func main() {

	reader := bufio.NewScanner(os.Stdin)
	reader.Scan()
	var buf []byte
	_, err := fmt.Sscanf(reader.Text(), "initial state: %s", &buf)
	if err != nil {
		log.Fatal(err)
	}

	ll := 3

	empties := bytes.Repeat([]byte{'.'}, len(buf)+ll+25)
	init := make([]byte, len(empties))
	copy(init, empties)

	copy(init[ll:], buf)

	fmt.Printf("initial:\t   %s\nnew:\t\t%s\n", buf, init)

	reader.Scan()
	rules := make(map[string]string)

	for reader.Scan() {
		var matches, state string

		_, err = fmt.Sscanf(reader.Text(), "%s => %s", &matches, &state)
		if err != nil {
			log.Fatal("fuck " + err.Error())
		}
		if state == "#" {
			rules[matches] = state
		}
		fmt.Printf("'%s' => '%s'\n", matches, state)
	}

	t := time.Now()
	lastGen := init
	thisGen := make([]byte, len(init))
	fmt.Printf("%d\t%s\n", 0, init)
	for gen := 1; gen <= 20; gen++ {
		copy(thisGen, empties)
		for i := 2; i < len(init)-2; i++ {
			chunk := string(lastGen[i-2 : i+3])
			if m, ok := rules[chunk]; ok {
				thisGen[i] = byte(m[0])
			}
		}
		fmt.Printf("%d\t%s\n", gen, thisGen)
		copy(lastGen, thisGen)

	}
	tot := 0
	for i, v := range lastGen {
		if v == byte('#') {
			tot += i - ll
		}
	}

	elapsed := time.Since(t)

	fmt.Printf("total: %d , elapsed: %s\n", tot, elapsed)
}
