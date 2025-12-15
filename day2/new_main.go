package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main_1() {
	scanner := bufio.NewScanner(os.Stdin)
	var allRanges string
	if scanner.Scan() {
		allRanges = scanner.Text()
	}
	ranges := strings.Split(allRanges, ",")

	result := 0
	for i := range ranges {
		r := ranges[i]
		rSplit := strings.Split(r, "-")
		rangeResult, err := analyzeRange_1(rSplit)
		if err != nil {
			log.Fatalf("Failure: %v", err)
		}
		result += rangeResult
	}

	fmt.Printf("Result: %v", result)
}

func analyzeRange_1(r []string) (int, error) {
	_, _, err := getLowerAndUpper(r)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func validateRangeSlice(r []string) error {
	rangeLen := len(r)
	if rangeLen != 2 {
		return fmt.Errorf("Range must have length 2. Has length %v", rangeLen)
	}
	return nil
}

func getLower(r []string) (string, error) {
	err := validateRangeSlice(r)
	if err != nil {
		return "", err
	}
	return r[0], nil
}

func getUpper(r []string) (string, error) {
	err := validateRangeSlice(r)
	if err != nil {
		return "", err
	}
	return r[1], nil
}

func getLowerAndUpper(r []string) (string, string, error) {
	lower, err := getLower(r)
	if err != nil {
		return "", "", err
	}
	upper, err := getUpper(r)
	if err != nil {
		return "", "", err
	}
	return lower, upper, nil
}
