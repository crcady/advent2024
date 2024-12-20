package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"

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

func newPoint(s string) point {
	re := regexp.MustCompile(`\d+`)
	match := re.FindAllString(s, -1)

	x, err := strconv.Atoi(match[0])
	if err != nil {
		log.Fatal("Failed parsing line", err)
	}
	y, err := strconv.Atoi(match[1])
	if err != nil {
		log.Fatal("Failed parsing line", err)
	}

	return point{x, y}
}

type wayfinder struct {
	corrupted map[point]bool
	width     int
}

func (w *wayfinder) open(p point) bool {
	if p.x < 0 || p.y < 0 || p.x >= w.width || p.y >= w.width {
		return false
	}

	return !w.corrupted[p]
}

func (w *wayfinder) solve1() int {
	openPoints := make(map[point]int64)
	for x := 0; x < w.width; x++ {
		for y := 0; y < w.width; y++ {
			p := point{x, y}
			if w.open(p) {
				openPoints[p] = int64(len(openPoints))
			}
		}
	}

	g := simple.NewUndirectedGraph()

	for p, pid := range openPoints {
		if rid, ok := openPoints[p.right()]; ok {
			g.SetEdge(g.NewEdge(simple.Node(pid), simple.Node(rid)))
		}

		if did, ok := openPoints[p.down()]; ok {
			g.SetEdge(g.NewEdge(simple.Node(pid), simple.Node(did)))
		}
	}

	pth := path.DijkstraFrom(g.Node(openPoints[point{0, 0}]), g)
	_, length := pth.To(openPoints[point{w.width - 1, w.width - 1}])

	return int(length)
}

func (w *wayfinder) solve2(points []point) point {
	for _, p := range points {
		w.corrupted[p] = true
		res := w.solve1()
		if res < 0 {
			return p
		}
	}
	log.Fatal("Didn't find a point that broke it...")
	return point{-1, -1}
}

func newWayfinder(points []point, width int) *wayfinder {
	corrupted := make(map[point]bool)
	for _, p := range points {
		corrupted[p] = true
	}
	return &wayfinder{corrupted, width}
}

func main() {
	fname := "example.txt"
	bytesToRead := 12
	memoryWidth := 7

	if len(os.Args) > 1 {
		fname = os.Args[1]
		bytesToRead = 1024
		memoryWidth = 71
	}

	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	points := make([]point, 0)

	for scanner.Scan() {
		points = append(points, newPoint(scanner.Text()))
	}

	wf := newWayfinder(points[:bytesToRead], memoryWidth)
	ans1 := wf.solve1()
	ans2 := wf.solve2(points[bytesToRead:])

	log.Println("Found a path throuth the corruption of length", ans1)
	log.Printf("All good until: %d,%d\n", ans2.x, ans2.y)

}
