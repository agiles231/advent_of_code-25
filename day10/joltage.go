package main

import (
	"bufio"
	"cmp"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main2() {
	scanner := bufio.NewScanner(os.Stdin)
	i := 0
	for scanner.Scan() {
		i++
		line := scanner.Text()
		buttons, joltage, err := parseButtonsAndJoltage(line)
		if err != nil {
			log.Fatalf("Failure: %v\n", err)
		}
		fmt.Printf("Target: %v, Buttons: %v\n", joltage, buttons)
		minPressed, err := greedyTraversal(joltage, buttons)
		if err != nil {
			log.Fatalf("Failure: %v\n", err)
		}
		fmt.Printf("Min pressed for %v: %v\n", i, minPressed)
	}
}

func parseButtonsAndJoltage(line string) ([]Button, Joltage, error) {
	components := strings.Split(line, " ")
	// discard the first one (not related to joltage)
	components = components[1:]
	joltageStr := components[len(components)-1]
	joltageValues, err := parseWrappedStrToIntSlice(joltageStr)
	if err != nil {
		return nil, Joltage{}, err
	}
	joltage := Joltage{values: joltageValues}
	length := joltage.Length()
	buttonStrs := components[:len(components)-1]
	buttons := make([]Button, len(buttonStrs))
	for i, buttonStr := range buttonStrs {
		activeIndexes, err := parseWrappedStrToIntSlice(buttonStr)
		if err != nil {
			return nil, Joltage{}, err
		}
		buttons[i] = MakeButton(activeIndexes, length)
	}
	return buttons, joltage, nil
}

func MakeButton(activeIndexes []int, length int) Button {
	activeSwitches := make([]bool, length)
	for _, i := range activeIndexes {
		activeSwitches[i] = true
	}
	return Button{activeSwitches: activeSwitches}
}

type Stack[T any] struct {
	data []T
}

func NewStack[T any](data []T) *Stack[T] {
	return &Stack[T]{
		data: data,
	}
}

func (s *Stack[T]) Push(d T) {
	s.data = append(s.data, d)
}

func (s *Stack[T]) Pop() T {
	length := len(s.data)
	datum := s.data[length-1]
	s.data = s.data[:length-1]
	return datum
}

func (s *Stack[T]) Peek() T {
	return s.data[len(s.data)-1]
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.data) == 0
}

type Queue[T any] struct {
	data []T
}

func NewQueue[T any](data []T) *Queue[T] {
	return &Queue[T]{
		data: data,
	}
}

func (q *Queue[T]) Push(d T) {
	q.data = append(q.data, d)
}

func (q *Queue[T]) Pop() T {
	datum := q.data[0]
	q.data = q.data[1:]
	return datum
}

func (q *Queue[T]) Peek() T {
	return q.data[0]
}

func (q *Queue[T]) Length() int {
	return len(q.data)
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.data) == 0
}

func greedyTraversal(target Joltage, buttons []Button) (int, error) {
	joltageLength := target.Length()
	startJoltage := Joltage{values: make([]int, joltageLength)}
	joltageStack := NewStack[Joltage]([]Joltage{startJoltage})
	startButtons := SortButtonsByDistance(target, startJoltage, buttons)
	buttonStack := NewStack[*Queue[Button]]([]*Queue[Button]{NewQueue[Button](startButtons)})
	tree := make(map[string]*Node[Joltage])
	tree[startJoltage.String()] = &Node[Joltage]{parent: nil, data: startJoltage}
	count := 0
	for !joltageStack.IsEmpty() {
		// Get joltage
		count++
		joltage := joltageStack.Peek()
		if count%10000 == 0 {
			fmt.Println(joltage)
			fmt.Println(count)
		}
		parentNode, ok := tree[joltage.String()]
		if !ok {
			return 0, fmt.Errorf("Error: %v not found in tree", joltage)
		}
		depth := parentNode.Depth()
		depth++
		// Get button
		buttonQueue := buttonStack.Peek()
		// fmt.Printf("Queue len: %v\n", buttonQueue.Length())
		button := buttonQueue.Pop()
		if buttonQueue.IsEmpty() {
			// This joltage branch is exhausted
			joltageStack.Pop()
			buttonStack.Pop()
			// fmt.Printf("Popping %v\n", joltage)
		}
		// Make new joltage using button on current joltage
		newJoltage, err := button.TransformJoltage(joltage)
		if err != nil {
			return 0, err
		}
		// Return if target found
		if newJoltage.Equals(target) {
			return depth, nil
		}
		// Add to stacks and tree if viable and not already in tree
		_, ok = tree[newJoltage.String()]
		if !ok && !newJoltage.Exceeds(target) {
			// Add to queue with full button set
			tree[newJoltage.String()] = &Node[Joltage]{parent: parentNode, data: newJoltage}
			joltageStack.Push(newJoltage)
			sortedButtons := SortButtonsByDistance(target, startJoltage, buttons)
			buttonStack.Push(NewQueue[Button](sortedButtons))
			// fmt.Printf("Added: %v\n", newJoltage)
		} else if !ok {
			// fmt.Printf("%v exceeded target %v\n", newJoltage, target)
		} else {
			// fmt.Printf("%v already found\n", newJoltage)
		}
	}
	return 0, errors.New("Target was never found")
}

func parseCsvToIntSlice(intsStr string) ([]int, error) {
	// remove outer brackets
	intStrings := strings.Split(intsStr, ",")
	length := len(intStrings)
	ints := make([]int, length)
	for i, intString := range intStrings {
		integer, err := strconv.Atoi(intString)
		if err != nil {
			return nil, err
		}
		ints[i] = integer
	}
	return ints, nil
}

type Button struct {
	activeSwitches []bool
}

func ParseButtonStrToIntSlice(buttonStr string) ([]int, error) {
	// remove outer parens
	intsStr := buttonStr[1 : len(buttonStr)-1]
	return parseCsvToIntSlice(intsStr)
}

func NewButton(activeSwitches []bool) Button {
	return Button{activeSwitches: activeSwitches}
}

func (b Button) Equals(b2 Button) bool {
	if b.Length() != b2.Length() {
		return false
	}

	for i := range b.activeSwitches {
		v1 := b.activeSwitches[i]
		v2 := b2.activeSwitches[i]
		if v1 != v2 {
			return false
		}
	}
	return true
}

func (b Button) TransformJoltage(j Joltage) (Joltage, error) {
	if b.Length() != j.Length() {
		return Joltage{}, fmt.Errorf("Button cannot interact with joltage of different length. Button length: %v, joltage length: %v", b.Length(), j.Length())
	}
	joltageValues := make([]int, j.Length())
	for i, active := range b.activeSwitches {
		joltageValues[i] = j.values[i]
		if active {
			joltageValues[i]++
		}
	}
	return Joltage{values: joltageValues}, nil
}

func ButtonGreed(b Button) int {
	greed := 0
	for _, active := range b.activeSwitches {
		if active {
			greed++
		}
	}
	return greed
}

func (b Button) Length() int {
	return len(b.activeSwitches)
}

func (b Button) String() string {
	return fmt.Sprintf("Button [activeSwitches: %v, length: %v]", b.activeSwitches, b.Length())
}

type Joltage struct {
	values []int
}

// Parse a string representation of []int back into []int.
// Returns error if parsing fails
//
// Steps:
// 1. Remove outer chars
// 2. Split on comma
// 3. For each value
// 3.i   Parse into int value
// 3.ii  Add to slice
// 4. Return slice
func parseWrappedStrToIntSlice(intsStr string) ([]int, error) {
	// remove outer characters ([] or () or {}, etc.)
	intsStr = intsStr[1 : len(intsStr)-1]
	return parseCsvToIntSlice(intsStr)
}

func NewJoltage(values []int) Joltage {
	return Joltage{values: values}
}

func (j Joltage) Equals(j2 Joltage) bool {
	if j.Length() != j2.Length() {
		return false
	}

	for i := range j.values {
		v1 := j.values[i]
		v2 := j2.values[i]
		if v1 != v2 {
			return false
		}
	}
	return true
}

func (j Joltage) Length() int {
	return len(j.values)
}

func (j Joltage) String() string {
	return fmt.Sprintf("Joltage [values: %v]", j.values)
}

func (j Joltage) Exceeds(j2 Joltage) bool {
	if j.Length() != j2.Length() {
		return true
	}
	for i := range j.values {
		if j.values[i] > j2.values[i] {
			return true
		}
	}
	return false
}

func (j Joltage) MaxUnitDistance(j2 Joltage) int {
	maxDist := 0
	for i := range j.values {
		v1 := j.values[i]
		v2 := j2.values[i]
		diff := int(math.Abs(float64(v2) - float64(v1)))
		if diff > maxDist {
			maxDist = diff
		}
	}
	return maxDist
}

func (j Joltage) MaxSortedUnitDistances(j2 Joltage) []int {
	unitDistances := make([]int, len(j.values))
	for i := range j.values {
		v1 := j.values[i]
		v2 := j2.values[i]
		unitDistances[i] = int(math.Abs(float64(v2) - float64(v1)))
	}
	slices.SortFunc(unitDistances, func(a, b int) int {
		return cmp.Compare(b, a)
	})
	return unitDistances
}

func (j Joltage) Distance(j2 Joltage) int {
	sum := 0
	for i := range j.values {
		v1 := j.values[i]
		v2 := j2.values[i]
		diff := v2 - v1
		diffSquared := diff * diff
		sum += diffSquared
	}
	sqrt := math.Sqrt(float64(sum))
	return int(sqrt)
}

func SortButtonsByDistance(target Joltage, j Joltage, buttons []Button) []Button {
	sortedButtons := slices.Clone(buttons)
	slices.SortFunc(sortedButtons, func(a, b Button) int {
		aJ, _ := a.TransformJoltage(j)
		bJ, _ := b.TransformJoltage(j)
		aExceeds := aJ.Exceeds(target)
		bExceeds := bJ.Exceeds(target)
		if aExceeds && bExceeds {
			return 0
		} else if aExceeds {
			return 1
		} else if bExceeds {
			return -1
		}
		aUnitDists := aJ.MaxSortedUnitDistances(target)
		bUnitDists := bJ.MaxSortedUnitDistances(target)
		// fmt.Printf("A: %v\n", aUnitDists)
		// fmt.Printf("b: %v\n", bUnitDists)
		for i := range aUnitDists {
			v1 := aUnitDists[i]
			v2 := bUnitDists[i]
			if v1 != v2 {
				return v1 - v2
			}
		}
		return 0
	})
	return sortedButtons
}
