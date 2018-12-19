package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type goblinElf int

const (
	goblin = iota
	elf
)

func (ge goblinElf) String() string {
	switch ge {
	case goblin:
		return "Goblin"
	case elf:
		return "Elf"
	default:
		return "NA"
	}
}

type point struct {
	x, y int
}

type map2d [][]int

func (m2d map2d) String() string {
	b := new(strings.Builder)
	line := make([]byte, len(m2d[0]))
	for y := range m2d {
		for x, v := range m2d[y] {
			if v > 127 {
				line[x] = '-'
			} else {
				line[x] = byte(v)
			}
		}
		b.Write(line)
		b.WriteByte('\n')
	}
	return b.String()
}

func (m2d map2d) Output(contestants fighters) string {
	b := new(strings.Builder)
	line := make([]byte, len(m2d[0]))
	ci := 0
	for y := range m2d {
		for x, v := range m2d[y] {
			line[x] = byte(v)
		}
		b.Write(line)
		var bits []string
		for ci < len(contestants) && contestants[ci].y == y {
			if contestants[ci].HP > 0 {
				bits = append(bits, contestants[ci].String())
			}
			ci++
		}
		if len(bits) > 0 {
			b.WriteString("  ")
			b.WriteString(strings.Join(bits, ", "))
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func (m map2d) Clone() map2d {
	ret := make(map2d, len(m))
	for y := range m {
		line := make([]int, len(m[y]))
		copy(line, m[y])
		ret[y] = line
	}
	return ret
}

func mapcopy(dst, src map2d) {
	if len(src) != len(dst) {
		panic("invalid lengthts on mapcopy")
	}
	for y := range src {
		copy(dst[y], src[y])
	}
}

func hamming(a, b point) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

type fighter struct {
	T  goblinElf
	HP int
	point
}

type path struct {
	*point
	ct int
}

func (f *fighter) Enemy(other *fighter) bool {
	return f.T != other.T
}

func (f *fighter) String() string {
	var t string
	switch f.T {
	case goblin:
		t = "G"
	case elf:
		t = "E"
	}
	return fmt.Sprintf("%s(%d,%d,%d)", t, f.HP, f.x, f.y)
}

func main() {

	var filePath string

	flag.StringVar(&filePath, "input", "", "specify input file.  defaults to stdin")
	flag.Parse()
	var reader *bufio.Scanner
	close := func() {}
	if filePath == "" {
		reader = bufio.NewScanner(os.Stdin)
	} else {
		fh, err := os.Open(filePath)
		if err != nil {
			log.Fatal("error opeining %s: %s", filePath, err)
		}
		reader = bufio.NewScanner(fh)
		close = func() { fh.Close() }
	}

	var grid map2d
	var contestants fighters
	for y := 0; reader.Scan(); y++ {
		buf := reader.Bytes()
		line := make([]int, len(buf))
		for i := range buf {
			line[i] = int(buf[i])
		}
		for x := 0; x < len(line); x++ {
			i := bytes.IndexAny(buf[x:], "GE")
			if i == -1 {
				break
			}
			x += i

			f := &fighter{
				HP: 200,
				point: point{
					x: x,
					y: y,
				},
			}

			switch line[x] {
			case 'G':
				f.T = goblin
			case 'E':
				f.T = elf
			}
			contestants = append(contestants, f)

		}
		grid = append(grid, line)
	}

	close()
	origGrid := grid.Clone()
	scratch := grid.Clone()
	origContestants := contestants.Clone()
	var rounds int
	for elfPowah := 4; ; elfPowah++ {
		mapcopy(grid, origGrid)
		contestants = origContestants.Clone()
		elfDied := false
		rounds = 0

		fmt.Printf("elf powah %d\n", elfPowah)
	OUTER:
		for ; ; rounds++ {
			fmt.Printf("going round %d\n", rounds)
			ci := 0
			for _, c := range contestants {
				if c.HP <= 0 {
					fmt.Printf("killing %d %#v\n", ci, c)
				} else {
					contestants[ci] = c
					ci++
				}
			}
			contestants = contestants[:ci]
			sort.Sort(contestants)

			fmt.Printf("%s\n", grid.Output(contestants))
			for i, me := range contestants {
				if me.HP <= 0 {
					continue
				}
				foundEnemies := false

				// find in range
				findInRange := func(me *fighter) fighters {
					var inRange fighters
					for j, them := range contestants {
						if i == j || them.HP <= 0 || !me.Enemy(them) {
							continue
						}
						if hamming(me.point, them.point) == 1 {
							inRange = append(inRange, them)
						}
					}
					return inRange
				}
				inRange := findInRange(me)
				fight := func(inRange fighters, me *fighter, elfPowah int) bool {
					fmt.Printf("fight!\n")
					minHP := -1
					var p *fighter
					for i := range inRange {
						if minHP == -1 ||
							inRange[i].HP < minHP ||
							inRange[i].HP == minHP && readingLess(&inRange[i].point, &p.point) {
							minHP = inRange[i].HP
							p = inRange[i]
						}
					}
					switch me.T {
					case elf:
						p.HP -= elfPowah
					default:
						p.HP -= 3
					}
					if p.HP <= 0 {
						switch p.T {
						case elf:
							return true
						}
						fmt.Printf("died!\n")
						grid[p.y][p.x] = '.'
					}
					return false
				}
				if len(inRange) > 0 {
					foundEnemies = true
					if elfDied = fight(inRange, me, elfPowah); elfDied {
						break OUTER
					}
					continue
				}
				// walk
				minDistance := -1
				var next *point
				var last *point
				for j, them := range contestants {
					if i == j || them.HP <= 0 || !me.Enemy(them) {
						continue
					}
					mapcopy(scratch, grid)

					foundEnemies = true
					ok := findPath(them.point, me.point, scratch)
					fmt.Printf("path done\n")

					if ok {
						fmt.Printf("path was okay!\n")
						// fmt.Printf("%#v => %#v\n%s\n", me, them, scratch)
						pt, end, ml := nextMove(me.point, them.point, scratch)
						fmt.Printf("ml was %d pt: %#v %#v\n", ml, pt, end)
						if minDistance == -1 ||
							ml < minDistance ||
							(ml == minDistance && readingLess(end, last)) {
							minDistance = ml
							next = pt
							last = end
						}
					}
				}
				if !foundEnemies {
					break OUTER
				}
				if minDistance == -1 {
					fmt.Printf("can't move!\n")
					continue
				}
				fmt.Printf("moving %#v to %#v\n", me, next)
				if me.point != *next {
					t := grid[me.y][me.x]
					switch t {
					case '.':
						panic("fuck shit")
					}
					grid[me.y][me.x] = '.'
					me.point = *next
					switch grid[me.y][me.x] {
					case '.':
					default:
						panic("shitfuck")

					}
					grid[me.y][me.x] = t
					mapcopy(scratch, grid)
				}

				inRange = findInRange(me)
				if len(inRange) > 0 {
					if elfDied = fight(inRange, me, elfPowah); elfDied {
						break OUTER
					}
				}
			}
		}
		if !elfDied {
			break
		}
	}

	winner := contestants[0].T
	hp := 0
	for _, c := range contestants {
		if c.HP <= 0 {
			continue
		}
		if c.T != winner {
			fmt.Printf("dammit!, winner was %d , this was %d\n", winner, c.T)
		}
		hp += c.HP
	}
	sort.Sort(contestants)
	fmt.Printf("%s\n", grid.Output(contestants))
	fmt.Printf("went %d round, %s wins with total score of %d * %d = %d\n", rounds, winner, rounds, hp, rounds*hp)
}

//l := grid[me.y][me.x]
//grid[me.y][me.x] = '.'
//me.y = pt.y
//me.x = pt.x
//grid[me.y][me.x] = l

func readingLess(a, b *point) bool {
	switch {
	case a.y < b.y:
		return true
	case a.y == b.y && a.x < b.x:
		return true
	default:
		return false
	}
}

type fighters []*fighter

func (f fighters) Len() int { return len(f) }

func (f fighters) Less(i, j int) bool {
	switch {
	case f[i].y < f[j].y:
		return true
	case f[i].y == f[j].y && f[i].x < f[j].x:
		return true
	}
	return false
}

func (f fighters) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f fighters) Clone() fighters {
	cp := make(fighters, len(f))
	for i := range f {
		e := &fighter{}
		*e = *f[i]
		cp[i] = e
	}
	return cp
}
func nextMove(start, end point, grid map2d) (*point, *point, int) {
	minScore := -1

	a := start
	var pl pathlist
	for {
		var pt *point
		if a.y > 0 && grid[a.y-1][a.x] > int('H') {
			g := grid[a.y-1][a.x]
			if minScore == -1 || g-int('H') < minScore {
				pt = &point{a.x, a.y - 1}
				minScore = g - int('H')
			}
		}
		if a.x > 0 && grid[a.y][a.x-1] > int('H') {
			g := grid[a.y][a.x-1]
			if minScore == -1 || g-int('H') < minScore {
				pt = &point{a.x - 1, a.y}
				minScore = g - int('H')
			}
		}
		if a.x < len(grid[0])-1 && grid[a.y][a.x+1] > int('H') {
			g := grid[a.y][a.x+1]
			if minScore == -1 || g-int('H') < minScore {
				pt = &point{a.x + 1, a.y}
				minScore = g - int('H')
			}
		}
		if a.y < len(grid)-1 && grid[a.y+1][a.x] > int('H') {
			g := grid[a.y+1][a.x]
			if minScore == -1 || g-int('H') < minScore {
				pt = &point{a.x, a.y + 1}
				minScore = g - int('H')
			}
		}
		pl = append(pl, path{point: pt, ct: minScore})
		a = *pt
		if minScore == 1 {
			break
		}

	}

	return pl[0].point, pl[len(pl)-1].point, pl[0].ct
}

func findPath(a, b point, grid map2d) bool {

	grid[a.y][a.x] = 'H'
	return findPathHelper(a, b, grid, 1)
}

func findPathHelper(a, b point, grid map2d, ct int) bool {
	found := false
	goodBail := false
	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == 'H'+ct-1 {
				if y > 0 && grid[y-1][x] == '.' {
					found = true
					grid[y-1][x] = 'H' + ct
				}
				if x > 0 && grid[y][x-1] == '.' {
					found = true
					grid[y][x-1] = 'H' + ct
				}
				if x < len(grid[y])-1 && grid[y][x+1] == '.' {
					found = true
					grid[y][x+1] = 'H' + ct
				}
				if y < len(grid)-1 && grid[y+1][x] == '.' {
					found = true
					grid[y+1][x] = 'H' + ct
				}
				h := hamming(point{x, y}, b)
				if h == 1 {
					goodBail = true
				}
			}
		}
	}

	if goodBail {
		return true
	}
	if !found {
		return false
	}
	return findPathHelper(a, b, grid, ct+1)
}

type pathlist []path

func (pl pathlist) Len() int { return len(pl) }

func (pl pathlist) Less(i, j int) bool {
	switch {
	case pl[i].ct < pl[j].ct:
		return true
	case pl[i].ct == pl[j].ct:
		switch {
		case pl[i].y < pl[j].y:
			return true
		case pl[i].y == pl[j].y && pl[i].x < pl[j].x:
			return true
		}
	}
	return false
}

func (pl pathlist) Swap(i, j int) {
	pl[i], pl[j] = pl[j], pl[i]
}
