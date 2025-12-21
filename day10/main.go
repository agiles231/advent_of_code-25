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

/*
How can I make this more efficient? I cannot explore every permutation - there is no upper bound - some patterns are 8 slots long, with 100s of variations for each slot, resulting in too many permutations to calculate all

I could potentially find multiples of existing permutations. For a new permutation, check if it's a multiple of a seen perm - if so, discard. This should greatly reduce the number of perms...but I don't konw if it's enough. However, I'm also not sure what other optimizations I can make.
I considered that some buttons are multiples of other buttons, and if they are truly an entire multiple, MOST of the time, they are interchangeable. But they aren't interchangeable when we get to the last bit.

I think the actual approach I'll have to use is greedy with back tracking:
1. Always use the button with the most total value, it it can be used. If I get to a part where no button can be used where I don't "bust" and I can't get the desired numbers, back track and try other combos. Yep, this is the way. So it really comes down to switching to a greedy DFS vs a BFS
*/

func main() {
	main2()
}

func main1() {
	scanner := bufio.NewScanner(os.Stdin)
	lineNum := 1
	totalButtonPresses := 0
	for scanner.Scan() {
		line := scanner.Text()
		desiredPattern, buttons, _, err := parseMachine(line)
		if err != nil {
			log.Fatalf("Failure: %v\n", err)
		}
		buttonPresses, err := determineLeastButtonPresses(desiredPattern, buttons)
		if err != nil {
			log.Fatalf("Failure: %v\n", err)
		}
		fmt.Printf("Min button presses for line %2v: %v\n", lineNum, buttonPresses)
		totalButtonPresses += buttonPresses
	}
	fmt.Printf("Min button presses for all lines: %v\n", totalButtonPresses)
}

func parseMachine(line string) (string, []string, string, error) {
	components := strings.Split(line, " ")
	numComponents := len(components)
	desiredPattern := components[0]
	// Remove the brackets
	desiredPattern = desiredPattern[1 : len(desiredPattern)-1]
	patternLen := len(desiredPattern)
	joltage := components[numComponents-1]
	joltage = joltage[1 : len(joltage)-1]
	rawButtons := components[1 : numComponents-1]
	numButtons := len(rawButtons)
	buttons := make([]string, numButtons)
	for i, rawButton := range rawButtons {
		button := strings.Repeat(".", patternLen)
		buttonNoParens := rawButton[1 : len(rawButton)-1]
		indexStrs := strings.Split(buttonNoParens, ",")
		for _, indexStr := range indexStrs {
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return "", nil, "", err
			}
			button = replaceAtIndex(button, 't', index)
		}
		buttons[i] = button
	}
	fmt.Printf("Desired pattern: %v, buttons: %v, joltage: %v\n", desiredPattern, buttons, joltage)
	return desiredPattern, buttons, joltage, nil
}

func determineLeastButtonPressesForJoltage(joltage string, availableButtons []string) (int, error) {
	joltageInts, err := StrToInts(joltage)
	if err != nil {
		return 0, err
	}
	patternLen := len(joltageInts)
	startPattern := "0" + strings.Repeat(",0", patternLen-1)
	patternStack := []string{startPattern}
	seenPatterns := make(map[string]*Node[string])
	seenPatterns[startPattern] = &Node[string]{parent: nil, data: startPattern}
	for len(patternStack) > 0 {
		// BFS
		pattern := patternStack[0]
		patternStack = patternStack[1:]
		// fmt.Printf("patternStack: %v, seenPatterns: %v, pattern: %v\n", patternStack, seenPatterns, pattern)
		node, ok := seenPatterns[pattern]
		if !ok {
			return 0, fmt.Errorf("Something wrong happened - %v, search pattern :%v", seenPatterns, pattern)
		}
		depth := node.Depth()
		depth++
		for _, button := range availableButtons {
			newPattern, err := evolvePatternIncrement(pattern, button)
			if err != nil {
				return 0, err
			}
			// fmt.Println(newPattern)
			exceeds, err := DoesPatternExceedTarget(joltage, pattern)
			if err != nil {
				return 0, err
			}
			if exceeds {
				continue
			}
			existingNode, ok := seenPatterns[newPattern]
			if ok {
				existingDepth := existingNode.Depth()
				if depth < existingDepth {
					seenPatterns[newPattern].parent = node
				}
			} else {
				seenPatterns[newPattern] = &Node[string]{parent: node, data: newPattern}
				patternStack = append(patternStack, newPattern)
			}
		}
	}
	joltageNode, ok := seenPatterns[joltage]
	if !ok {
		return 0, errors.New("Joltage pattern wasn't seen!!")
	}
	return joltageNode.Depth(), nil
}

func determineLeastButtonPresses(desiredPattern string, availableButtons []string) (int, error) {
	patternLen := len(desiredPattern)
	startPattern := strings.Repeat(".", patternLen)
	patternStack := []string{startPattern}
	seenPatterns := make(map[string]*Node[string])
	seenPatterns[startPattern] = &Node[string]{parent: nil, data: startPattern}
	for len(patternStack) > 0 {
		// BFS
		pattern := patternStack[0]
		patternStack = patternStack[1:]
		node, ok := seenPatterns[pattern]
		if !ok {
			return 0, fmt.Errorf("Something wrong happened - %v, search pattern :%v", seenPatterns, pattern)
		}
		depth := node.Depth()
		depth++
		for _, button := range availableButtons {
			newPattern := evolvePattern(pattern, button)
			fmt.Println(newPattern)
			existingNode, ok := seenPatterns[newPattern]
			if ok {
				existingDepth := existingNode.Depth()
				if depth < existingDepth {
					seenPatterns[newPattern].parent = node
				}
			} else {
				seenPatterns[newPattern] = &Node[string]{parent: node, data: newPattern}
				patternStack = append(patternStack, newPattern)
			}
		}
	}
	desiredPattnerNode, ok := seenPatterns[desiredPattern]
	if !ok {
		return 0, errors.New("Desired pattern wasn't seen!!")
	}
	return desiredPattnerNode.Depth(), nil
}

func evolvePatternIncrement(pattern string, button string) (string, error) {
	intParts, err := StrToInts(pattern)
	if err != nil {
		return "", err
	}
	for i := range intParts {
		if button[i] == 't' {
			intParts[i] += 1
		}
	}
	newPattern := IntsToStr(intParts)
	return newPattern, nil
}

func StrToInts(pattern string) ([]int, error) {
	parts := strings.Split(pattern, ",")
	ints := make([]int, len(parts))
	for i, part := range parts {
		integer, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		ints[i] = integer
	}
	return ints, nil
}

func IntsToStr(ints []int) string {
	pattern := ""
	for i, integer := range ints {
		if i != 0 {
			pattern += ","
		}
		pattern += strconv.Itoa(integer)
	}
	return pattern
}

func DoesPatternExceedTarget(targetPattern string, pattern string) (bool, error) {
	targetInts, err := StrToInts(targetPattern)
	if err != nil {
		return false, err
	}
	ints, err := StrToInts(pattern)
	if err != nil {
		return false, err
	}

	for i := range targetInts {
		if ints[i] > targetInts[i] {
			return true, nil
		}
	}
	return false, nil
}

func evolvePattern(pattern string, button string) string {
	newPattern := pattern
	for i := range newPattern {
		if button[i] == 't' {
			if newPattern[i] == '.' {
				newPattern = replaceAtIndex(newPattern, '#', i)
			} else {
				newPattern = replaceAtIndex(newPattern, '.', i)
			}
		}
	}
	return newPattern
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

type Node[T any] struct {
	parent *Node[T]
	data   T
}

func (n *Node[T]) Depth() int {
	cur := n
	depth := 0
	for cur.parent != nil {
		depth++
		cur = cur.parent
	}
	return depth
}
