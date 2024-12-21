package main

import (
	"bufio"
	"log"
	"os"

	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

type point struct {
	x, y int
}

func (p point) up() point {
	return point{p.x, p.y - 1}
}

func (p point) down() point {
	return point{p.x, p.y + 1}
}

func (p point) left() point {
	return point{p.x - 1, p.y}
}

func (p point) right() point {
	return point{p.x + 1, p.y}
}

type racetrack struct {
	course     map[point]bool
	walls      map[point]bool
	start, end point
	nextY      int
}

func newRacetrack() *racetrack {
	return &racetrack{
		course: make(map[point]bool),
		walls:  make(map[point]bool),
		start:  point{-1, -1},
		end:    point{-1, -1},
		nextY:  0,
	}
}

func (rt *racetrack) addLine(line string) {
	y := rt.nextY
	rt.nextY++

	for x, r := range line {
		p := point{x, y}
		switch r {
		case 'S':
			rt.start = p
			rt.course[p] = true
		case 'E':
			rt.end = p
			rt.course[p] = true
		case '.':
			rt.course[p] = true
		case '#':
			rt.walls[p] = true
		default:
			log.Fatal("Unrecognized symbol", r)
		}
	}
}

func (rt *racetrack) baseline() int {
	ids := make(map[point]int64)

	g := simple.NewUndirectedGraph()
	for p := range rt.course {
		ids[p] = int64(len(ids))
	}

	for p, pid := range ids {
		if rt.course[p.right()] {
			rid := ids[p.right()]
			g.SetEdge(g.NewEdge(simple.Node(pid), simple.Node(rid)))
		}

		if rt.course[p.down()] {
			did := ids[p.down()]
			g.SetEdge(g.NewEdge(simple.Node(pid), simple.Node(did)))
		}
	}

	pth := path.DijkstraFrom(simple.Node(ids[rt.start]), g)

	_, cost := pth.To(ids[rt.end])

	return int(cost)
}

func (rt *racetrack) cheat(p point) *racetrack {
	newCourse := make(map[point]bool, len(rt.course)+1)
	for pc := range rt.course {
		newCourse[pc] = true
	}
	newCourse[p] = true

	newWalls := make(map[point]bool, len(rt.walls)-1)

	for wc := range rt.walls {
		if wc != p {
			newWalls[wc] = true
		}
	}

	return &racetrack{
		course: newCourse,
		walls:  newWalls,
		start:  rt.start,
		end:    rt.end,
		nextY:  0,
	}
}

func (rt *racetrack) cheats() map[int]int {
	freqs := make(map[int]int)
	baseline := rt.baseline()

	for p := range rt.walls {
		cheatTrack := rt.cheat(p)
		cost := cheatTrack.baseline()
		if cost < baseline {
			freqs[baseline-cost]++
		}
	}

	return freqs
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

	rt := newRacetrack()

	for scanner.Scan() {
		line := scanner.Text()
		rt.addLine(line)
	}

	log.Println("The baseline is:", rt.baseline())
	ans1 := 0
	for savings, count := range rt.cheats() {
		if savings >= 100 {
			ans1 += count
		}
	}
	log.Println(ans1)

}
