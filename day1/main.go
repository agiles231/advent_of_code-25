package main

import "fmt"
import "log"
import "bufio"
import "os"
import "strconv"

func main() {
	fmt.Println("Hello, world")
	scanner := bufio.NewScanner(os.Stdin)
	spaceMax := 100
	current := 50
	count := 0
	for scanner.Scan() {
		input := scanner.Text()
		direction := input[0]
		numStr := input[1:]
		num, err := strconv.Atoi(numStr)
		if err != nil {
			log.Fatal("Atoi failed")
		}

		// handle full rotations
		fullRotations := num / spaceMax
		count += fullRotations
		num = num % spaceMax
		if num == 0 {
			continue
		}

		// handle sign
		sign := 1
		if direction == 'L' {
			sign = -1
		}
		num = num * sign

		prevCur := current
		current = ((current + num) + spaceMax) % spaceMax
		if sign == -1 && prevCur <= current && prevCur != 0 {
			count++
		} else if sign == 1 && prevCur >= current {
			count++
		} else if current == 0 {
			count++
		}
	}
	fmt.Printf("Total count: %v", count)
}
