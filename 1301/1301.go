package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

type direction int

const (
	left direction = iota
	up
	right
	down
)

type turn int

const (
	leftTurn turn = iota
	straight
	rightTurn
)

type cart struct {
	x, y      int
	direction direction
	nextTurn  turn
}

func (c *cart) Tick(grid [][]byte) {

	switch c.direction {
	case left:
		c.x--
	case right:
		c.x++
	case down:
		c.y++
	case up:
		c.y--
	}

	switch grid[c.y][c.x] {
	case '+':
		switch c.nextTurn {
		case leftTurn:
			c.direction--
			if c.direction < 0 {
				c.direction = down
			}
		case rightTurn:
			c.direction++
		}
		c.direction %= 4
		c.nextTurn++
		c.nextTurn %= 3
	case '\\':
		switch c.direction {
		case left:
			c.direction = up
		case right:
			c.direction = down
		case up:
			c.direction = left
		case down:
			c.direction = right
		default:
			panic("fuck")
		}
	case '/':
		switch c.direction {
		case left:
			c.direction = down
		case right:
			c.direction = up
		case up:
			c.direction = right
		case down:
			c.direction = left
		default:
			panic("fuck")
		}
	}
}

type Carts []*cart

func (c Carts) Len() int { return len(c) }

func (c Carts) Less(i, j int) bool {
	switch {
	case c[i].y < c[j].y:
		return true
	case c[i].y == c[j].y && c[i].x < c[j].x:
		return true
	}
	return false
}

func (c Carts) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Carts) Collision() (bool, Carts) {
	m := make(map[string]*cart)
	f := func(x, y int) string {
		return strconv.Itoa(x) + ":" + strconv.Itoa(y)
	}
	collision := false
	var deletes []string
	for _, crt := range c {
		key := f(crt.x, crt.y)
		if _, ok := m[key]; ok {
			collision = true
			deletes = append(deletes, key)
			continue
		}
		m[key] = crt
	}
	if !collision {
		return false, c
	}
	for _, k := range deletes {
		delete(m, k)
	}
	newCarts := make(Carts, 0, len(m))
	for _, v := range m {
		newCarts = append(newCarts, v)
	}
	return true, newCarts

}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	var grid [][]byte

	var carts Carts
	t1 := time.Now()
	for y := 0; reader.Scan(); y++ {
		buf := reader.Bytes()
		line := make([]byte, len(buf))
		copy(line, buf)
		grid = append(grid, line)
		for {
			i := bytes.IndexFunc(line, func(r rune) bool {
				switch r {
				case '>', '<', '^', 'v':
					return true
				}
				return false
			})
			if i == -1 {
				break
			}

			cart := &cart{
				x:         i,
				y:         y,
				direction: getDirection(line[i]),
			}
			carts = append(carts, cart)
			switch cart.direction {
			case left, right:
				line[i] = '-'
			case up, down:
				line[i] = '|'
			}

			if i+1 == len(line) {
				break
			}
			line = line[i+1:]
		}
	}

	t2 := time.Now()

	for len(carts) > 1 {
		sort.Sort(carts)

		for passCarts := carts; len(passCarts) > 0; passCarts = passCarts[1:] {
			c := passCarts[0]
			c.Tick(grid)
			var collided bool
			if collided, carts = carts.Collision(); collided {
				fmt.Printf("collision at: %d,%d\n", c.x, c.y)
			}

		}

	}
	t1Elapsed := time.Since(t1)
	t2Elapsed := time.Since(t2)
	fmt.Printf("pos: %d,%d , time since before reading file: %s, time after: %s\n", carts[0].x, carts[0].y,
		t1Elapsed, t2Elapsed)
}

func getDirection(b byte) direction {
	switch b {
	case '<':
		return left
	case '>':
		return right
	case '^':
		return up
	case 'v':
		return down
	}
	panic("fuck")
}
