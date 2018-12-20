package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type instruction func(input []int, registers []int) []int

var instructions = map[string]instruction{
	"addr": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[3]] = registers[input[1]] + registers[input[2]]
		return output
	},
	"addi": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[3]] = registers[input[1]] + input[2]
		return output
	},
	"mulr": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[3]] = registers[input[1]] * registers[input[2]]
		return output
	},
	"muli": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[3]] = registers[input[1]] * input[2]
		return output
	},
	"banr": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[3]] = registers[input[1]] & registers[input[2]]
		return output
	},
	"bani": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[3]] = registers[input[1]] & input[2]
		return output
	},
	"borr": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[3]] = registers[input[1]] | registers[input[2]]
		return output
	},
	"bori": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[3]] = registers[input[1]] | input[2]

		return output
	},
	"setr": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[3]] = registers[input[1]]
		return output
	},
	"seti": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[3]] = input[1]
		return output
	},
	"gtir": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		switch input[1] > registers[input[2]] {
		case true:
			output[input[3]] = 1
		default:
			output[input[3]] = 0
		}
		return output
	},
	"gtri": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		switch registers[input[1]] > input[2] {
		case true:
			output[input[3]] = 1
		default:
			output[input[3]] = 0
		}
		return output
	},
	"gtrr": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		switch registers[input[1]] > registers[input[2]] {
		case true:
			output[input[3]] = 1
		default:
			output[input[3]] = 0
		}
		return output
	},
	"eqir": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		switch input[1] == registers[input[2]] {
		case true:
			output[input[3]] = 1
		default:
			output[input[3]] = 0
		}
		return output
	},
	"eqri": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		switch registers[input[1]] == input[2] {
		case true:
			output[input[3]] = 1
		default:
			output[input[3]] = 0
		}
		return output
	},
	"eqrr": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
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

var opmap map[int]string

func call(k string, input, registers, output []int) (passed bool) {
	defer func() {
		if pan := recover(); pan != nil {
			passed = false
		}
	}()
	o := instructions[k](input, registers)
	for i := range o {
		if o[i] != output[i] {
			return false
		}
	}
	return true

}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	before := make([]int, 4)
	data := make([]int, 4)
	after := make([]int, 4)
	opmap = make(map[int]string)

	possibilities := make(map[int]map[string]struct{})
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
			&data[0], &data[1], &data[2], &data[3])

		if err != nil {
			log.Fatalf("fuck data %s", err)
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
		var lastK string
		op := data[0]

		for k := range instructions {
			if call(k, data, before, after) {
				if _, ok := possibilities[op]; !ok {
					possibilities[op] = make(map[string]struct{})
				}
				possibilities[op][k] = struct{}{}
				works++
				lastK = k
			}
		}
		if works == 1 {
			if v, ok := opmap[op]; ok && v != lastK {
				panic("fuck")
			}
			opmap[op] = lastK
			delete(possibilities, op)
			reduce(lastK, possibilities)
		}
	}

	for len(possibilities) > 0 {
		for k := range possibilities {
			if len(possibilities[k]) == 1 {
				for v := range possibilities[k] {
					opmap[k] = v
					delete(possibilities, k)
					reduce(v, possibilities)
				}
			}
		}
		for k := range possibilities {
			for v := range possibilities[k] {
				vct := 1
				for k2 := range possibilities {
					if k2 == k {
						continue
					}
					if _, ok := possibilities[k2][v]; ok {
						vct++
						break
					}
				}
				if vct == 1 {
					if _, ok := opmap[k]; ok {
						panic("fuck it")
					}
					opmap[k] = v
					delete(possibilities, k)
					reduce(v, possibilities)
				}
			}
		}
	}

	reader.Scan()
	registers := make([]int, 4)
	lines := 0
	for reader.Scan() {
		txt := reader.Text()
		_, err := fmt.Sscanf(txt,
			"%d %d %d %d",
			&data[0], &data[1], &data[2], &data[3])

		if err != nil {
			fmt.Printf("fuck data %s", err)
			continue
		}
		// fmt.Printf("data: %s %#v\n", txt, data)
		lines++

		key := opmap[data[0]]
		// fmt.Printf("executing %s with %#v registers %#v\n", key, data, registers)
		output := instructions[key](data, registers)
		if output == nil {
			log.Fatal("died on line %d\n", lines)
		}
		// fmt.Printf("after: %#v\n", output)
		registers = output
	}

	fmt.Printf("lines: %d\n", lines)

	fmt.Printf("register[0]: %d\n", registers[0])
}

func reduce(name string, m map[int]map[string]struct{}) {
	for _, outer := range m {
		delete(outer, name)
	}
}
