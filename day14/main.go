package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type point struct {
	x, y int
}

type robot struct {
	pos, vel point
}

func newRobot(line string) *robot {
	re := regexp.MustCompile(`-?\d+`)
	matches := re.FindAllString(line, -1)
	if len(matches) != 4 {
		log.Fatal("Didn't parse the input line properly", line, matches)
	}

	nums := make([]int, 4)
	for i, s := range matches {
		num, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal("Couldn't parse into number:", num, err)
		}
		nums[i] = num
	}

	return &robot{
		point{nums[0], nums[1]},
		point{nums[2], nums[3]},
	}
}

func (r *robot) tick(height, width int) {
	r.pos.x = (r.pos.x + r.vel.x) % width
	r.pos.y = (r.pos.y + r.vel.y) % height

	if r.pos.x < 0 {
		r.pos.x = width + r.pos.x
	}

	if r.pos.y < 0 {
		r.pos.y = height + r.pos.y
	}
}

type board struct {
	robots        []*robot
	height, width int
}

func newBoard(h, w int) *board {
	return &board{[]*robot{}, h, w}
}

func (b *board) addRobot(r *robot) {
	b.robots = append(b.robots, r)
}

func (b *board) tick() {
	for _, r := range b.robots {
		r.tick(b.height, b.width)
	}
}

func (b *board) print() {
	counts := map[point]int{}

	for _, r := range b.robots {
		counts[r.pos]++
	}

	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			k := counts[point{j, i}]
			if k > 0 {
				fmt.Print(counts[point{j, i}])
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

func (bp *board) quadCounts() [4]int {
	// a b
	// c d
	a, b, c, d := 0, 0, 0, 0

	for _, r := range bp.robots {
		if r.pos.x < (bp.width-1)/2 { // left
			if r.pos.y < (bp.height-1)/2 { // top
				a++
			} else if r.pos.y > (bp.height-1)/2 { // bottom
				c++
			}
		} else if r.pos.x > (bp.width-1)/2 { //right
			if r.pos.y < (bp.height-1)/2 { // top
				b++
			} else if r.pos.y > (bp.height-1)/2 { // bottom
				d++
			}
		}
	}

	return [4]int{a, b, c, d}
}

func (b *board) safetyFactor() int {
	qc := b.quadCounts()
	return qc[0] * qc[1] * qc[2] * qc[3]
}

func findTriangle(pts map[point]bool, tip point, size int) bool {
	// A triangle of size 1 is a point
	// A triangle of size n has n-squared points

	for i := 0; i < size; i++ {
		for j := 0; j <= i; j++ {
			// Check left
			if !pts[point{tip.x - j, tip.y + i}] {
				return false
			}
			// Check right
			if !pts[point{tip.x + j, tip.y + i}] {
				return false
			}
		}
	}
	return true
}

func (b *board) isTree() bool {

	points := map[point]bool{}
	for _, r := range b.robots {
		points[r.pos] = true
	}

	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			if findTriangle(points, point{x, y}, 3) {
				return true
			}
		}
	}
	return false
}

func main() {
	fname := "example.txt"
	height, width := 7, 11

	if len(os.Args) > 1 {
		fname = os.Args[1]
		height, width = 103, 101
	}

	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	board := newBoard(height, width)
	board2 := newBoard(height, width)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		r := newRobot(line)
		board.addRobot(r)
		r2 := newRobot(line)
		board2.addRobot(r2)

	}

	for i := 0; i < 100; i++ {
		board.tick()
	}

	log.Println(board.safetyFactor())

	ticks := 0
	for !board2.isTree() {
		board2.tick()
		ticks++

		if ticks > 10_000 {
			log.Println("Failed!")
			break
		}
	}

	board2.print()
	log.Println(ticks)

}
