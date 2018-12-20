package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

type instruction func(input []byte, registers []byte) []byte

var instructions = map[string]instruction{
	"addr": func(input []byte, registers []byte) []byte {
		output := make([]byte, len(registers))
		copy(output, registers)
		output[input[3]] = registers[input[1]] + registers[input[2]]
		return output
	},
	"addi": func(input []byte, registers []byte) []byte {
		output := make([]byte, len(registers))
		copy(output, registers)
		output[input[3]] = registers[input[1]] + input[2]
		return output
	},
	"mulr": func(input []byte, registers []byte) []byte {
		output := make([]byte, len(registers))
		copy(output, registers)
		output[input[3]] = registers[input[1]] * registers[input[2]]
		return output
	},
	"muli": func(input []byte, registers []byte) []byte {
		output := make([]byte, len(registers))
		copy(output, registers)
		output[input[3]] = registers[input[1]] * input[2]
		return output
	},
	"banr": func(input []byte, registers []byte) []byte {
		output := make([]byte, len(registers))
		copy(output, registers)
		output[input[3]] = registers[input[1]] & registers[input[2]]
		return output
	},
	"bani": func(input []byte, registers []byte) []byte {
		output := make([]byte, len(registers))
		copy(output, registers)
		output[input[3]] = registers[input[1]] & input[2]
		return output
	},
	"borr": func(input []byte, registers []byte) []byte {
		output := make([]byte, len(registers))
		copy(output, registers)
		output[input[3]] = registers[input[1]] | registers[input[2]]
		return output
	},
	"bori": func(input []byte, registers []byte) []byte {
		output := make([]byte, len(registers))
		copy(output, registers)
		output[input[3]] = registers[input[1]] | input[2]

		return output
	},
	"setr": func(input []byte, registers []byte) []byte {
		output := make([]byte, len(registers))
		copy(output, registers)
		output[input[3]] = registers[input[1]]
		return output
	},
	"seti": func(input []byte, registers []byte) []byte {
		output := make([]byte, len(registers))
		copy(output, registers)
		output[input[3]] = input[1]
		return output
	},
	"gtir": func(input []byte, registers []byte) []byte {
		output := make([]byte, len(registers))
		copy(output, registers)
		switch input[1] > registers[input[2]] {
		case true:
			output[input[3]] = 1
		default:
			output[input[3]] = 0
		}
		return output
	},
	"gtri": func(input []byte, registers []byte) []byte {
		output := make([]byte, len(registers))
		copy(output, registers)
		switch registers[input[1]] > input[2] {
		case true:
			output[input[3]] = 1
		default:
			output[input[3]] = 0
		}
		return output
	},
	"gtrr": func(input []byte, registers []byte) []byte {
		output := make([]byte, len(registers))
		copy(output, registers)
		switch registers[input[1]] > registers[input[2]] {
		case true:
			output[input[3]] = 1
		default:
			output[input[3]] = 0
		}
		return output
	},
	"eqir": func(input []byte, registers []byte) []byte {
		output := make([]byte, len(registers))
		copy(output, registers)
		switch input[1] == registers[input[2]] {
		case true:
			output[input[3]] = 1
		default:
			output[input[3]] = 0
		}
		return output
	},
	"eqri": func(input []byte, registers []byte) []byte {
		output := make([]byte, len(registers))
		copy(output, registers)
		switch registers[input[1]] == input[2] {
		case true:
			output[input[3]] = 1
		default:
			output[input[3]] = 0
		}
		return output
	},
	"eqrr": func(input []byte, registers []byte) []byte {
		output := make([]byte, len(registers))
		copy(output, registers)
		switch registers[input[1]] == registers[input[2]] {
		case true:
			output[input[3]] = 1
		default:
			output[input[3]] = 0
		}
		return output
	},
}

func call(k string, input, registers, output []byte) (passed bool) {
	defer func() {
		if pan := recover(); pan != nil {
			passed = false
		}
	}()
	o := instructions[k](input, registers)
	return bytes.Compare(output, o) == 0

}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	before := make([]byte, 4)
	op := make([]byte, 4)
	after := make([]byte, 4)
	threes := 0
	for reader.Scan() {
		_, err := fmt.Sscanf(reader.Text(),
			"Before: [%d, %d, %d, %d]",
			&before[0], &before[1], &before[2], &before[3])
		if err != nil {
			break
		}
		reader.Scan()
		_, err = fmt.Sscanf(reader.Text(),
			"%d %d %d %d",
			&op[0], &op[1], &op[2], &op[3])

		if err != nil {
			log.Fatalf("fuck op %s", err)
		}

		reader.Scan()
		_, err = fmt.Sscanf(reader.Text(),
			"After: [%d, %d, %d, %d]",
			&after[0], &after[1], &after[2], &after[3])

		if err != nil {
			log.Fatalf("fuck after %s", err)
		}
		reader.Scan()
		works := 0
		for k := range instructions {
			if call(k, op, before, after) {
				works++
				if works == 3 {
					threes++
					break
				}
			}
		}
	}

	for reader.Scan() {
	}
	fmt.Printf("threes: %d\n", threes)
}
