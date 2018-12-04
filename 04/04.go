package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"time"
)

func main() {
	const format = "2006-01-02 15:04"
	fh, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()

	reader := bufio.NewScanner(fh)
	re := regexp.MustCompile(`^\[([^\]]+)\]\s+(.*)$`)

	var es entries
	for reader.Scan() {
		e := entry{}
		t := reader.Text()

		matches := re.FindStringSubmatch(t)
		if len(matches) != 3 {
			panic(fmt.Sprintf("fuck: %s, %d", t, len(matches)))
		}

		e.ts, err = time.Parse(format, matches[1])
		if err != nil {
			panic(err)
		}
		e.vs = matches[2]
		es = append(es, e)
	}
	sort.Sort(es)

	fsm := Fsm()
	for _, e := range es {
		fsm.process(e)
	}

	g := fsm.sleepiest()
	fmt.Printf("sleepiest guard: %#v\n", g)

	m := g.sleepiestMinute()
	fmt.Printf("sleepiest minute: %d\n", m)

	fmt.Printf("id * minute: %d\n", g.id*m)

	g, m = fsm.sleepiestGuardMinute()

	fmt.Printf("gid: %d, minute: %d, id*m: %d\n", g.id, m, g.id*m)
}

type fsm struct {
	lastG  *guard
	guards map[int]*guard
}

func Fsm() fsm {
	return fsm{
		guards: make(map[int]*guard),
	}
}

type state int

const (
	awake state = iota
	asleep
)

type guard struct {
	id         int
	sleepMins  map[int]int
	sleepTotal int
	state      state
	lastMin    int
}

func Guard(id int) *guard {
	return &guard{
		id:        id,
		sleepMins: make(map[int]int),
	}
}

func (g *guard) transition(s state, minute int) {
	switch g.state {
	case asleep:
		switch s {
		case awake:
			for t := g.lastMin; t < minute; t++ {
				g.sleepMins[t]++
				g.sleepTotal++
			}
		}
	case awake:
		switch s {
		case asleep:
			g.lastMin = minute
		}
	}
	g.state = s
}

func (g *guard) sleepiestMinute() int {
	var topMinute, topMinuteTotal int

	for h, c := range g.sleepMins {
		if c > topMinuteTotal {
			topMinute = h
			topMinuteTotal = c
		}
	}
	return topMinute
}

func (f *fsm) process(e entry) {
	var newGuardId int
	_, err := fmt.Sscanf(e.vs, "Guard #%d begins shift", &newGuardId)
	if err == nil {
		if f.lastG != nil {
			f.lastG.transition(awake, 60)
		}
		var m *guard
		if m = f.guards[newGuardId]; m == nil {
			m = Guard(newGuardId)
			f.guards[newGuardId] = m
		}
		f.lastG = m
		return
	}
	switch e.vs {
	case "falls asleep":
		f.lastG.transition(asleep, e.ts.Minute())
	case "wakes up":
		f.lastG.transition(awake, e.ts.Minute())
	default:
		panic("unexpected state")
	}
}

func (f *fsm) sleepiest() *guard {
	var m *guard
	max := 0
	for _, guard := range f.guards {
		if guard.sleepTotal > max {
			m = guard
			max = guard.sleepTotal
		}
	}
	return m
}

func (f *fsm) sleepiestGuardMinute() (*guard, int) {
	var mon *guard
	var maxMinute, maxMinuteCt int

	for _, guard := range f.guards {
		for min, ct := range guard.sleepMins {
			if ct > maxMinuteCt {
				mon = guard
				maxMinuteCt = ct
				maxMinute = min
			}
		}
	}
	return mon, maxMinute
}

type entry struct {
	ts time.Time
	vs string
}

type entries []entry

func (e entries) Len() int {
	return len(e)
}

func (e entries) Less(i, j int) bool {
	return e[i].ts.Before(e[j].ts)
}
func (e entries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
