package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var reports [][]int
	var filename string

	if len(os.Args) < 2 {
		filename = "example.txt"
	} else {
		filename = os.Args[1]
	}

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var currentReport []int
		line := scanner.Text()
		words := strings.Fields(line)
		for _, w := range words {
			n, _ := strconv.Atoi(w)
			currentReport = append(currentReport, n)
		}
		reports = append(reports, currentReport)
	}

	safeCount := 0
	for _, report := range reports {
		if checkReport(report) {
			safeCount++
		}
	}

	fmt.Println("The first answer is:", safeCount)

	sortaSafeCount := 0
	for _, report := range reports {
		if checkReport(report) {
			sortaSafeCount++
			continue
		}

		if len(report) < 3 {
			continue
		}

		foundTolerant := false
		for i := range report {
			var newReport []int
			newReport = append(newReport, report...)
			newReport = append(newReport[:i], newReport[i+1:]...)
			if checkReport(newReport) {
				foundTolerant = true
			}
		}
		if foundTolerant {
			sortaSafeCount++
		}
	}
	fmt.Println("The second answer is:", sortaSafeCount)
}

func checkReport(report []int) bool {
	if report[1] == report[0] {
		return false
	}

	if report[1] > report[0] {
		for i := range report {
			if i == 0 {
				continue
			}
			delta := report[i] - report[i-1]
			if delta > 3 || delta < 1 {
				break
			}
			if i == len(report)-1 {
				return true
			}
		}
	} else {
		for i := range report {
			if i == 0 {
				continue
			}
			delta := report[i-1] - report[i]
			if delta > 3 || delta < 1 {
				break
			}
			if i == len(report)-1 {
				return true
			}
		}
	}
	return false
}
