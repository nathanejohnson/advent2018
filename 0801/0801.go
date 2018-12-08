package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type node struct {
	childCt    int
	metadataCt int
	children   []*node
}

func main() {
	fh, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	reader := bufio.NewScanner(fh)

	reader.Scan()
	data := strings.Split(reader.Text(), " ")
	t := time.Now()
	root := &node{}
	_, count := root.Read(data)
	elapsed := time.Since(t)
	fmt.Printf("part1: %d, elapsed: %s\n", count, elapsed.String())
}

func (n *node) Read(data []string) ([]string, int) {
	header := data[:2]
	data = data[2:]
	n.childCt = digit(header[0])
	n.metadataCt = digit(header[1])
	value := 0
	for i := 0; i < n.childCt; i++ {
		child := &node{}
		var ct int
		data, ct = child.Read(data)
		value += ct
		n.children = append(n.children, child)
	}

	for i := 0; i < n.metadataCt; i++ {
		ct := digit(data[i])
		value += ct
	}
	data = data[n.metadataCt:]
	return data, value
}

func digit(s string) int {
	d, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return d
}
