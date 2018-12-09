package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

type entry struct {
	cw  *entry
	ccw *entry
	val int
}

type list struct {
	entries []*entry
	current *entry
	len     int
}

func (l *list) removeCurrent() {
	c := l.current
	c.ccw.cw = c.cw
	c.cw.ccw = c.ccw
	l.current = c.cw
	l.len--
}

func (l *list) insert(v int) {
	e := &entry{
		val: v,
	}
	c := l.current
	e.ccw = c
	e.cw = c.cw
	c.cw = e
	e.cw.ccw = e
	l.current = e
	l.len++
}

func InitList() *list {
	root := new(entry)
	root.cw = root
	root.ccw = root
	return &list{
		entries: []*entry{root},
		current: root,
		len:     1,
	}
}

func main() {
	fh, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	reader := bufio.NewScanner(fh)

	reader.Scan()
	var players, lastPts int
	_, err = fmt.Sscanf(reader.Text(), "%d players; last marble is worth %d points", &players, &lastPts)
	if err != nil {
		panic(err)
	}

	l := InitList()

	scores := make([]int, players)
	for i := 1; i < lastPts; i++ {
		// fmt.Printf("%d\n", i)
		if i%23 == 0 {
			pid := i % players
			scores[pid] += i
			fmt.Printf("scores[pid] after i: %d\n", scores[pid])
			for x := 0; x < 7; x++ {
				l.current = l.current.ccw
			}
			scores[pid] += l.current.val
			fmt.Printf("scores[pid] after current: %d, val: %d\n", scores[pid], l.current.val)
			l.removeCurrent()
		} else {
			l.current = l.current.cw
			// fmt.Printf("%d: before %#v, %#p\n", i, l.current, l.current)
			l.insert(i)
			// fmt.Printf("%d: after %#v, %#p\n", i, l.current, l.current)
		}
		// c := l.current
		//for x := 0; x < l.len; x++ {
		//	fmt.Printf("%d ", c.val)
		//	c = c.cw
		//
		//}
		//fmt.Println()
	}
	sort.Ints(scores)
	fmt.Printf("part1: highest score %d\n", scores[players-1])

}
