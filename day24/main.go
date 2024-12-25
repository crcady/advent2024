package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	OP_AND   = "AND"
	OP_XOR   = "XOR"
	OP_OR    = "OR"
	OP_CONST = "EQ"

	STATUS_TRUE  = "1"
	STATUS_FALSE = "0"
	STATUS_UNK   = "Z"
)

type gate struct {
	id     string
	op     string
	status string

	input1 string
	input2 string
}

func (g gate) set() gate {
	return gate{
		id:     g.id,
		op:     g.op,
		status: STATUS_TRUE,

		input1: g.input1,
		input2: g.input2,
	}
}

func (g gate) unset() gate {
	return gate{
		id:     g.id,
		op:     g.op,
		status: STATUS_FALSE,

		input1: g.input1,
		input2: g.input2,
	}
}

func newConst(id, status string) gate {
	return gate{
		id:     id,
		op:     OP_CONST,
		status: status,

		input1: "",
		input2: "",
	}
}

func newOp(op, input1, input2, id string) gate {
	return gate{
		id:     id,
		op:     op,
		status: STATUS_UNK,

		input1: input1,
		input2: input2,
	}
}

type circuit map[string]gate

func newCircuit() circuit {
	return map[string]gate{}
}

func (c circuit) add(g gate) {
	c[g.id] = g
}

func (c circuit) settled() circuit {
	circ := newCircuit()
	for k, v := range c {
		circ[k] = v
	}

	for {
		unk_gates := []gate{}
		known_gates := map[string]string{}
		for id, g := range circ {
			if g.status == STATUS_UNK {
				unk_gates = append(unk_gates, g)
			} else {
				known_gates[id] = g.status
			}
		}

		if len(unk_gates) == 0 {
			break
		}

		for _, g := range unk_gates {
			var status1, status2 string

			if status, ok := known_gates[g.input1]; ok {
				status1 = status
			} else {
				continue
			}

			if status, ok := known_gates[g.input2]; ok {
				status2 = status
			} else {
				continue
			}

			switch g.op {
			case OP_XOR:
				if status1 != status2 {
					circ[g.id] = gate{g.id, g.op, STATUS_TRUE, g.input1, g.input2}
				} else {
					circ[g.id] = gate{g.id, g.op, STATUS_FALSE, g.input1, g.input2}
				}
			case OP_OR:
				if status1 == STATUS_TRUE || status2 == STATUS_TRUE {
					circ[g.id] = gate{g.id, g.op, STATUS_TRUE, g.input1, g.input2}
				} else {
					circ[g.id] = gate{g.id, g.op, STATUS_FALSE, g.input1, g.input2}
				}
			case OP_AND:
				if status1 == STATUS_TRUE && status2 == STATUS_TRUE {
					circ[g.id] = gate{g.id, g.op, STATUS_TRUE, g.input1, g.input2}
				} else {
					circ[g.id] = gate{g.id, g.op, STATUS_FALSE, g.input1, g.input2}
				}
			}
		}
	}

	return circ
}

func (c circuit) read(prefix string) int64 {
	bitstr := ""
	i := 0
	for {
		name := prefix + fmt.Sprintf("%02d", i)
		if g, ok := c[name]; ok {
			if g.status == STATUS_TRUE {
				bitstr = "1" + bitstr
			} else {
				bitstr = "0" + bitstr
			}
			i++
		} else {
			break
		}
	}
	log.Println(bitstr)
	res, err := strconv.ParseInt(bitstr, 2, 64)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func (c circuit) findOp(op, x, y string) (string, bool) {
	for id, g := range c {
		if g.op != op {
			continue
		}

		if g.input1 == x && g.input2 == y {
			return id, true
		}

		if g.input1 == y && g.input2 == x {
			return id, true
		}
	}

	return "", false
}

type adder struct {
	carryIn  string
	carryOut string

	inputXor string
	inputAnd string
	carryAnd string
}

func findHalfAdder(c circuit) adder {
	// The only half adder is the zero-th position
	a := adder{"", "", "", "", ""}

	if xorGate, ok := c.findOp(OP_XOR, "x00", "y00"); ok {
		a.inputXor = xorGate
	} else {
		log.Fatal("Need to re-map gates to find the half adder; didn't find XOR")
	}

	if carryOut, ok := c.findOp(OP_AND, "x00", "y00"); ok {
		a.carryOut = carryOut
	} else {
		log.Fatal("Need to re-map gates to find the half adder, didn't find AND")
	}

	return a
}

func findFullAdder(c circuit, idx int, carryHint string) adder {
	xInput := fmt.Sprintf("x%02d", idx)
	yInput := fmt.Sprintf("y%02d", idx)
	zOutput := fmt.Sprintf("z%02d", idx)

	inputXor, ok := c.findOp(OP_XOR, xInput, yInput)
	if !ok {
		log.Fatal("Need to do some swapping of the inputs to find XOR", xInput, yInput)
	}

	outputCandidate := c[zOutput]
	if outputCandidate.op != OP_XOR {
		outputCandidate2, ok := c.findOp(OP_XOR, carryHint, inputXor)
		if !ok {
			log.Fatal(zOutput, "gate is not a XOR and no substitute found")
		}
		log.Fatal("Need to swap ", zOutput, " and ", outputCandidate2)
	}

	if outputCandidate.input1 != carryHint && outputCandidate.input2 != carryHint {
		log.Fatal("Carry-in", carryHint, "not part of", zOutput)
	}

	if outputCandidate.input1 != inputXor && outputCandidate.input2 != inputXor {
		needToSwap := ""
		if outputCandidate.input1 == carryHint {
			needToSwap = outputCandidate.input2
		} else {
			needToSwap = outputCandidate.input1
		}
		log.Println("Output of ", xInput, " and ", yInput, " doesn't go to ", zOutput)
		log.Fatal("Try swapping ", needToSwap, " and ", inputXor)
	}

	inputAnd, ok := c.findOp(OP_AND, xInput, yInput)
	if !ok {
		log.Fatal("Need to do some swapping of input to find AND", xInput, yInput)
	}

	carryAnd, ok := c.findOp(OP_AND, inputXor, carryHint)
	if !ok {
		log.Fatal("Didn't find a carryAnd between", inputXor, "and", carryHint)
	}

	carryOut, ok := c.findOp(OP_OR, carryAnd, inputAnd)
	if !ok {
		log.Fatal("Didn't find a carrout OR gate between", carryAnd, "and", inputAnd)
	}

	return adder{
		carryIn:  carryHint,
		carryOut: carryOut,

		inputXor: inputXor,
		inputAnd: inputAnd,
		carryAnd: carryAnd,
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
	circ := newCircuit()

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		splits := strings.Split(line, ": ")
		id := splits[0]
		initial_condition := splits[1]
		circ.add(newConst(id, initial_condition))
	}

	for scanner.Scan() {
		line := scanner.Text()
		re := regexp.MustCompile(`([a-z0-9]+) ([A-Z]+) ([a-z0-9]+) -> ([a-z0-9]+)`)
		matches := re.FindAllStringSubmatch(line, -1)[0]
		input1 := matches[1]
		op := matches[2]
		input2 := matches[3]
		id := matches[4]

		circ.add(newOp(op, input1, input2, id))
	}

	solution1 := circ.settled()
	ans1 := solution1.read("z")

	log.Println("The answer to part one is", ans1)

	// Upon inspection, this is a n-bit ripple-carry adder. That has one half adder (z00 <- x00 XOR y00; carry00 <- x00 AND y00), and n-1 full adders.
	// Each full adder is composed of:
	// - A partial computation of zNN by computing xNN XOR yNN, possibly with input order reversed
	// - An XOR of that partial computation with the carry of the adder for the next-lowest bit
	//
	// There are a few ways that swaps can need to occur:
	// - If the wrong input bits are XOR'd together, we need to swap an input. To do that we have to figure out which swap to take
	// - If the z bits appear out of order, we need to swap those
	// - Probably the same thing can happen with the partials and the carries

	circ2 := circuit{}
	for k, v := range circ {
		circ2[k] = v
	}

	halfAdder := findHalfAdder(circ2)
	delete(circ2, halfAdder.inputXor)
	delete(circ2, "x00")
	delete(circ2, "y00")

	lastCarry := halfAdder.carryOut

	for i := 1; i < 45; i++ {
		fullerAdder := findFullAdder(circ2, i, lastCarry)
		lastCarry = fullerAdder.carryOut
		log.Println("Successful on iteration", i)
	}

}
