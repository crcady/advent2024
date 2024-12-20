package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/aclements/go-z3/z3"
)

const (
	op_adv = 0
	op_bxl = 1
	op_bst = 2
	op_jnz = 3
	op_bxc = 4
	op_out = 5
	op_bdv = 6
	op_cdv = 7
)

func mustConvert(v any) int {
	var s string
	switch v.(type) {
	case []byte:
		s = string(v.([]byte))
	case string:
		s = v.(string)
	}
	num, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal("Couldn't convert to int:", v)
	}
	return num
}

type machine struct {
	regA, regB, regC, ip int
	program              []int
	outputs              []int
}

func (m *machine) step() bool {
	if m.ip >= len(m.program) {
		return false
	}

	opcode := m.program[m.ip]
	operand := m.program[m.ip+1]
	m.ip += 2

	switch opcode {
	case op_adv:
		res := m.regA / (1 << m.combo(operand))
		m.regA = res
	case op_bdv:
		res := m.regA / (1 << m.combo(operand))
		m.regB = res
	case op_cdv:
		res := m.regA / (1 << m.combo(operand))
		m.regC = res
	case op_bxl:
		res := m.regB ^ operand
		m.regB = res
	case op_bst:
		res := m.combo(operand) % 8
		m.regB = res
	case op_jnz:
		if m.regA != 0 {
			m.ip = operand
		}
	case op_bxc:
		res := m.regB ^ m.regC
		m.regB = res
	case op_out:
		res := m.combo(operand) % 8
		m.outputs = append(m.outputs, res)
	}

	return true
}

func (m *machine) combo(op int) int {
	switch op {
	case 0, 1, 2, 3:
		return op
	case 4:
		return m.regA
	case 5:
		return m.regB
	case 6:
		return m.regC
	default:
		panic("Improper operand!")
	}
}

func (m *machine) output() string {
	outs := make([]string, len(m.outputs))
	for i, n := range m.outputs {
		outs[i] = strconv.Itoa(n)
	}

	return strings.Join(outs, ",")
}

func newMachine(data []byte) *machine {
	re := regexp.MustCompile(`Register A: (?P<regA>\d+)\nRegister B: (?P<regB>\d+)\nRegister C: (?P<regC>\d+)\n\nProgram: (?P<prog>[0-9,]+)`)

	matches := re.FindSubmatch(data)

	if matches == nil {
		log.Fatal("Failed to parse input", string(data))
	}

	regA := mustConvert(matches[re.SubexpIndex("regA")])
	regB := mustConvert(matches[re.SubexpIndex("regB")])
	regC := mustConvert(matches[re.SubexpIndex("regC")])

	pStrings := strings.Split(string(matches[re.SubexpIndex("prog")]), ",")

	prog := make([]int, len(pStrings))
	for i, s := range pStrings {
		prog[i] = mustConvert(s)
	}

	return &machine{
		regA:    regA,
		regB:    regB,
		regC:    regC,
		ip:      0,
		program: prog,
		outputs: make([]int, 0),
	}
}

type symSolver struct {
	regA, regB, regC z3.BV
	ip               int
	program          []int
	outputs          int
	ctx              *z3.Context
}

func newSymSolver(m *machine) *symSolver {
	ctx := z3.NewContext(nil)
	symA := ctx.BVConst("A", 64)
	symB := ctx.FromInt(int64(m.regB), ctx.BVSort(64)).(z3.BV)
	symC := ctx.FromInt(int64(m.regC), ctx.BVSort(64)).(z3.BV)

	pCopy := make([]int, len(m.program))
	copy(pCopy, m.program)

	return &symSolver{
		regA:    symA,
		regB:    symB,
		regC:    symC,
		ip:      0,
		program: pCopy,
		outputs: 0,
		ctx:     ctx,
	}
}

func (ss *symSolver) combo(n int) z3.BV {
	switch n {
	case 0, 1, 2, 3:
		return ss.literal(n)
	case 4:
		return ss.regA
	case 5:
		return ss.regB
	case 6:
		return ss.regC
	default:
		panic("Invalid Operand!")
	}
}

func (ss *symSolver) literal(n int) z3.BV {
	return ss.ctx.FromInt(int64(n), ss.ctx.BVSort(64)).(z3.BV)
}

func (ss *symSolver) solve() int {
	solver := z3.NewSolver(ss.ctx)

	targetValue := ss.regA

outer:
	for {
		opcode := ss.program[ss.ip]
		operand := ss.program[ss.ip+1]
		ss.ip += 2

		switch opcode {
		case op_adv:
			res := ss.regA.SRsh(ss.combo(operand))
			ss.regA = res
		case op_bdv:
			res := ss.regA.SRsh(ss.combo(operand))
			ss.regB = res
		case op_cdv:
			res := ss.regA.SRsh(ss.combo(operand))
			ss.regC = res
		case op_bxl:
			res := ss.regB.Xor(ss.literal(operand))
			ss.regB = res
		case op_bxc:
			res := ss.regB.Xor(ss.regC)
			ss.regB = res
		case op_bst:
			res := ss.combo(operand).SMod(ss.literal(8))
			ss.regB = res
		case op_jnz:
			if ss.outputs == len(ss.program) {
				solver.Assert(ss.regA.Eq(ss.literal(0)))
				break outer
			}

			solver.Assert(ss.regA.NE(ss.literal(0)))
			ss.ip = operand // This is always zero...
		case op_out:
			outVal := ss.combo(operand).SMod(ss.literal(8))
			solver.Assert(outVal.Eq(ss.literal(ss.program[ss.outputs])))
			ss.outputs++
		}
	}
	solver.Assert(targetValue.SGE(ss.literal(0)))
	sat, err := solver.Check()
	if !sat {
		log.Fatal("Failed to find an initial value for A")
	}
	lastFound := -1

	for sat && err == nil {
		model := solver.Model()
		solvedVal := model.Eval(targetValue, true)
		if val, _, ok := solvedVal.(z3.BV).AsInt64(); ok {
			lastFound = int(val)
			log.Println("Found a value:", val)
			solver.Assert(targetValue.SLT(ss.literal(int(val))))
			sat, err = solver.Check()
		} else {
			log.Fatal("Failed to retrieve A from model")
		}
	}

	if err != nil {
		log.Fatal("Error while solving constraints:", err)
	}

	return lastFound
}

func main() {
	fname := "example.txt"

	if len(os.Args) > 1 {
		fname = os.Args[1]
	}

	data, err := os.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}

	m := newMachine(data)
	ss := newSymSolver(m) // Safe; doesn't depend on m

	stepCount := 0

	for m.step() {
		stepCount++
	}

	log.Printf("Took %d steps and have %d outputs\n", stepCount, len(m.outputs))
	log.Println(m.output())

	ans2 := ss.solve()
	log.Println(ans2)
}
