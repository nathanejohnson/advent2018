package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

type worker struct {
	working bool
	job     string
	length  int
	elapsed int
}

func (w *worker) tick(data map[string][]string) string {
	if w.working {
		w.elapsed++
		if w.elapsed >= w.length {
			remove(w.job, data)
			w.working = false
			return w.job
		}
	}
	return ""
}

func (w *worker) takeJerb(job string) bool {
	if !w.working {
		w.job = job
		w.length = int(job[0] - 'A' + 61)
		w.elapsed = 0
		w.working = true

		return true
	}
	return false
}

func main() {

	fh, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	reader := bufio.NewScanner(fh)

	depList := make(map[string][]string)
	var empty []string
	for reader.Scan() {
		var a, b string
		_, err := fmt.Sscanf(reader.Text(), "Step %s must be finished before step %s can begin.", &a, &b)
		if err != nil {
			panic(err)
		}
		depList[b] = append(depList[b], a)
		depList[a] = append(depList[a], empty...)
	}

	var output string
	var workers [5]*worker
	var waiting []string

	for i := range workers {
		workers[i] = new(worker)
	}
	for i := 0; ; i++ {
		hits := readies(depList)
		waiting = append(waiting, hits...)
		sort.Strings(waiting)
		idleCount := 0
		for _, w := range workers {
			if !w.working {
				idleCount++
			}
		}
		if len(waiting) == 0 && idleCount == 5 {
			fmt.Printf("elapsed: %d , output: %s\n", i, output)
			break
		}
		for len(waiting) > 0 {
			taken := false
			for i := range workers {
				if taken = workers[i].takeJerb(waiting[0]); taken {
					break
				}
			}
			if !taken {
				break
			}
			delete(depList, waiting[0])
			waiting = waiting[1:]
		}
		for _, w := range workers {
			j := w.tick(depList)
			if j != "" {
				output += j
			}

		}
	}

}

func readies(data map[string][]string) []string {
	var out []string
	for k, v := range data {
		if len(v) == 0 {
			out = append(out, k)
		}
	}
	return out
}

func remove(key string, data map[string][]string) {
	for k, v := range data {
		switch k {
		case key:
			delete(data, k)
		default:
			var newV []string
			for _, s := range v {
				if s != key {
					newV = append(newV, s)
				}
			}
			data[k] = newV
		}
	}
}
