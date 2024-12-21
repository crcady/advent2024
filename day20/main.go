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

type cheat struct {
	start, end point
}

func (c cheat) length() int {
	dx := c.end.x - c.start.x
	dy := c.end.y - c.start.y

	if dx < 0 {
		dx = -dx
	}

	if dy < 0 {
		dy = -dy
	}

	return dx + dy

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

func (rt *racetrack) baseline() ([]point, int) {
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

	nodes, cost := pth.To(ids[rt.end])

	backmap := make(map[int64]point, len(ids))
	for p, n := range ids {
		backmap[n] = p
	}

	points := make([]point, len(nodes))
	for i, n := range nodes {
		points[i] = backmap[n.ID()]
	}

	return points, int(cost)
}

func (rt *racetrack) cheats(length int) map[int]int {
	savings := make(map[cheat]int)
	path, _ := rt.baseline()

	for i, first := range path[:len(path)-1] {
		// All cheats end on the path, or else the program seg faults
		for j := i + 1; j < len(path); j++ {
			second := path[j]

			current := cheat{first, second}
			cheatLength := current.length()
			if cheatLength > length {
				continue
			}

			stepsSkipped := j - i
			stepsSaved := stepsSkipped - cheatLength
			if stepsSaved > savings[current] {
				savings[current] = stepsSaved
			}
		}
	}

	freqs := make(map[int]int)
	for _, saved := range savings {
		freqs[saved]++
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

	_, baseline := rt.baseline()
	log.Println("The baseline is:", baseline)

	ans1 := 0
	for savings, count := range rt.cheats(2) {
		if savings >= 100 {
			ans1 += count
		}
	}
	log.Println(ans1)

	ans2 := 0
	for savings, count := range rt.cheats(20) {
		if savings >= 100 {
			ans2 += count
		}
	}

	log.Println(ans2)
}
