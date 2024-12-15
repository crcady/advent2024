package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type point struct {
	x, y int
}

func (p point) left() point {
	return point{p.x - 1, p.y}
}

func (p point) right() point {
	return point{p.x + 1, p.y}
}

func (p point) up() point {
	return point{p.x, p.y - 1}
}

func (p point) down() point {
	return point{p.x, p.y + 1}
}

type warehouse struct {
	positions [][]byte
	robot     point
}

type doublehouse struct {
	positions [][]byte
	robot     point
}

func NewWarehouse() *warehouse {
	return &warehouse{[][]byte{}, point{-1, -1}}
}

func NewDoubleHouse() *doublehouse {
	return &doublehouse{[][]byte{}, point{-1, -1}}
}

func (wh *warehouse) addRow(r []byte) {
	newRow := make([]byte, len(r))
	copy(newRow, r)

	wh.positions = append(wh.positions, newRow)

	for i, c := range newRow {
		if c == '@' {
			wh.robot = point{i, len(wh.positions) - 1}
			wh.positions[len(wh.positions)-1][i] = '.'
		}
	}
}

func (dh *doublehouse) addRow(r []byte) {
	newRow := make([]byte, len(r)*2)
	for i, b := range r {
		switch b {
		case '#':
			newRow[2*i] = '#'
			newRow[2*i+1] = '#'
		case 'O':
			newRow[2*i] = '['
			newRow[2*i+1] = ']'
		case '.':
			newRow[2*i] = '.'
			newRow[2*i+1] = '.'
		case '@':
			newRow[2*i] = '.'
			newRow[2*i+1] = '.'
			dh.robot = point{2 * i, len(dh.positions)}
		}
	}

	dh.positions = append(dh.positions, newRow)
}

func (wh *warehouse) get(p point) byte {
	return wh.positions[p.y][p.x]
}

func (dh *doublehouse) get(p point) byte {
	return dh.positions[p.y][p.x]
}

func (wh *warehouse) set(p point, b byte) {
	wh.positions[p.y][p.x] = b
}

func (dh *doublehouse) set(p point, b byte) {
	dh.positions[p.y][p.x] = b
}

func (wh *warehouse) tryMove(dir rune) {
	var current point
	switch dir {
	case '>':
		current = wh.robot.right()
		for wh.get(current) == 'O' {
			current = current.right()
		}

		if wh.get(current) == '.' {
			wh.set(current, 'O')
			wh.robot = wh.robot.right()
			wh.set(wh.robot, '.')
		}

	case '<':
		current = wh.robot.left()
		for wh.get(current) == 'O' {
			current = current.left()
		}

		if wh.get(current) == '.' {
			wh.set(current, 'O')
			wh.robot = wh.robot.left()
			wh.set(wh.robot, '.')
		}

	case '^':
		current = wh.robot.up()
		for wh.get(current) == 'O' {
			current = current.up()
		}

		if wh.get(current) == '.' {
			wh.set(current, 'O')
			wh.robot = wh.robot.up()
			wh.set(wh.robot, '.')
		}

	case 'v':
		current = wh.robot.down()
		for wh.get(current) == 'O' {
			current = current.down()
		}

		if wh.get(current) == '.' {
			wh.set(current, 'O')
			wh.robot = wh.robot.down()
			wh.set(wh.robot, '.')
		}
	}
}

func (dh *doublehouse) tryMove(dir rune) {
	boxes := []point{}
	var current point
	switch dir {
	case '>':
		current = dh.robot.right()
		for dh.get(current) == '[' || dh.get(current) == ']' {
			boxes = append(boxes, current)
			current = current.right()
		}

		if dh.get(current) == '.' {
			dh.set(current, ']')
			for i, b := range boxes {
				if i%2 == 0 {
					dh.set(b, ']')
				} else {
					dh.set(b, '[')
				}
			}
			dh.robot = dh.robot.right()
			dh.set(dh.robot, '.')
		}

	case '<':
		current = dh.robot.left()
		for dh.get(current) == '[' || dh.get(current) == ']' {
			boxes = append(boxes, current)
			current = current.left()
		}

		if dh.get(current) == '.' {
			dh.set(current, '[')
			for i, b := range boxes {
				if i%2 == 0 {
					dh.set(b, '[')
				} else {
					dh.set(b, ']')
				}
			}
			dh.robot = dh.robot.left()
			dh.set(dh.robot, '.')
		}

	case '^':
		rows := []map[point]bool{}
		rows = append(rows, map[point]bool{dh.robot: true})
		keepGoing := true
		blocked := false
		for keepGoing {
			newRow := map[point]bool{}
			keepGoing = false
			for p := range rows[len(rows)-1] {
				switch dh.get(p.up()) {
				case '#':
					blocked = true
					keepGoing = false
					break
				case '[':
					newRow[p.up()] = true
					newRow[p.up().right()] = true
					keepGoing = true
				case ']':
					newRow[p.up()] = true
					newRow[p.up().left()] = true
					keepGoing = true
				}
			}

			rows = append(rows, newRow)
		}

		if !blocked {
			for i := len(rows) - 2; i > 0; i-- {
				for p := range rows[i] {
					dh.set(p.up(), dh.get(p))
					dh.set(p, '.')
				}
			}

			dh.robot = dh.robot.up()
			dh.set(dh.robot, '.')
		}

	case 'v':
		rows := []map[point]bool{}
		rows = append(rows, map[point]bool{dh.robot: true})
		keepGoing := true
		blocked := false
		for keepGoing {
			newRow := map[point]bool{}
			keepGoing = false
			for p := range rows[len(rows)-1] {
				switch dh.get(p.down()) {
				case '#':
					blocked = true
					keepGoing = false
					break
				case '[':
					newRow[p.down()] = true
					newRow[p.down().right()] = true
					keepGoing = true
				case ']':
					newRow[p.down()] = true
					newRow[p.down().left()] = true
					keepGoing = true
				}
			}

			rows = append(rows, newRow)
		}

		if !blocked {
			for i := len(rows) - 2; i > 0; i-- {
				for p := range rows[i] {
					dh.set(p.down(), dh.get(p))
					dh.set(p, '.')
				}
			}

			dh.robot = dh.robot.down()
			dh.set(dh.robot, '.')
		}

	}
}

func (wh *warehouse) GPS() int {
	res := 0
	for y, row := range wh.positions {
		for x, b := range row {
			if b == 'O' {
				res += (y * 100) + x
			}
		}
	}
	return res
}

func (dh *doublehouse) GPS() int {
	res := 0
	for y, row := range dh.positions {
		for x, b := range row {
			if b == '[' {
				res += (y * 100) + x
			}
		}
	}
	return res
}

func (wh *warehouse) print() {
	for _, row := range wh.positions {
		fmt.Println(string(row))
	}
}

func (dh *doublehouse) print() {
	for y, row := range dh.positions {
		for x, b := range row {
			if dh.robot.x == x && dh.robot.y == y {
				fmt.Print("@")
			} else {
				fmt.Print(string(b))
			}
		}
		fmt.Println()
	}
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
	wh := NewWarehouse()
	dh := NewDoubleHouse()

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			break
		}
		wh.addRow(line)
		dh.addRow(line)
	}

	moves := ""

	for scanner.Scan() {
		line := scanner.Text()
		moves = moves + line
	}

	//dh.print()

	for _, dir := range moves {
		wh.tryMove(dir)
		dh.tryMove(dir)
		//fmt.Println("After", string(dir))
		//dh.print()
		//fmt.Println()
	}

	//dh.print()

	log.Println(wh.GPS())
	log.Println(dh.GPS())
}
