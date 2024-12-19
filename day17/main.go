package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
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
	stepCount := 0

	for m.step() {
		stepCount++
	}

	log.Printf("Took %d steps and have %d outputs\n", stepCount, len(m.outputs))
	log.Println(m.output())
}
