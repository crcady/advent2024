package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/aclements/go-z3/z3"
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

func makeInt(ctx *z3.Context, num int) z3.Int {
	val := ctx.FromInt(int64(num), ctx.IntSort())
	return val.(z3.Int)
}

func (cm clawMachine) solve2() int {
	ctx := z3.NewContext(nil)

	A := ctx.IntConst("A")
	B := ctx.IntConst("B")
	Cost := ctx.IntConst("Cost")

	Ax := makeInt(ctx, cm.ax)
	Ay := makeInt(ctx, cm.ay)
	Bx := makeInt(ctx, cm.bx)
	By := makeInt(ctx, cm.by)

	Px := makeInt(ctx, cm.px)
	Py := makeInt(ctx, cm.py)

	x1 := A.Mul(Ax)
	x2 := B.Mul(Bx)
	xSum := x1.Add(x2)

	y1 := A.Mul(Ay)
	y2 := B.Mul(By)
	ySum := y1.Add(y2)

	// Now that they are all declared, we need to tell the solver we care about them
	solver := z3.NewSolver(ctx)
	solver.Assert(Px.Eq(xSum))
	solver.Assert(Py.Eq(ySum))
	solver.Assert(Cost.Eq(B.Add(A.Mul(makeInt(ctx, 3)))))

	// Now the expensive part
	sat, err := solver.Check()

	if err != nil {
		log.Fatal("Failed to check constraints", err)
	}

	if sat {
		m := solver.Model()
		solvedVal := m.Eval(Cost, true)
		if val, _, ok := solvedVal.(z3.Int).AsInt64(); ok {
			return int(val)
		} else {
			log.Fatal("Failed to retrieve Cost from model")
		}
	}

	return 0
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
	machines2 := []clawMachine{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			m := clawMachineFromLines(lines)
			machines = append(machines, m)

			m2 := clawMachineFromLines(lines)
			m2.px += 10000000000000
			m2.py += 10000000000000
			machines2 = append(machines2, m2)

			lines = []string{}
		} else {
			lines = append(lines, line)
		}
	}
	m := clawMachineFromLines(lines)
	machines = append(machines, m)

	m2 := clawMachineFromLines(lines)
	m2.px += 10000000000000
	m2.py += 10000000000000
	machines2 = append(machines2, m2)

	ans1 := 0
	for _, cm := range machines {
		ans1 += cm.solve()
	}
	log.Println(ans1)

	ans2 := 0
	for _, cm := range machines2 {
		ans2 += cm.solve2()
	}
	log.Println(ans2)

}
