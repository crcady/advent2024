package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func line2ints(s string) []int {
	res := make([]int, len(s))
	for i := range s {
		n, err := strconv.Atoi(string(s[i]))
		if err != nil {
			panic(err)
		}
		res[i] = n
	}
	return res
}

type point struct {
	x int
	y int
}

func (p point) up() point {
	return point{p.x - 1, p.y}
}

func (p point) down() point {
	return point{p.x + 1, p.y}
}

func (p point) left() point {
	return point{p.x, p.y - 1}
}

func (p point) right() point {
	return point{p.x, p.y + 1}
}

type topoMap [][]int

func (tm topoMap) zeros() []point {
	res := make([]point, 0)
	for i, slc := range tm {
		for j, n := range slc {
			if n == 0 {
				res = append(res, point{i, j})
			}
		}
	}
	return res
}

func (tm topoMap) valid(p point) bool {
	maxX, maxY := len(tm)-1, len(tm[0])-1
	return p.x >= 0 && p.y >= 0 && p.x <= maxX && p.y <= maxY
}

func (tm topoMap) height(p point) int {
	return tm[p.x][p.y]
}

func (tm topoMap) score(p point, onlyVisitOnce bool) int {
	visited := make(map[point]bool)
	toVisit := make([]point, 0)
	res := 0

	toVisit = append(toVisit, p)
	for len(toVisit) > 0 {
		p := toVisit[len(toVisit)-1]
		toVisit = toVisit[:len(toVisit)-1]

		if onlyVisitOnce {
			visited[p] = true
		}

		if tm.height(p) == 9 {
			res++
			continue
		}

		candidates := []point{p.up(), p.down(), p.left(), p.right()}
		for _, c := range candidates {
			if !tm.valid(c) {
				continue
			}

			if visited[c] {
				continue
			}

			if tm.height(c) == tm.height(p)+1 {
				toVisit = append(toVisit, c)
			}

		}

	}
	return res
}

func main() {
	fname := "example.txt"

	if len(os.Args) > 1 {
		fname = os.Args[1]
	}

	file, err := os.Open(fname)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.Println("Opened", fname)

	scanner := bufio.NewScanner(file)

	var tm topoMap
	for scanner.Scan() {
		line := scanner.Text()
		tm = append(tm, line2ints(line))
	}

	ans1 := 0
	for _, p := range tm.zeros() {
		ans1 += tm.score(p, true)
	}
	log.Println("Answer to first half:", ans1)

	ans2 := 0
	for _, p := range tm.zeros() {
		ans2 += tm.score(p, false)
	}
	log.Println("Answer to second half:", ans2)

}
