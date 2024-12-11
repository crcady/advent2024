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

type stoneGroup struct {
	ss     stones
	blinks int
}

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

func (ss stones) blink() stones {
	res := make(stones, 0, len(ss))
	for _, s := range ss {
		res = append(res, s.blink()...)
	}
	return res
}

func (sg stoneGroup) blink() stoneGroup {
	return stoneGroup{sg.ss.blink(), sg.blinks + 1}
}

func (sg stoneGroup) split() []stoneGroup {
	stones := sg.ss
	i := 0
	res := make([]stoneGroup, 0, 1000)

	for {
		j := i + 1000
		if j > len(stones) {
			j = len(stones)
			res = append(res, stoneGroup{stones[i:j], sg.blinks})
			break
		} else {
			res = append(res, stoneGroup{stones[i:j], sg.blinks})
			i = j
		}
	}

	return res
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
	for i := 0; i < 25; i++ {
		stones = stones.blink()
	}

	log.Println("After 25 blinks there are", len(stones), "stones")

	notDone := []stoneGroup{{stones, 25}}
	count := 0
	const threshold = 1_000_000

	for len(notDone) > 0 {
		sg := notDone[len(notDone)-1]
		notDone = notDone[:len(notDone)-1]
		sg = sg.blink()

		if sg.blinks == 75 {
			count += len(sg.ss)
			continue
		}

		if len(sg.ss) > threshold {
			notDone = append(notDone, sg.split()...)
		} else {
			notDone = append(notDone, sg)
		}
	}
	log.Println("After 75 blinks there are", count, "stones")
}
