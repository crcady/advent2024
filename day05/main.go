package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type rules map[int][]int

func (r rules) addString(s string) error {

	splits := strings.Split(s, "|")
	if len(splits) != 2 {
		return errors.New("couldn't parse line into rule")
	}

	x, err := strconv.Atoi(splits[0])
	if err != nil {
		return err
	}

	y, err := strconv.Atoi(splits[1])
	if err != nil {
		return err
	}

	if slc, ok := r[y]; ok {
		r[y] = append(slc, x)
	} else {
		r[y] = []int{x}
	}
	return nil
}

func (r rules) check(ints []int) int {
	banned := make(map[int]bool)

	for i, n := range ints {
		//Check to see if we encountered a banned number
		if _, ok := banned[n]; ok {
			return i
		}

		//We didn't so we need to incorporate the existing rules
		if b, ok := r[n]; ok {
			for _, n2 := range b {
				banned[n2] = true
			}
		}
	}
	return -1
}

func (r rules) fix(ints []int) {
	idx := r.check(ints)
	var temp int
	for idx != -1 {
		temp = ints[idx]
		ints[idx] = ints[idx-1]
		ints[idx-1] = temp

		idx = r.check(ints)
	}
}

func split2ints(s string) []int {
	nums := strings.Split(s, ",")
	res := make([]int, len(nums))

	for i, num := range nums {
		n, _ := strconv.Atoi(num)
		res[i] = n
	}

	return res

}

func main() {
	var fname string
	var r = make(rules)

	if len(os.Args) > 1 {
		fname = os.Args[1]
	} else {
		fname = "example.txt"
	}

	file, err := os.Open(fname)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	secondHalf := false
	ans1, ans2 := 0, 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		switch {
		case l == "":
			secondHalf = true
		case secondHalf:
			ints := split2ints(l)
			if r.check(ints) == -1 {
				ans1 += ints[(len(ints)-1)/2]
			} else {
				r.fix(ints)
				ans2 += ints[(len(ints)-1)/2]
			}
		default:
			err := r.addString(l)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	fmt.Println(ans1)
	fmt.Println(ans2)

}
