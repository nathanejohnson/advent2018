package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	fh, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	haystack := make([]string, 0, 250)
	defer fh.Close()
	reader := bufio.NewScanner(fh)
	for lineNo := 0; reader.Scan(); lineNo++ {
		needle := reader.Text()
		for oldLineNo, prevLine := range haystack {
			score := 0
			diffIdx := -1
			if len(prevLine) != len(needle) {
				continue
			}
			for i := range needle {
				if prevLine[i] != needle[i] {
					score++
					diffIdx = i
				}
				if score > 1 {
					break
				}
			}
			if score <= 1 {
				var newStr string
				switch diffIdx {
				case -1:
					fmt.Printf("this shouldn't happen: %s, %s\n", prevLine, needle)
					return
				default:
					newStr = needle[:diffIdx] + needle[diffIdx+1:]
				}
				fmt.Printf("Holy shit: needle: %d:%d %s : %s, %s idx: %d, %d\n",
					lineNo, oldLineNo, prevLine, needle, newStr, diffIdx, len(prevLine))
				return
			}
		}
		haystack = append(haystack, needle)
	}
}
