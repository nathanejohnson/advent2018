package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"time"
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

	lastPts *= 100

	l := InitList()

	scores := make([]int, players)

	t := time.Now()
	for i := 1; i < lastPts; i++ {
		// fmt.Printf("%d\n", i)
		if i%23 == 0 {
			pid := i % players
			scores[pid] += i
			for x := 0; x < 7; x++ {
				l.current = l.current.ccw
			}
			scores[pid] += l.current.val
			l.removeCurrent()
		} else {
			l.current = l.current.cw
			l.insert(i)
		}
	}
	sort.Ints(scores)
	elapsed := time.Since(t)
	fmt.Printf("part1: highest score %d, elapsed: %s\n", scores[players-1], elapsed)

}
