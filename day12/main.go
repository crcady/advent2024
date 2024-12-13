package main

import (
	"bufio"
	"log"
	"os"
)

type point struct {
	x int
	y int
}

func (p point) adjacent() []point {
	return []point{{p.x - 1, p.y}, {p.x + 1, p.y}, {p.x, p.y - 1}, {p.x, p.y + 1}}
}

type region struct {
	points map[point]bool
}

func NewRegion(p point) *region {
	return &region{map[point]bool{p: true}}
}

func (r region) add(p point) {
	r.points[p] = true
}

func (r region) area() int {
	return len(r.points)
}

func (r region) perimeter() int {
	res := 0
	for p1 := range r.points {
		for _, p2 := range p1.adjacent() {
			if !r.points[p2] {
				res++
			}
		}
	}
	return res
}

func (r region) price() int {
	return r.perimeter() * r.area()
}

type gardenMap struct {
	rows   []string
	height int
}

func newGardenMap() *gardenMap {
	return &gardenMap{make([]string, 0), 0}
}
func (gm *gardenMap) addRow(r string) {
	gm.rows = append(gm.rows, r)
	gm.height++
}

func (gm *gardenMap) get(p point) byte {
	return gm.rows[p.x][p.y]
}

func (gm *gardenMap) onMap(p point) bool {
	return p.x >= 0 && p.y >= 0 && p.x < gm.height && p.y < gm.height
}

func (gm *gardenMap) regions() []region {
	res := make([]region, 0)
	visited := make(map[point]bool)
	toVisit := []point{{0, 0}}

	for len(visited) < gm.height*gm.height {
		nextCandidate := toVisit[len(toVisit)-1]
		toVisit = toVisit[:len(toVisit)-1]

		for visited[nextCandidate] {
			nextCandidate = toVisit[len(toVisit)-1]
			toVisit = toVisit[:len(toVisit)-1]
		}

		// Now we have an unvisited point, meaning it is in a new region
		r := NewRegion(nextCandidate)
		inCurrentRegion := []point{nextCandidate}

		for len(inCurrentRegion) > 0 {
			p := inCurrentRegion[len(inCurrentRegion)-1]
			inCurrentRegion = inCurrentRegion[:len(inCurrentRegion)-1]
			if visited[p] {
				continue
			}

			r.add(p)
			visited[p] = true

			for _, n := range p.adjacent() {
				if !gm.onMap(n) {
					continue
				}

				if gm.get(n) == gm.get(p) { //same region
					inCurrentRegion = append(inCurrentRegion, n)
				} else {
					toVisit = append(toVisit, n)
				}
			}
		}
		res = append(res, *r)
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

	scanner := bufio.NewScanner(file)

	gm := newGardenMap()

	for scanner.Scan() {
		gm.addRow(scanner.Text())
	}

	log.Printf("Read %d rows from %s\n", gm.height, fname)
	regions := gm.regions()
	log.Printf("Found %d regions\n", len(regions))

	ans1 := 0
	for _, r := range regions {
		ans1 += r.price()
	}
	log.Printf("The price for part one is %d", ans1)
}
