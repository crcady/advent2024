package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
)

func mustConvert(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return num
}

type clawMachine struct {
	ax, ay, bx, by, px, py int
}

func clawMachineFromLines(ls []string) clawMachine {
	re := regexp.MustCompile(`\d+`)
	a := re.FindAllString(ls[0], 2)
	b := re.FindAllString(ls[1], 2)
	p := re.FindAllString(ls[2], 2)

	return clawMachine{mustConvert(a[0]), mustConvert(a[1]), mustConvert(b[0]), mustConvert(b[1]), mustConvert(p[0]), mustConvert(p[1])}
}

func (cm clawMachine) solve() int {
	x, y := 0, 0
	cost := 0
	for {
		switch {
		case x > cm.px || y > cm.py:
			return 0
		case (cm.px-x)%cm.bx == 0 && (cm.py-y)%cm.by == 0 && (cm.px-x)/cm.bx == (cm.py-y)/cm.by && (cm.px-x)/cm.bx <= 100:
			return cost + (cm.py-y)/cm.by
		default:
			cost += 3
			x += cm.ax
			y += cm.ay
			if cost > 300 {
				return 0
			}
		}
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

	lines := []string{}
	machines := []clawMachine{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			m := clawMachineFromLines(lines)
			machines = append(machines, m)

			lines = []string{}
		} else {
			lines = append(lines, line)
		}
	}
	m := clawMachineFromLines(lines)
	machines = append(machines, m)

	ans1 := 0
	for _, cm := range machines {
		ans1 += cm.solve()
	}
	log.Println(ans1)

}
