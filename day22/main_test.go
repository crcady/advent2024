package main

import "testing"

func TestMix(t *testing.T) {
	ans := secret(42).mix(secret(15))
	if ans != 37 {
		t.Errorf("Mixed 42 and 15, expected 37, got %d", ans)
	}
}

func TestPrune(t *testing.T) {
	ans := secret(100000000).prune()
	if ans != 16113920 {
		t.Errorf("Pruned 100000000, expected 16113920, got %d", ans)
	}
}

func TestEvolve(t *testing.T) {
	ans := secret(123).evolve()
	if ans != 15887950 {
		t.Errorf("Evolved 123, expected 15887950, got %d", ans)
	}
}

func TestBuyer(t *testing.T) {
	ans := []uint8{7, 7, 0, 9}

	for i, n := range []uint64{1, 2, 3, 2024} {
		s := secret(n)
		pc := newPriceChecker()
		pc.addPrice(s.price())
		for i := 0; i < 2000; i++ {
			s = s.evolve()
			pc.addPrice(s.price())
		}

		buyPrice := pc.check([4]int8{-2, 1, -1, 3})
		if buyPrice != ans[i] {
			t.Errorf("Expected a buy price of %d and got %d", ans[i], buyPrice)
		}
	}
}
