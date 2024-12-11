package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type stone int
type stones []stone

func (s stone) blink() stones {
	if s == 0 {
		return []stone{1}
	}

	if digits := strconv.Itoa(int(s)); len(digits)%2 == 0 {
		left := digits[:len(digits)/2]

		right := digits[len(digits)/2:]
		return []stone{newStone(left), newStone(right)}
	}

	return []stone{s * 2024}
}

func newStone(s string) stone {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return stone(num)
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
	scanner.Scan()
	line := scanner.Text()

	var stones stones
	for _, s := range strings.Split(line, " ") {
		stones = append(stones, newStone(s))
	}

	log.Println("There are", len(stones), "stones. About to blink...")

	stone_map := make(map[stone]int)
	for _, s := range stones {
		stone_map[s] = stone_map[s] + 1
	}
	for i := 0; i < 25; i++ {
		temp_map := make(map[stone]int)
		for k, v := range stone_map {
			temp_map[k] = v
		}

		stone_map = make(map[stone]int)

		for k, v := range temp_map {
			ss := k.blink()
			for _, s := range ss {
				stone_map[s] = stone_map[s] + v
			}
		}
	}
	ans1 := 0
	for _, v := range stone_map {
		ans1 += v
	}
	log.Println("After 25 blinks there are", len(stone_map), "unique stones and", ans1, "total stones")

	for i := 0; i < 50; i++ {
		temp_map := make(map[stone]int)
		for k, v := range stone_map {
			temp_map[k] = v
		}

		stone_map = make(map[stone]int)

		for k, v := range temp_map {
			ss := k.blink()
			for _, s := range ss {
				stone_map[s] = stone_map[s] + v
			}
		}
	}
	ans2 := 0
	for _, v := range stone_map {
		ans2 += v
	}
	log.Println("After 75 blinks there are", len(stone_map), "unique stones and", ans2, "total stones")
}
