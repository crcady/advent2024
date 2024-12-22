package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

type secret uint64

func (s secret) mix(val secret) secret {
	return secret(s ^ val)
}

func (s secret) prune() secret {
	return secret(s % 16777216)
}

func (s secret) evolve() secret {
	s1 := s.mix(s << 6).prune()
	s2 := s1.mix(s1 >> 5).prune()
	s3 := s2.mix(s2 << 11).prune()

	return s3
}

func (s secret) price() uint8 {
	return uint8(uint64(s) % 10)
}

type priceChecker struct {
	prices []uint8
	cache  map[[4]int8]uint8
}

func newPriceChecker() *priceChecker {
	return &priceChecker{make([]uint8, 0, 2001), nil}
}

func (pc *priceChecker) addPrice(price uint8) {
	pc.prices = append(pc.prices, price)
	pc.cache = nil //Just in case I somehow call addPrice() after check()
}

func (pc *priceChecker) check(seq [4]int8) uint8 {
	if pc.cache != nil {
		return pc.cache[seq]
	}

	pc.cache = make(map[[4]int8]uint8, 0)

	currentSeq := [4]int8{0, 0, 0, 0}
	for i := 1; i < 4; i++ {
		currentSeq[i] = int8(pc.prices[i] - pc.prices[i-1])
	}

	// currentSeq is now {0, x, x, x}

	last := pc.prices[3]

	for _, p := range pc.prices[4:] {
		// Shift the current deltas over by 1
		for i := 0; i < 3; i++ {
			currentSeq[i] = currentSeq[i+1]
		}

		newDelta := int8(p) - int8(last)
		currentSeq[3] = newDelta

		if _, ok := pc.cache[currentSeq]; !ok {
			pc.cache[currentSeq] = p
		}

		last = p
	}

	return pc.cache[seq]
}

func main() {
	fname := "example.txt"

	if len(os.Args) > 1 {
		fname = os.Args[1]
	}

	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var ans1 uint64 = 0

	checkers := []*priceChecker{}

	for scanner.Scan() {
		line := scanner.Text()
		num, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}

		s := secret(uint64(num))

		checker := newPriceChecker()
		checker.addPrice(s.price())

		for i := 0; i < 2000; i++ {
			s = s.evolve()
			checker.addPrice(s.price())
		}

		ans1 += uint64(s)
		checkers = append(checkers, checker)
	}
	log.Println("Sum of the new secrets is", ans1)

	mostBananas := 0

	for i := -9; i < 10; i++ {
		for j := -9; j < 10; j++ {
			for k := -9; k < 10; k++ {
				for l := -9; l < 10; l++ {
					seq := [4]int8{int8(i), int8(j), int8(k), int8(l)}
					bananas := 0
					for _, pc := range checkers {
						bananas += int(pc.check(seq))
					}
					if bananas > mostBananas {
						mostBananas = bananas
					}
				}
			}
		}
	}

	log.Println("Sold my secrets for", mostBananas, "bananas. Gross!")

}
