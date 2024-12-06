package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type point struct {
	x int
	y int
}

type state struct {
	obstacles []point
	height    int
	width     int
}

type guard struct {
	loc       point
	direction string
}

func (s state) onMap(g guard) bool {
	return g.loc.x >= 0 && g.loc.y >= 0 && g.loc.x < s.height && g.loc.y < s.width
}

func (s state) step(g guard) guard {
	var next point
	switch g.direction {
	case "left":
		next = point{g.loc.x, g.loc.y - 1}
	case "right":
		next = point{g.loc.x, g.loc.y + 1}
	case "up":
		next = point{g.loc.x - 1, g.loc.y}
	case "down":
		next = point{g.loc.x + 1, g.loc.y}
	default:
		log.Fatal("Unknown direction")
	}

	collision := false

	for _, p := range s.obstacles {
		if next == p {
			collision = true
		}
	}

	if collision {
		return guard{
			loc:       g.loc,
			direction: map[string]string{"left": "up", "up": "right", "right": "down", "down": "left"}[g.direction],
		}
	} else {
		return guard{
			loc:       next,
			direction: g.direction,
		}
	}
}

func main() {
	fname := "example.txt"

	if len(os.Args) == 2 {
		fname = os.Args[1]
	}

	file, err := os.Open(fname)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	x := 0
	obstacles := make([]point, 0)
	candidates := make([]point, 0)

	var origin point
	var l string

	for scanner.Scan() {
		l = scanner.Text()
		for y, c := range l {
			switch c {
			case '#':
				obstacles = append(obstacles, point{x, y})
			case '^':
				origin = point{x, y}
			case '.':
				candidates = append(candidates, point{x, y})
			}
		}
		x++
	}

	s0 := state{
		obstacles: obstacles,
		height:    x,
		width:     len(l),
	}

	g := guard{
		loc:       origin,
		direction: "up",
	}

	visited := make(map[point]bool)
	visited[origin] = true

	for s0.onMap(g) {
		g = s0.step(g)
		visited[g.loc] = true
	}

	fmt.Println(len(visited) - 1) //includes the off-the-board position

	count := 0

	fmt.Println("Brute forcing", len(candidates), "candidate positions")

	for _, c := range candidates {
		obs := make([]point, len(obstacles)+1)
		copy(obs, obstacles)
		obs[len(obstacles)] = c

		g := guard{origin, "up"}
		s := state{obs, s0.height, s0.width}

		seen := make(map[guard]bool)
		for {
			if seen[g] {
				count++
				break
			}

			if !s.onMap(g) {
				break
			}

			seen[g] = true
			g = s.step(g)
		}
	}

	fmt.Println(count)

}
