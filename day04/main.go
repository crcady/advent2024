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

type board []string

func (b board) getPoint(p point) byte {
	x, y := p.x, p.y
	height := len(b)
	var width int
	if x <= height-1 && x >= 0 {
		width = len(b[x])
	} else {
		return '.'
	}

	if y < 0 || y > width-1 {
		return '.'
	}

	return b[x][y]
}

func main() {
	var fname string
	var lines board

	if len(os.Args) > 1 {
		fname = os.Args[1]
	} else {
		fname = "example.txt"
	}

	file, err := os.Open(fname)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	fmt.Printf("Read %d lines from %s\n", len(lines), fname)

	var xlocs []point

	for i, l := range lines {
		for j, c := range l {
			if c == 'X' {
				xlocs = append(xlocs, point{i, j})
			}
		}
	}

	ans1 := 0
	for _, loc := range xlocs {
		ans1 += checkFromX(lines, loc)
	}

	fmt.Printf("Found %d XMASes\n", ans1)

	var alocs []point

	for i, l := range lines {
		for j, c := range l {
			if c == 'A' {
				alocs = append(alocs, point{i, j})
			}
		}
	}

	ans2 := 0

	for _, loc := range alocs {
		if checkA(lines, loc) {
			ans2++
		}
	}

	fmt.Printf("Found %d X-MASes\n", ans2)

}

func checkFromX(b board, p point) int {
	x, y := p.x, p.y
	count := 0

	if b.getPoint(point{x - 1, y}) == 'M' && b.getPoint(point{x - 2, y}) == 'A' && b.getPoint(point{x - 3, y}) == 'S' {
		count++
	}

	if b.getPoint(point{x - 1, y - 1}) == 'M' && b.getPoint(point{x - 2, y - 2}) == 'A' && b.getPoint(point{x - 3, y - 3}) == 'S' {
		count++
	}

	if b.getPoint(point{x - 1, y + 1}) == 'M' && b.getPoint(point{x - 2, y + 2}) == 'A' && b.getPoint(point{x - 3, y + 3}) == 'S' {
		count++
	}

	if b.getPoint(point{x, y - 1}) == 'M' && b.getPoint(point{x, y - 2}) == 'A' && b.getPoint(point{x, y - 3}) == 'S' {
		count++
	}

	if b.getPoint(point{x, y + 1}) == 'M' && b.getPoint(point{x, y + 2}) == 'A' && b.getPoint(point{x, y + 3}) == 'S' {
		count++
	}

	if b.getPoint(point{x + 1, y}) == 'M' && b.getPoint(point{x + 2, y}) == 'A' && b.getPoint(point{x + 3, y}) == 'S' {
		count++
	}

	if b.getPoint(point{x + 1, y - 1}) == 'M' && b.getPoint(point{x + 2, y - 2}) == 'A' && b.getPoint(point{x + 3, y - 3}) == 'S' {
		count++
	}

	if b.getPoint(point{x + 1, y + 1}) == 'M' && b.getPoint(point{x + 2, y + 2}) == 'A' && b.getPoint(point{x + 3, y + 3}) == 'S' {
		count++
	}

	return count
}

func checkA(b board, p point) bool {
	x, y := p.x, p.y

	ul, ur, bl, br := b.getPoint(point{x - 1, y - 1}), b.getPoint(point{x - 1, y + 1}), b.getPoint(point{x + 1, y - 1}), b.getPoint(point{x + 1, y + 1})

	if (ul == 'M' && br == 'S') || (ul == 'S' && br == 'M') {
		if (ur == 'M' && bl == 'S') || (ur == 'S' && bl == 'M') {
			return true
		}
	}

	return false
}
