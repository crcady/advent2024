package main

import (
	"bufio"
	"log"
	"os"
)

type schematic struct {
	isLock  bool
	heights []int
}

func newSchematic(lines []string) schematic {
	numCols := len(lines[0])
	if numCols != 5 {
		log.Fatal("Wrong number of cylinders", lines[0])
	}

	heights := make([]int, numCols)
	for _, line := range lines {
		for i, x := range line {
			if x == '#' {
				heights[i]++
			}
		}
	}

	//All of the heights are one too many, because the filled row doesn't count
	for i := range heights {
		heights[i]--
	}

	isLock := lines[0][0] == '#'

	return schematic{isLock, heights}
}

func checkFit(lock, key schematic) bool {
	for i := range lock.heights {
		if lock.heights[i]+key.heights[i] > 5 {
			return false
		}
	}
	return true
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

	currentLines := []string{}
	schematics := []schematic{}

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			schematics = append(schematics, newSchematic(currentLines))
			currentLines = []string{}
		} else {
			currentLines = append(currentLines, line)
		}
	}

	// There's no blank line at the end of the file
	schematics = append(schematics, newSchematic(currentLines))

	log.Printf("There are %d different schematics\n", len(schematics))
	locks := []schematic{}
	keys := []schematic{}

	for _, s := range schematics {
		if s.isLock {
			locks = append(locks, s)
		} else {
			keys = append(keys, s)
		}
	}

	log.Printf("There are %d locks and %d keys\n", len(locks), len(keys))

	ans1 := 0
	for _, lock := range locks {
		for _, key := range keys {
			if checkFit(lock, key) {
				ans1++
			}
		}
	}

	log.Printf("There are %d unique matches\n", ans1)
}
