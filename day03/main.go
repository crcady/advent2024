package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	input := "xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))"

	if len(os.Args) > 1 {
		fname := os.Args[1]
		dat, err := os.ReadFile(fname)
		if err != nil {
			log.Fatal(err)
		}
		input = string(dat)
	}

	r1, err := regexp.Compile(`mul\((\d\d?\d?),(\d\d?\d?)\)`)
	if err != nil {
		fmt.Println("Error compiling RegEx:", err)
	}

	matches := r1.FindAllStringSubmatch(input, -1)
	sum := 0

	for _, m := range matches {
		sum += handleMatch(m[1:])
	}

	fmt.Println("Answer to first half is:", sum)

	r2 := regexp.MustCompile(`(do)\(\)|(don't)\(\)|mul\((\d\d?\d?),(\d\d?\d?)\)`)
	if len(os.Args) < 2 {
		input = "xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))"
	}

	matches = r2.FindAllStringSubmatch(input, -1)
	enabled, sum := true, 0

	for _, m := range matches {
		switch {
		case m[1] == "do":
			enabled = true
		case m[2] == "don't":
			enabled = false
		default:
			if enabled {
				sum += handleMatch(m[3:])
			}
		}
	}

	fmt.Println("The answer to the second half is:", sum)

}

func handleMatch(m []string) int {
	x, err := strconv.Atoi(m[0])
	if err != nil {
		log.Fatal(err)
	}
	y, err := strconv.Atoi(m[1])
	if err != nil {
		log.Fatal(err)
	}

	return x * y
}
