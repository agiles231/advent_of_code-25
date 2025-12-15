package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Range struct {
	lower int
	upper int
}

func NewRange(lower, upper int) Range {
	return Range{lower: lower, upper: upper}
}

func (r Range) InRange(num int) bool {
	return r.lower <= num && r.upper >= num
}

func (r Range) RangeLen() int {
	return (r.upper - r.lower) + 1
}

func main() {
	fmt.Println("Hello, world!")
	scanner := bufio.NewScanner(os.Stdin)
	/*
		Read the ranges
		Consolidate ranges
		Read the ingredients
		Iterate ingredients and determine if in any ranges
	*/
	ranges, err := ReadRanges(scanner)
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	ranges = ConsolidateRanges(ranges)
	PrintRanges(ranges)
	items, err := ReadItems(scanner)
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	PrintItems(items)
	itemLen := len(items)
	rangeLen := len(ranges)
	inRangeCount := 0
	for i := 0; i < itemLen; i++ {
		item := items[i]
		for j := 0; j < rangeLen; j++ {
			r := ranges[j]
			if r.InRange(item) {
				inRangeCount++
				break
			}
		}
	}
	fmt.Printf("In range count: %v", inRangeCount)

	totalSpan := 0
	for i := 0; i < rangeLen; i++ {
		r := ranges[i]
		totalSpan += r.RangeLen()
	}
	fmt.Printf("Total range span: %v", totalSpan)
}

func ReadItems(scanner *bufio.Scanner) ([]int, error) {
	items := make([]int, 0)
	for scanner.Scan() {
		line := scanner.Text()
		item, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func ReadRanges(scanner *bufio.Scanner) ([]Range, error) {
	ranges := make([]Range, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		rangeParts := strings.Split(line, "-")
		lower, err := strconv.Atoi(rangeParts[0])
		if err != nil {
			return nil, err
		}
		upper, err := strconv.Atoi(rangeParts[1])
		if err != nil {
			return nil, err
		}
		ranges = append(ranges, NewRange(lower, upper))
	}
	return ranges, nil
}

func PrintRanges(ranges []Range) {
	rangeLen := len(ranges)
	for i := 0; i < rangeLen; i++ {
		r := ranges[i]
		fmt.Printf("Range %3v: [%15v - %15v]\n", i, r.lower, r.upper)
	}
}

func PrintItems(items []int) {
	itemLen := len(items)
	for i := 0; i < itemLen; i++ {
		item := items[i]
		fmt.Printf("Ingredient %5v: %15v\n", i, item)
	}
}

func ConsolidateRanges(ranges []Range) []Range {
	lenRanges := len(ranges)
	SortRanges(ranges)
	consolidatedRanges := make([]Range, 0)
	for i := 0; i < lenRanges; i++ {
		r := ranges[i]
		lower := r.lower
		upper := r.upper
		jump := 0
		for j := i + 1; j < lenRanges; j++ {
			upperR := ranges[j]
			upperLower := upperR.lower
			upperUpper := upperR.upper
			if upperLower <= upper {
				jump++
				if upperUpper > upper {
					upper = upperUpper
				}
			} else {
				break
			}
		}
		i += jump
		consolidatedRanges = append(consolidatedRanges, NewRange(lower, upper))
	}
	return consolidatedRanges
}

func SortRanges(ranges []Range) {
	slices.SortFunc(ranges, func(a, b Range) int {
		return a.lower - b.lower
	})
}
