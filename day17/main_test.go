package main

import (
	"fmt"
	"testing"
)

func Test_machine(t *testing.T) {
	var tests = []struct {
		a, b, c          int
		program          []int
		outA, outB, outC int
		outputs          []int
	}{
		{0, 0, 9, []int{2, 6}, 0, 1, 9, []int{}},
		{10, 0, 0, []int{5, 0, 5, 1, 5, 4}, 10, 0, 0, []int{0, 1, 2}},
		{2024, 0, 0, []int{0, 1, 5, 4, 3, 0}, 0, 0, 0, []int{4, 2, 5, 6, 7, 7, 7, 7, 3, 1, 0}},
		{0, 29, 0, []int{1, 7}, 0, 26, 0, []int{}},
		{0, 2024, 43690, []int{4, 0}, 0, 44354, 43690, []int{}},
	}

	for i, tt := range tests {
		testname := fmt.Sprintf("test%d", i)
		t.Run(testname, func(t *testing.T) {
			m := machine{tt.a, tt.b, tt.c, 0, tt.program, make([]int, 0)}
			for m.step() {
			}
			if len(tt.outputs) != len(m.outputs) {
				t.Fatal("Wrong number of outputs!")
			}

			for i := range tt.outputs {
				if tt.outputs[i] != m.outputs[i] {
					t.Errorf("Outputs don't match! %d != %d", tt.outputs[i], m.outputs[i])
				}
			}

			if tt.outA != m.regA {
				t.Errorf("Register A doesn't match")
			}

			if tt.outB != m.regB {
				t.Errorf("Register B doesn't match")
			}

			if tt.outC != m.regC {
				t.Errorf("Register C doesn't match")
			}

		})

	}
}
