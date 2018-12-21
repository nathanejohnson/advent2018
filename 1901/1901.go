package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type instruction func(input []int, registers []int) []int

var instructions = map[string]instruction{
	"addr": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[2]] = registers[input[0]] + registers[input[1]]
		return output
	},
	"addi": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[2]] = registers[input[0]] + input[1]
		return output
	},
	"mulr": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[2]] = registers[input[0]] * registers[input[1]]
		return output
	},
	"muli": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[2]] = registers[input[0]] * input[1]
		return output
	},
	"banr": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[2]] = registers[input[0]] & registers[input[1]]
		return output
	},
	"bani": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[2]] = registers[input[0]] & input[1]
		return output
	},
	"borr": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[2]] = registers[input[0]] | registers[input[1]]
		return output
	},
	"bori": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[2]] = registers[input[0]] | input[1]

		return output
	},
	"setr": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[2]] = registers[input[0]]
		return output
	},
	"seti": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		output[input[2]] = input[0]
		return output
	},
	"gtir": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		switch input[0] > registers[input[1]] {
		case true:
			output[input[2]] = 1
		default:
			output[input[2]] = 0
		}
		return output
	},
	"gtri": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		switch registers[input[0]] > input[1] {
		case true:
			output[input[2]] = 1
		default:
			output[input[2]] = 0
		}
		return output
	},
	"gtrr": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		switch registers[input[0]] > registers[input[1]] {
		case true:
			output[input[2]] = 1
		default:
			output[input[2]] = 0
		}
		return output
	},
	"eqir": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		switch input[0] == registers[input[1]] {
		case true:
			output[input[2]] = 1
		default:
			output[input[2]] = 0
		}
		return output
	},
	"eqri": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		switch registers[input[0]] == input[1] {
		case true:
			output[input[2]] = 1
		default:
			output[input[2]] = 0
		}
		return output
	},
	"eqrr": func(input []int, registers []int) []int {
		output := make([]int, len(registers))
		copy(output, registers)
		switch registers[input[0]] == registers[input[1]] {
		case true:
			output[input[2]] = 1
		default:
			output[input[2]] = 0
		}
		return output
	},
}

type step struct {
	call     string
	operands []int
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	var steps []step
	var ipLoc int

	reader.Scan()
	_, err := fmt.Sscanf(reader.Text(), "#ip %d", &ipLoc)
	if err != nil {
		log.Fatal("error reading ipLoc: %s", err)
	}

	for reader.Scan() {
		txt := reader.Text()
		bits := strings.Split(txt, " ")
		if len(bits) != 4 {
			log.Fatal("bad read: %s\n", txt)
		}

		st := step{
			operands: make([]int, 3),
		}
		st.call = bits[0]
		for i := 0; i < 3; i++ {
			var err error
			st.operands[i], err = strconv.Atoi(bits[i+1])
			if err != nil {
				log.Fatal("error parsing %s\n", txt)
			}
		}
		steps = append(steps, st)
	}

	registers := make([]int, 6)
	for registers[ipLoc] < len(steps) {
		st := steps[registers[ipLoc]]
		// fmt.Printf("calling %s(%#v, %#v\n", st.call, st.operands, registers)
		registers = instructions[st.call](st.operands, registers)
		registers[ipLoc]++
		// fmt.Printf("registers: %#v\n", registers)
	}
	fmt.Printf("part1: %#v\n", registers)

	for i, st := range steps {
		fmt.Printf("%02d: %s %d %d %d\n",
			i, st.call, st.operands[0], st.operands[1], st.operands[2])
	}
	registers = make([]int, 6)
	registers[0] = 1

	//// registers = []int{0, 10551292, 9, 10551293, 0, 1000}
	//registers = []int{0, 10551292, 3, 10551292, 0, 1}
	//for registers[ipLoc] < len(steps) && ct < 100000 {
	//	ct++
	//	st := steps[registers[ipLoc]]
	//	fmt.Printf("calling %s(%#v, %#v)\n", st.call, st.operands, registers)
	//	registers = instructions[st.call](st.operands, registers)
	//	fmt.Printf("registers: %#v\n", registers)
	//	registers[ipLoc]++
	//}
	// find sum of all pairs that multiply together to make 10551292
	num := 10551292
	sum := 0
	pairs := findFactorPairs(num)
	for _, f := range pairs {
		fmt.Printf("factor pair: %#v\n", f)
		sum = sum + f[0] + f[1]
	}
	fmt.Printf("sum: %d\n", sum)
}

func clone(i []int) []int {
	out := make([]int, len(i))
	copy(out, i)
	return out
}

func findFactorPairs(bigNum int) [][2]int {
	var out [][2]int
	stopVal := int(math.Sqrt(float64(bigNum)))
	for i := 1; i <= stopVal; i++ {
		if bigNum%i == 0 {
			div := bigNum / i
			out = append(out, [2]int{i, div})

		}
	}
	return out
}
