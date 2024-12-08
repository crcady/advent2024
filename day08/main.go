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

func (p1 point) antinodes(p2 point) []point {
	dx := p2.x - p1.x
	dy := p2.y - p1.y

	a1 := point{p1.x - dx, p1.y - dy}
	a2 := point{p2.x + dx, p2.y + dy}

	return []point{a1, a2}
}

func (p point) check(height int, width int) bool {
	return p.x >= 0 && p.y >= 0 && p.x < height && p.y < width
}

func (p1 point) antinodes2(p2 point, height int, width int) []point {
	ans := make([]point, 0, 2)
	dx := p2.x - p1.x
	dy := p2.y - p1.y

	p := p1
	for p.check(height, width) {
		ans = append(ans, p)
		p = point{p.x - dx, p.y - dy}
	}

	p = point{p1.x + dx, p1.y + dy}
	for p.check(height, width) {
		ans = append(ans, p)
		p = point{p.x + dx, p.y + dy}
	}

	p = p2
	for p.check(height, width) {
		ans = append(ans, p)
		p = point{p.x + dx, p.y + dy}
	}

	return ans

}

func main() {
	fname := "example.txt"
	antennas := make(map[rune][]point)
	antinodes := make(map[point]bool)
	antinodes2 := make(map[point]bool)

	if len(os.Args) > 1 {
		fname = os.Args[1]
	}

	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.Printf("Processing %s...\n", fname)
	scanner := bufio.NewScanner(file)
	x := 0
	var l string
	for scanner.Scan() {
		l = scanner.Text()
		for y, b := range l {
			if b != '.' && b != '#' {
				if a, ok := antennas[b]; ok {
					antennas[b] = append(a, point{x, y})
				} else {
					antennas[b] = []point{{x, y}}
				}
			}
		}
		x++
	}
	height, width := x, len(l)

	for _, v := range antennas {

		for i, p1 := range v {
			for _, p2 := range v[i+1:] {
				candidates := p1.antinodes(p2)
				for _, c := range candidates {
					if c.check(height, width) {
						antinodes[c] = true
					}
				}
				candidates2 := p1.antinodes2(p2, height, width)
				for _, c := range candidates2 {
					antinodes2[c] = true
				}
			}
		}
	}
	fmt.Println(len(antinodes))
	fmt.Println(len(antinodes2))

	// for i := range make([]bool, height) {
	// 	for j := range make([]bool, width) {
	// 		if antinodes2[point{i, j}] {
	// 			fmt.Print("#")
	// 		} else {
	// 			fmt.Print(".")
	// 		}
	// 	}
	// 	fmt.Print("\n")
	// }
}
