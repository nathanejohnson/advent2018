package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type node struct {
	childCt    int
	metadataCt int
	children   []*node
	value      int
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
	root := new(node)
	root.Read(data)
	fmt.Printf("part2: %d\n", root.value)
}

func (n *node) Read(data []string) []string {
	header := data[:2]
	data = data[2:]
	n.childCt = digit(header[0])
	n.metadataCt = digit(header[1])

	for i := 0; i < n.childCt; i++ {
		child := new(node)
		data = child.Read(data)
		n.children = append(n.children, child)
	}
	value := 0
	for i := 0; i < n.metadataCt; i++ {
		ct := digit(data[i])
		switch {
		case n.childCt == 0:
			value += ct
		case 0 < ct && ct <= n.childCt:
			value += n.children[ct-1].value
		}
	}
	n.value = value
	data = data[n.metadataCt:]
	return data
}

func digit(s string) int {
	d, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return d
}
