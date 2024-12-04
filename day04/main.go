package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	var fname string
	var lines []string

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

	var xlocs [][]int

	for i, l := range lines {
		for j, c := range l {
			if c == 'X' {
				xlocs = append(xlocs, []int{i, j})
			}
		}
	}

	ans1 := 0
	for _, loc := range xlocs {
		ans1 += checkFromX(lines, loc[0], loc[1])
	}

	fmt.Printf("Found %d XMASes\n", ans1)

	var alocs [][]int

	for i, l := range lines {
		for j, c := range l {
			if c == 'A' {
				alocs = append(alocs, []int{i, j})
			}
		}
	}

	ans2 := 0

	for _, loc := range alocs {
		if checkA(lines, loc[0], loc[1]) {
			ans2++
		}
	}

	fmt.Printf("Found %d X-MASes\n", ans2)

}

func checkFromX(lines []string, x int, y int) int {
	count := 0

	if x >= 3 {
		if lines[x-1][y] == 'M' && lines[x-2][y] == 'A' && lines[x-3][y] == 'S' {
			count++
		}

		if y >= 3 {
			if lines[x-1][y-1] == 'M' && lines[x-2][y-2] == 'A' && lines[x-3][y-3] == 'S' {
				count++
			}
		}

		if y <= len(lines)-4 {
			if lines[x-1][y+1] == 'M' && lines[x-2][y+2] == 'A' && lines[x-3][y+3] == 'S' {
				count++
			}
		}

	}

	if x <= len(lines[x])-4 {
		if lines[x+1][y] == 'M' && lines[x+2][y] == 'A' && lines[x+3][y] == 'S' {
			count++
		}

		if y >= 3 {
			if lines[x+1][y-1] == 'M' && lines[x+2][y-2] == 'A' && lines[x+3][y-3] == 'S' {
				count++
			}
		}

		if y <= len(lines)-4 {
			if lines[x+1][y+1] == 'M' && lines[x+2][y+2] == 'A' && lines[x+3][y+3] == 'S' {
				count++
			}
		}

	}

	if y >= 3 {
		if lines[x][y-1] == 'M' && lines[x][y-2] == 'A' && lines[x][y-3] == 'S' {
			count++
		}
	}

	if y <= len(lines)-4 {
		if lines[x][y+1] == 'M' && lines[x][y+2] == 'A' && lines[x][y+3] == 'S' {
			count++
		}
	}

	return count
}

func checkA(lines []string, x int, y int) bool {
	if x == 0 || y == 0 || x > len(lines)-2 || y > len(lines[x])-2 {
		return false
	}

	a, b, c, d := lines[x-1][y-1], lines[x-1][y+1], lines[x+1][y-1], lines[x+1][y+1]

	if (a == 'M' && d == 'S') || (a == 'S' && d == 'M') {
		if (b == 'M' && c == 'S') || (b == 'S' && c == 'M') {
			return true
		}
	}

	return false
}
