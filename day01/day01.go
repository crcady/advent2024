package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	var leftNums []int
	var rightNums []int
	var deltas []int
	var fname string

	if len(os.Args) < 2 {
		fname = "example.txt"
	} else {
		fname = os.Args[1]
	}

	fmt.Println("Running algorithm against", fname)
	exampleFile, err := os.Open(fname)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer exampleFile.Close()

	scanner := bufio.NewScanner(exampleFile)
	for scanner.Scan() {
		line := scanner.Text()
		var left, right int
		fmt.Sscanf(line, "%d   %d", &left, &right)
		leftNums = append(leftNums, left)
		rightNums = append(rightNums, right)
	}

	sort.Ints(leftNums)
	sort.Ints(rightNums)

	for i := range leftNums {
		deltas = append(deltas, rightNums[i]-leftNums[i])
	}

	sum := 0
	for _, n := range deltas {
		if n < 0 {
			n = -n
		}
		sum += n
	}

	fmt.Println("The first is", sum)

	var similarities []int

	for _, n := range leftNums {
		similarities = append(similarities, n*countOccurences(rightNums, n))
	}

	similaritySum := 0
	for _, n := range similarities {
		similaritySum += n
	}

	fmt.Println("The second answer is", similaritySum)

}

func countOccurences(slice []int, target int) int {
	count := 0
	for _, n := range slice {
		if n == target {
			count++
		}
	}

	return count
}
