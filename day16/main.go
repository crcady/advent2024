package main

import (
	"bufio"
	"log"
	"math"
	"os"

	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

type reindeer struct {
	x, y      int
	direction string
}

func (r reindeer) cw() reindeer {
	dmap := map[string]string{"up": "right", "right": "down", "down": "left", "left": "up"}
	return reindeer{r.x, r.y, dmap[r.direction]}
}

func (r reindeer) ccw() reindeer {
	dmap := map[string]string{"up": "left", "left": "down", "down": "right", "right": "up"}
	return reindeer{r.x, r.y, dmap[r.direction]}
}

func (r reindeer) point() point {
	return point{r.x, r.y}
}

func (r reindeer) fwd() reindeer {
	switch r.direction {
	case "up":
		return reindeer{r.x, r.y - 1, r.direction}
	case "down":
		return reindeer{r.x, r.y + 1, r.direction}
	case "left":
		return reindeer{r.x - 1, r.y, r.direction}
	case "right":
		return reindeer{r.x + 1, r.y, r.direction}
	default:
		log.Fatal("Invalid direction when moving forward")
	}
	return r // Just to make the compiler happy
}

type point struct {
	x, y int
}

type maze struct {
	startX, startY, endX, endY int
	linesSeen                  int
	tiles                      map[point]bool
}

func newMaze() *maze {
	return &maze{0, 0, 0, 0, 0, make(map[point]bool)}
}

func (m *maze) addLine(line string) {
	y := m.linesSeen

	for x, c := range line {
		switch c {
		case 'S':
			m.startX = x
			m.startY = y
			m.tiles[point{x, y}] = true
		case 'E':
			m.endX = x
			m.endY = y
			m.tiles[point{x, y}] = true
		case '.':
			m.tiles[point{x, y}] = true
		}
	}
	m.linesSeen++
}

func (m *maze) open(p point) bool {
	return m.tiles[p]
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

	m := newMaze()

	for scanner.Scan() {
		line := scanner.Text()
		m.addLine(line)
	}

	rIDs := make(map[reindeer]int64)
	rIDs[reindeer{m.startX, m.startY, "right"}] = 0

	for p := range m.tiles {
		if p.x == m.endX && p.y == m.endY {
			continue
		}

		rUp := reindeer{p.x, p.y, "up"}
		rDown := reindeer{p.x, p.y, "down"}
		rLeft := reindeer{p.x, p.y, "left"}
		rRight := reindeer{p.x, p.y, "right"}

		dirs := []reindeer{rUp, rDown, rLeft, rRight}

		if p.x == m.startX && p.y == m.startY {
			dirs = dirs[:3]
		}

		for _, r := range dirs {
			rIDs[r] = int64(len(rIDs)) + 1
		}
	}

	rIDs[reindeer{m.endX, m.endY, "up"}] = -1
	rIDs[reindeer{m.endX, m.endY, "down"}] = -1
	rIDs[reindeer{m.endX, m.endY, "left"}] = -1
	rIDs[reindeer{m.endX, m.endY, "right"}] = -1

	g := simple.NewWeightedDirectedGraph(math.Inf(1), math.Inf(1))

	for r := range rIDs {
		if r.x == m.endX && r.y == m.endY {
			continue
		}

		cw := r.cw()
		ccw := r.ccw()
		fwd := r.fwd()

		//log.Println("Adding edge from:", r, rIDs[r], "to", cw, rIDs[cw])
		g.SetWeightedEdge(simple.WeightedEdge{F: simple.Node(rIDs[r]), T: simple.Node(rIDs[cw]), W: 1000})
		g.SetWeightedEdge(simple.WeightedEdge{F: simple.Node(rIDs[r]), T: simple.Node(rIDs[ccw]), W: 1000})

		if m.open(fwd.point()) {
			g.SetWeightedEdge(simple.WeightedEdge{F: simple.Node(rIDs[r]), T: simple.Node(rIDs[fwd]), W: 1})
		}
	}

	pth := path.DijkstraFrom(simple.Node(0), g)

	_, cost := pth.To(-1)

	log.Println(cost)

	pths := path.DijkstraAllFrom(simple.Node(0), g)
	paths, _ := pths.AllTo(-1)

	goodSeats := make(map[point]bool)
	pointMap := make(map[int64]point)
	for r, i := range rIDs {
		pointMap[i] = r.point()
	}

	for _, path := range paths {
		for _, n := range path {
			goodSeats[pointMap[n.ID()]] = true
		}
	}

	log.Println(len(goodSeats))

}
