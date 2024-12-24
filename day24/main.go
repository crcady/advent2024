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

}
