package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

func check(patterns []string, design string, idx int, cache map[int]int) int {
	if res, ok := cache[idx]; ok {
		return res
	}
	//fmt.Println(idx)
	if idx == len(design) {
		cache[idx] = 1
		return 1
	}

	res := 0
outer:
	for _, p := range patterns {
		if idx+len(p) > len(design) {
			continue
		}

		for i, r := range p {
			if design[i+idx] != byte(r) {
				continue outer
			}
		}
		res += check(patterns, design, idx+len(p), cache)

	}

	cache[idx] = res
	return res
}

func main() {
	fname := "example.txt"

	if len(os.Args) > 1 {
		fname = os.Args[1]
	}

	file, err := os.Open(fname)
	if err != nil {
		log.Fatal("Couldn't open input", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	patterns := strings.Split(scanner.Text(), ", ")
	scanner.Scan() // Blank line

	designs := make([]string, 0)

	for scanner.Scan() {
		designs = append(designs, scanner.Text())

	}

	newRegEx := "^(" + strings.Join(patterns, "|") + ")*$"
	re := regexp.MustCompile(newRegEx)

	possible := make([]string, 0)

	for _, design := range designs {
		match := re.FindString(design)
		if len(match) > 0 {
			possible = append(possible, design)
		}
	}

	log.Println("Can make", len(possible), "designs")

	ans2 := 0

	for _, design := range possible {
		cache := make(map[int]int, len(design)) // This should guarantee no additional allocations occur for the memoizer
		ways := check(patterns, design, 0, cache)
		ans2 += ways
	}

	log.Println(ans2, "different arrangements")

}
