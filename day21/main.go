package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
)

var hits, misses int = 0, 0

type point struct {
	x, y int
}

func (p point) move(b byte) point {
	switch b {
	case '<':
		return point{p.x - 1, p.y}
	case '>':
		return point{p.x + 1, p.y}
	case '^':
		return point{p.x, p.y - 1}
	case 'v':
		return point{p.x, p.y + 1}
	default:
		log.Println("Move got a", string(b))
		panic("Invalid direction provided to move()")
	}
}

func newPoint(x, y int) point {
	return point{x, y}
}

type checker interface {
	check([]byte) int
}

type human struct {
}

func (h human) check(keys []byte) int {
	//fmt.Println("Human checking", string(keys))
	return len(keys)
}

type keyPad struct {
	keys      map[point]byte
	positions map[byte]point
	next      checker
	cache     map[string]int
}

// Checks whether a given point corresponds to a button on the pad
func (kp keyPad) valid(p point) bool {
	if _, ok := kp.keys[p]; ok {
		return true
	}
	return false
}

// Checks a set of moves to see if they would crash the robot
func (kp keyPad) checkMoves(startKey byte, moves []byte) bool {
	current := kp.positions[startKey]

	for _, m := range moves {
		current = current.move(m)
		if !kp.valid(current) {
			return false
		}
	}

	return true
}

func (kp keyPad) check(keys []byte) int {
	if res, ok := kp.cache[string(keys)]; ok {
		hits++
		return res
	} else {
		misses++
	}
	//fmt.Println("keyPad checking", string(keys))
	var fromKey byte = 'A'

	cum := 0

	for _, toKey := range keys {
		//fmt.Println("Looking at moves between", string(fromKey), "and", string(toKey))
		fromPoint := kp.positions[fromKey]
		toPoint := kp.positions[toKey]
		moves := map[byte]int{}

		if toPoint.x > fromPoint.x {
			moves['>'] = toPoint.x - fromPoint.x
		}

		if fromPoint.x > toPoint.x {
			moves['<'] = fromPoint.x - toPoint.x
		}

		if toPoint.y > fromPoint.y {
			moves['v'] = toPoint.y - fromPoint.y
		}

		if fromPoint.y > toPoint.y {
			moves['^'] = fromPoint.y - toPoint.y
		}

		perms := permute(moves)
		lowestCost := math.MaxInt

		for _, p := range perms {
			if kp.checkMoves(fromKey, p) {
				p = append(p, 'A')
				cost := kp.next.check(p)
				if cost < lowestCost {
					lowestCost = cost
				}
			}
		}
		cum += lowestCost
		fromKey = toKey
	}

	kp.cache[string(keys)] = cum
	return cum
}

func permute(counts map[byte]int) [][]byte {
	return recurPermute(counts, make([]byte, 0))
}

func recurPermute(counts map[byte]int, soFar []byte) [][]byte {
	res := make([][]byte, 0)

	for k, v := range counts {
		if v == 0 {
			continue
		}

		newCounts := make(map[byte]int, len(counts))
		for k1, v1 := range counts {
			newCounts[k1] = v1
		}
		newCounts[k]--

		newSoFar := make([]byte, len(soFar)+1)
		copy(newSoFar, soFar)
		newSoFar[len(newSoFar)-1] = k

		permutations := recurPermute(newCounts, newSoFar)
		res = append(res, permutations...)
	}

	if len(res) == 0 {
		res = append(res, soFar)
	}

	return res
}

func numPad(next checker) keyPad {
	keySlice := []byte("0123456789A")

	keys := map[point]byte{}
	positions := map[byte]point{}

	for _, b := range keySlice {
		var x, y int

		switch b {
		case '7', '4', '1':
			x = 0
		case '8', '5', '2', '0':
			x = 1
		case '9', '6', '3', 'A':
			x = 2
		}

		switch b {
		case '7', '8', '9':
			y = 0
		case '4', '5', '6':
			y = 1
		case '1', '2', '3':
			y = 2
		case '0', 'A':
			y = 3
		}

		p := newPoint(x, y)
		keys[p] = b
		positions[b] = p
	}

	return keyPad{
		keys:      keys,
		positions: positions,
		next:      next,
		cache:     make(map[string]int, 0),
	}
}

func dirPad(next checker) keyPad {
	keySlice := []byte("<>v^A")

	keys := map[point]byte{}
	positions := map[byte]point{}

	for _, b := range keySlice {
		var x, y int

		switch b {
		case '<':
			x = 0
		case '^', 'v':
			x = 1
		case 'A', '>':
			x = 2
		}

		switch b {
		case '^', 'A':
			y = 0
		case '<', 'v', '>':
			y = 1
		}

		p := newPoint(x, y)

		keys[p] = b
		positions[b] = p

	}

	return keyPad{
		keys:      keys,
		positions: positions,
		next:      next,
		cache:     make(map[string]int, 0),
	}
}

func main() {
	fname := "example.txt"

	if len(os.Args) > 1 {
		fname = os.Args[1]
	}

	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	ans1 := 0
	ans2 := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()

		h := human{}
		dp2 := dirPad(h)
		dp1 := dirPad(dp2)
		np := numPad(dp1)

		cost := np.check(line)

		num, err := strconv.Atoi(string(line[:len(line)-1]))
		if err != nil {
			panic(err)
		}

		ans1 += cost * num

		var lastChecker checker = h
		for i := 0; i < 25; i++ {
			newChecker := dirPad(lastChecker)
			lastChecker = newChecker
		}
		np2 := numPad(lastChecker)
		cost2 := np2.check(line)
		ans2 += cost2 * num

		log.Println("Finished", string(line))
	}

	log.Println("Computed", ans1, "and", ans2, "with", hits, "cache hits and", misses, "cache misses")
}
