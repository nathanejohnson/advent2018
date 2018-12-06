package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type Entry struct {
	v    byte
	next *Entry
	prev *Entry
}

func (l *Entry) Next() *Entry {
	if l == nil {
		return nil
	}
	return l.next
}

func (ll *LinkedList) Remove(ch byte) {
	lc := ch
	if lc <= 'Z' {
		lc += 32
	}
	uc := lc - 32
	cur := ll.front
	var prev *Entry
	for cur != nil {
		switch cur.v {
		case uc, lc:
			if cur == ll.front {
				cur = cur.next
				ll.front = cur
				cur.prev = nil
			} else {
				cur = cur.next
				prev.next = cur
				if cur != nil {
					cur.prev = prev
				}
			}
			ll.length--
		default:
			prev = cur
			cur = cur.Next()
		}
	}
}

func (ll *LinkedList) Reduce() {
	cur := ll.front
	var prev *Entry
	for cur != nil && cur.next != nil {
		if abs(int(cur.v)-int(cur.next.v)) == 32 {
			if cur == ll.front {
				ll.front = cur.Next().Next()
				cur = ll.front
				cur.prev = ll.front
			} else {
				prev.next = cur.Next().Next()
				prev.next.prev = prev
				cur = prev
				prev = prev.prev
			}

			ll.length -= 2
		} else {
			prev = cur
			cur = cur.Next()
		}
	}
}

func (ll *LinkedList) validateLen() int {
	n := ll.front
	ct := 0
	for n != nil {
		ct++
		n = n.next
	}
	return ct
}

func (ll *LinkedList) gimme() []byte {
	data := make([]byte, 0, ll.length)
	n := ll.front
	for n != nil {
		data = append(data, n.v)
		n = n.next
	}
	return data
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

type LinkedList struct {
	entries []Entry
	front   *Entry
	length  int
}

func (lls *LinkedList) Init(data []byte) {
	entries := lls.entries
	if len(entries) != len(data) {
		entries = make([]Entry, len(data))
	}
	lls.front = &entries[0]
	lls.length = len(entries)
	l := len(data)
	for i := 0; i < l; i++ {
		if i < l-1 {
			entries[i].next = &entries[i+1]
		} else {
			entries[i].next = nil
		}
		if i > 0 {
			entries[i].prev = &entries[i-1]
		}
		entries[i].v = data[i]
	}
	lls.entries = entries
}

func main() {
	fh, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	reader := bufio.NewScanner(fh)
	reader.Scan()
	input := reader.Bytes()

	startTime := time.Now()
	ll := &LinkedList{}
	ll.Init(input)
	ll.Reduce()
	fmt.Printf("part1: %d\n", ll.length)
	part1Solution := ll.gimme()
	shortest := len(part1Solution)
	for c := 'a'; c <= 'z'; c++ {
		ll.Init(part1Solution)
		ll.Remove(byte(c))
		ll.Reduce()
		if ll.length < shortest {
			shortest = ll.length
		}

	}
	fmt.Printf("part2: %d\n", shortest)
	fmt.Printf("elapsed: %s\n", time.Since(startTime))
}
