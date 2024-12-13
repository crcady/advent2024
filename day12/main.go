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

func (r region) sides() int {
	topfences := make(map[point]bool)
	bottomfences := make(map[point]bool)
	leftfences := make(map[point]bool)
	rightfences := make(map[point]bool)

	for p := range r.points {
		a := p.adjacent()
		up, down, left, right := a[0], a[1], a[2], a[3]
		if !r.points[up] {
			topfences[p] = true
		}
		if !r.points[down] {
			bottomfences[p] = true
		}
		if !r.points[left] {
			leftfences[p] = true
		}
		if !r.points[right] {
			rightfences[p] = true
		}
	}

	res := 0
	for len(topfences) > 0 {
		// This is an absurd way to do this
		var current point
		for k := range topfences {
			current = k
			break
		}
		// Remove it from our fences
		delete(topfences, current)

		// Move left
		left := point{current.x, current.y - 1}
		for topfences[left] {
			delete(topfences, left)
			left = point{left.x, left.y - 1}
		}

		// Move right
		right := point{current.x, current.y + 1}
		for topfences[right] {
			delete(topfences, right)
			right = point{right.x, right.y + 1}
		}

		// Catalog it
		res++
	}

	for len(bottomfences) > 0 {
		// This is an absurd way to do this
		var current point
		for k := range bottomfences {
			current = k
			break
		}
		// Remove it from our fences
		delete(bottomfences, current)

		// Move left
		left := point{current.x, current.y - 1}
		for bottomfences[left] {
			delete(bottomfences, left)
			left = point{left.x, left.y - 1}
		}

		// Move right
		right := point{current.x, current.y + 1}
		for bottomfences[right] {
			delete(bottomfences, right)
			right = point{right.x, right.y + 1}
		}

		// Catalog it
		res++
	}

	for len(leftfences) > 0 {
		// This is an absurd way to do this
		var current point
		for k := range leftfences {
			current = k
			break
		}
		// Remove it from our hfences
		delete(leftfences, current)

		// Move up
		up := point{current.x - 1, current.y}
		for leftfences[up] {
			delete(leftfences, up)
			up = point{up.x - 1, up.y}
		}

		// Move down
		down := point{current.x + 1, current.y}
		for leftfences[down] {
			delete(leftfences, down)
			down = point{down.x + 1, down.y}
		}

		// Catalog it
		res++
	}

	for len(rightfences) > 0 {
		// This is an absurd way to do this
		var current point
		for k := range rightfences {
			current = k
			break
		}
		// Remove it from our hfences
		delete(rightfences, current)

		// Move up
		up := point{current.x - 1, current.y}
		for rightfences[up] {
			delete(rightfences, up)
			up = point{up.x - 1, up.y}
		}

		// Move down
		down := point{current.x + 1, current.y}
		for rightfences[down] {
			delete(rightfences, down)
			down = point{down.x + 1, down.y}
		}

		// Catalog it
		res++
	}

	return res
}

func (r region) price() int {
	return r.perimeter() * r.area()
}

func (r region) price2() int {
	return r.sides() * r.area()
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

	ans2 := 0
	for _, r := range regions {
		ans2 += r.price2()
	}
	log.Printf("The price for part two is %d", ans2)

}
