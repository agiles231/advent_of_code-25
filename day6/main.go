package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type MathProblem struct {
	numbers   []int
	operation string
}

func NewMathProblem(numbers []int, operation string) (*MathProblem, error) {
	if operation != "*" && operation != "+" {
		return nil, fmt.Errorf("MathProblem operation must be '*' or '+'. Got %v", operation)
	}
	return &MathProblem{numbers: numbers, operation: operation}, nil
}

func (m MathProblem) Compute() int {
	numLen := len(m.numbers)
	total := 0
	if m.operation == "*" {
		total = 1
	}
	for i := 0; i < numLen; i++ {
		num := m.numbers[i]
		if m.operation == "+" {
			total += num
		} else if m.operation == "*" {
			total *= num
		}
	}
	return total
}

func (m MathProblem) Print() {
	numLen := len(m.numbers)
	fmt.Print("[")
	for i := 0; i < numLen-1; i++ {
		fmt.Printf("%6v, ", m.numbers[i])
	}
	fmt.Printf("%6v]", m.numbers[numLen-1])
	fmt.Printf(" :: Operation: %v\n", m.operation)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	// mathProblems, err := ReadMathProblems(scanner)
	mathProblems, err := ReadMathProblemsPart2(scanner)
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	numMathProblems := len(mathProblems)
	total := 0
	for i := 0; i < numMathProblems; i++ {
		mathProblem := mathProblems[i]
		mathProblem.Print()
		computedValue := mathProblem.Compute()
		total += computedValue
		fmt.Printf("Computed value: %v\n", computedValue)
	}
	fmt.Printf("Total: %v", total)
}

func ReadMathProblemsPart2(scanner *bufio.Scanner) ([]*MathProblem, error) {
	lines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	lineLen := len(lines[0])
	numLines := len(lines)
	numProblems := len(strings.Fields(lines[0]))
	mathProblems := make([]*MathProblem, numProblems)
	problemTerms := make([]int, 0)
	operator := ""
	problemIndex := 0
	for i := 0; i < lineLen; i++ {
		allBlank := true
		termCol := ""
		for j := 0; j < numLines; j++ {
			line := lines[j]
			if line[i] == '*' || line[i] == '+' {
				operator = string(line[i])
			} else if line[i] != ' ' {
				allBlank = false
				termCol += string(line[i])
			}
		}
		if allBlank {
			// this is a separator, create our MathProblem and start a new one
			problem, err := NewMathProblem(problemTerms, operator)
			if err != nil {
				return nil, err
			}
			mathProblems[problemIndex] = problem
			problemIndex++
			operator = ""
			problemTerms = make([]int, 0)
		} else {
			term, err := strconv.Atoi(termCol)
			if err != nil {
				return nil, err
			}
			problemTerms = append(problemTerms, term)
		}
	}
	problem, err := NewMathProblem(problemTerms, operator)
	if err != nil {
		return nil, err
	}
	mathProblems[problemIndex] = problem
	fmt.Printf("Created %3v math problems: ", problemIndex+1)
	return mathProblems, nil
}

func ReadMathProblems(scanner *bufio.Scanner) ([]*MathProblem, error) {
	numberRows := make([][]int, 0)
	operations := make([]string, 0)
	j := 0
	for scanner.Scan() {
		line := scanner.Text()
		numberStrs := strings.Fields(line)
		if numberStrs[0] == "*" || numberStrs[1] == "+" {
			operations = numberStrs
		} else {
			numberRow, err := ParseNumbers(numberStrs)
			if err != nil {
				return nil, err
			}
			numberRows = append(numberRows, numberRow)
		}
		j++
	}
	mathProblems, err := MakeMathProblems(numberRows, operations)
	if err != nil {
		return nil, err
	}
	return mathProblems, nil
}

func MakeMathProblemsPart2(numberRows [][]string, operations []string) ([]*MathProblem, error) {
	numTerms := len(numberRows)
	numProblems := len(numberRows[0])

	mathProblems := make([]*MathProblem, numProblems)
	for i := 0; i < numProblems; i++ {
		numberStrs := make([]string, numTerms)
		for j := 0; j < numTerms; j++ {
			numberStrs[j] = numberRows[j][i]
		}
		operation := operations[i]
		numbers, err := MakeNumbers(numberStrs)
		if err != nil {
			return nil, err
		}
		mathProblem, err := NewMathProblem(numbers, operation)
		if err != nil {
			return nil, err
		}
		mathProblems[i] = mathProblem
	}
	return mathProblems, nil
}

func MakeNumbers(numberStrs []string) ([]int, error) {
	colLen := len(numberStrs)
	terms := make([]int, 0)
	maxLen := 0
	for i := 0; i < colLen; i++ {
		maxLen = max(maxLen, len(numberStrs[i]))
	}
	for i := 0; i < maxLen; i++ {
		accumulatedNumber := ""
		for j := 0; j < colLen; j++ {
			numberStr := numberStrs[j]
			numLen := len(numberStr)
			index := numLen - maxLen + i
			if index < 0 {
				continue
			}
			accumulatedNumber += string(numberStr[index])
		}
		term, err := strconv.Atoi(accumulatedNumber)
		if err != nil {
			return nil, err
		}
		terms = append(terms, term)
	}
	return terms, nil
}

func MakeMathProblems(numberRows [][]int, operations []string) ([]*MathProblem, error) {
	numTerms := len(numberRows)
	numProblems := len(numberRows[0])

	mathProblems := make([]*MathProblem, numProblems)
	for i := 0; i < numProblems; i++ {
		numbers := make([]int, numTerms)
		for j := 0; j < numTerms; j++ {
			numbers[j] = numberRows[j][i]
		}
		operation := operations[i]
		mathProblem, err := NewMathProblem(numbers, operation)
		if err != nil {
			return nil, err
		}
		mathProblems[i] = mathProblem
	}
	return mathProblems, nil
}

func ParseNumbers(numberStrs []string) ([]int, error) {
	numLen := len(numberStrs)

	numbers := make([]int, numLen)
	for i := 0; i < numLen; i++ {
		numberStr := numberStrs[i]
		num, err := strconv.Atoi(numberStr)
		if err != nil {
			return nil, err
		}
		numbers[i] = num
	}
	return numbers, nil
}
