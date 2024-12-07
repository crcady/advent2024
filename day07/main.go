package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type eqTarget struct {
	target   int
	operands []int
}

func (et eqTarget) check(useConcat bool) bool {
	if sum(et.operands) == et.target {
		return true
	}

	if len(et.operands) < 2 {
		return false
	}

	if et.operands[0] > et.target {
		return false
	}

	lastNum := et.operands[len(et.operands)-1]
	if et.target%lastNum == 0 {
		timesTarget := eqTarget{et.target / lastNum, et.operands[:len(et.operands)-1]}
		if timesTarget.check(useConcat) {
			return true
		}
	}

	plusTarget := eqTarget{et.target - lastNum, et.operands[:len(et.operands)-1]}

	if plusTarget.check(useConcat) {
		return true
	} else {
		if !useConcat {
			return false
		}

		//Try concat-ing the first 2
		newOps := make([]int, len(et.operands)-1)
		for i, n := range et.operands[2:] {
			newOps[i+1] = n
		}

		if len(strconv.Itoa(et.operands[0]))+len(strconv.Itoa(et.operands[1])) <= len(strconv.Itoa(et.target)) {
			newOps[0] = mustConcat(et.operands[:2])

			concatTarget := eqTarget{et.target, newOps}

			if concatTarget.check(useConcat) {
				return true
			}
		}

		for i, n := range et.operands[2:] {
			newOps[i+1] = n
		}

		newOps[0] = et.operands[0] + et.operands[1]

		tgt := eqTarget{et.target, newOps}
		if tgt.check(useConcat) {
			return true
		}

		newOps[0] = et.operands[0] * et.operands[1]

		tgt = eqTarget{et.target, newOps}

		return tgt.check(useConcat)
	}
}

func mustConvert(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return num
}

func mustConcat(nums []int) int {
	if len(nums) != 2 {
		panic("mustConcat called with wrong number of args")
	}

	x, y := strconv.Itoa(nums[0]), strconv.Itoa(nums[1])

	res, err := strconv.Atoi(x + y) //This + is string concatenation!
	if err != nil {
		panic(err)
	}

	return res
}

func sum(nums []int) int {
	count := 0
	for _, n := range nums {
		count += n
	}
	return count
}

func line2tgt(l string) eqTarget {
	s1 := strings.Split(l, ": ")
	if len(s1) != 2 {
		log.Fatal("Error parsing line")
	}

	tgt := mustConvert(s1[0])
	s2 := strings.Split(s1[1], " ")
	if len(s2) < 2 {
		log.Fatal("Error parsing line")
	}

	ops := make([]int, len(s2))
	for i, s := range s2 {
		ops[i] = mustConvert(s)
	}

	return eqTarget{tgt, ops}
}

func main() {
	fname := "example.txt"
	var targets []eqTarget

	if len(os.Args) > 1 {
		fname = os.Args[1]
	}

	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		l := scanner.Text()
		targets = append(targets, line2tgt(l))
	}

	ans1 := 0
	ans2 := 0

	for _, et := range targets {
		if et.check(false) {
			ans1 += et.target
			ans2 += et.target
		} else {
			if et.check(true) {
				ans2 += et.target
			}
		}
	}

	fmt.Printf("Evaluated %s, got %d and %d\n", fname, ans1, ans2)
}
