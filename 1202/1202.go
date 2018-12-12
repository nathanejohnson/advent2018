package main

import (
	"bufio"
	"bytes"
	"flag"
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
	var lpad, rpad int
	flag.IntVar(&lpad, "lpad", 5, "left padding")
	flag.IntVar(&rpad, "rpad", 25, "right padding")

	flag.Parse()
	reader := bufio.NewScanner(os.Stdin)
	reader.Scan()
	var buf []byte
	_, err := fmt.Sscanf(reader.Text(), "initial state: %s", &buf)
	if err != nil {
		log.Fatal(err)
	}

	empties := bytes.Repeat([]byte{'.'}, len(buf)+lpad+rpad)
	init := make([]byte, len(empties))
	copy(init, empties)

	copy(init[lpad:], buf)

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
	}

	t := time.Now()
	lastGen := init
	thisGen := make([]byte, len(init))
	fmt.Printf("%d\t%s\n", 0, init)
	safety := []byte{'.', '.', '.'}
	gen := 0
	score := 0
	delta := 0
	for ; gen < 50000000000; gen++ {
		if bytes.Compare(lastGen[0:3], safety) != 0 {
			fmt.Printf("dammit before: %d: %s\n", gen, lastGen)
			return
		}
		if bytes.Compare(lastGen[len(lastGen)-3:], safety) != 0 {
			fmt.Printf("dammit after %d: %s\n", gen, lastGen)
			return
		}
		copy(thisGen, empties)
		for i := 2; i < len(init)-2; i++ {
			chunk := string(lastGen[i-2 : i+3])
			if m, ok := rules[chunk]; ok {
				thisGen[i] = byte(m[0])
			}
		}
		newScore := Score(thisGen, lpad)
		newDelta := newScore - score
		if newDelta == delta {
			fmt.Printf("we're stable!\n")
			break
		}
		delta = newDelta
		score = newScore
		copy(lastGen, thisGen)
		fmt.Printf("checkpoint gen %d: score: %d delta: %d\n", gen, score, delta)
	}

	score += delta * (50000000000 - gen)

	elapsed := time.Since(t)

	fmt.Printf("total: %d , elapsed: %s\n", score, elapsed)
}

func Score(data []byte, lpad int) int {
	plant := byte('#')
	score := 0
	for i := range data {
		if data[i] == plant {
			score += i - lpad
		}
	}
	return score
}
