package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	result := int64(0)
	for scanner.Scan() {
		batteryBank := scanner.Text()
		bankResult, err := handleBatteryBank(batteryBank, 12)
		if err != nil {
			log.Fatalf("Failed: %v", err)
		}
		result += int64(bankResult)
	}
	log.Printf("Result: %v", result)
}

func handleBatteryBank(bank string, n int) (int64, error) {
	const sentinel = -1
	bankLen := len(bank)
	if bankLen < n {
		return 0, fmt.Errorf("Battery bank lenght needs to be >= %v. Given %v", n, bankLen)
	}
	highest := make([]int, n)
	for i := 0; i < n; i++ {
		highest[i] = sentinel
	}
	for i := 0; i < bankLen; i++ {
		num, err := getNum(bank, i)
		if err != nil {
			return 0, err
		}
		firstSelectable := n - min(bankLen-i, n)
		maybeSelect(highest, firstSelectable, num, sentinel)
	}
	joltage, err := makeNumber(highest)
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	return joltage, nil
}

func makeNumber(highest []int) (int64, error) {
	n := len(highest)
	totalStr := ""
	for i := 0; i < n; i++ {
		totalStr += strconv.Itoa(highest[i])
	}
	total, err := strconv.ParseInt(totalStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func maybeSelect(highest []int, firstSelectable int, num int, sentinel int) (int, bool) {
	n := len(highest)
	for j := firstSelectable; j < n; j++ {
		if highest[j] < num {
			highest[j] = num
			zeroed := setSentinel(highest, j+1, sentinel)
			return zeroed, true
		}
	}
	return 0, false
}

func setSentinel(highest []int, start int, sentinel int) int {
	n := len(highest)
	for i := start; i < n; i++ {
		highest[i] = sentinel
	}
	return n - start
}

func getNum(slice string, i int) (int, error) {
	num, err := strconv.Atoi(slice[i : i+1])
	if err != nil {
		return 0, nil
	}
	return num, nil
}
