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

func main() {
	main_part2()
}

func main_part2() {
	scanner := bufio.NewScanner(os.Stdin)
	var allRanges string
	if scanner.Scan() {
		allRanges = scanner.Text()
	}
	ranges := strings.Split(allRanges, ",")
	result := int64(0)
	for i := range ranges {
		r := ranges[i]
		rSplit := strings.Split(r, "-")
		rLower := rSplit[0]
		rUpper := rSplit[1]
		log.Printf("Range # %3v [%12v - %12v]", i, rLower, rUpper)
		rangeResult, err := analyzeRange_part2(rLower, rUpper)
		if err != nil {
			log.Fatalf("Error analyzing range %v: %v", r, err)
		}
		prevResult := result
		result += rangeResult
		if result < prevResult {
			log.Fatalf("Went down")
		}
	}
	log.Printf("Result: %v\n", result)
}

func main_part1() {
	scanner := bufio.NewScanner(os.Stdin)
	var allRanges string
	if scanner.Scan() {
		allRanges = scanner.Text()
	}
	ranges := strings.Split(allRanges, ",")
	result := int64(0)
	for i := range ranges {
		r := ranges[i]
		log.Printf("Range # %v.2 %v", i, r)
		rSplit := strings.Split(r, "-")
		rLower := rSplit[0]
		rUpper := rSplit[1]
		rangeResult, err := analyzeRange(rLower, rUpper)
		if err != nil {
			log.Fatalf("Error analyzing range %v: %v", r, err)
		}
		prevResult := result
		result += rangeResult
		if result < prevResult {
			log.Fatalf("Went down")
		}
	}
	log.Printf("Result: %v\n", result)
}

func analyzeNumber(num int) (int64, error) {
	numStr := strconv.Itoa(num)
	numLen := len(numStr)
	for patternSize := 1; patternSize <= numLen/2; patternSize++ {
		// guard uneven division
		if numLen%patternSize != 0 {
			continue
		}
		numRepeat := numLen / patternSize
		pattern := numStr[:patternSize]
		matches := true
		for j := 1; j < numRepeat; j++ {
			if numStr[patternSize*j:(j+1)*patternSize] != pattern {
				matches = false
				break
			}
		}
		if matches {
			log.Printf("%12v matches with pattern size of %v", num, patternSize)
			return int64(num), nil
		}
	}
	return 0, nil
}

func analyzeRange_part2(lower, upper string) (int64, error) {
	/*
		Because we don't can't analytically solve this case, we will simply iterate every number in every range
	*/
	lowerNum, err := strconv.Atoi(lower)
	if err != nil {
		return 0, err
	}
	upperNum, err := strconv.Atoi(upper)
	if err != nil {
		return 0, err
	}
	results := int64(0)
	for i := lowerNum; i <= upperNum; i++ {
		numResult, err := analyzeNumber(i)
		if err != nil {
			return 0, err
		}
		results += numResult
	}
	return results, nil
}

func analyzeRange(lower, upper string) (int64, error) {
	/*
		For 2 numbers composed of digits "A1", "B1", ...; and "A2", "B2", ... of length m, there is a formula to compute the number of repeating patterns between them:
		Assume each number is broken into halves, H1_1 and H1_2, and H2_1 and H2_2
		The formula is this:
		(H2_1 - H1_1 - 1) + (1 if H1_1 > H1_2 else 0) + (1 if H2_2 > H2_1 else 0)

		For ranges crossing different lengths - compute for all of every in between length.
		For a given length, the total number for the entire range of that length is 10^(m / 2) where m is the length of the number. Well, not quite:
		m = 4
		1000 - 9999 -> 1010, 1111, 1212, 1313, ..., 2020, 2121, ..., 9898, 9999
		Notice, we are missing 0-9 from 10^2. So it's actually 10^2 - 10^1. What about l = 6?
		m = 6
		100000 - 999999 -> 100100, 101101, ..., 199199, 200200, ..., 901901, ..., 999999.
		So it's 10^3 - 10^2 == 9 * 10^2. So the final formula is: 10^((m / 2) - 1) * 9
		Note m is always even, because odd m's don't have twice repeated patterns

		So the final algorithm is:
		1. Sanitize the input - discard the lower range if it's odd, discard the upper range if it's odd
		2. If both lower and upper and are the same length, compute using the in range formula
		3. If both lower and upper are different lengths:
		3.1 Compute the number between the lower range and the length above it
		3.2 Compute the number between the upper range and the length below it
		3.3 Compute the number in all ranges between the lower and upper lengths
	*/

	// Sanitize inputs
	lower = sanitizeLower(lower)
	upper = sanitizeUpper(upper)
	log.Printf("lower: %v, upper: %v\n", lower, upper)

	// Guard invalid range (may have originally been valid, but invalid after sanitization)
	lowerNum, err := strconv.Atoi(lower)
	if err != nil {
		return 0, err
	}
	upperNum, err := strconv.Atoi(upper)
	if err != nil {
		return 0, err
	}
	if upperNum <= lowerNum {
		return 0, nil
	}

	// Determine case:
	lowerM := len(lower)
	upperM := len(upper)
	if lowerM == upperM {
		result, err := computeBetweenOfSameLength(lower, upper)
		if err != nil {
			return 0, err
		}
		return int64(result), nil
	} else {
		lowersUpper := strings.Repeat("9", lowerM)
		lowerResult, err := computeBetweenOfSameLength(lower, lowersUpper)
		if err != nil {
			return 0, err
		}
		uppersLower := "1" + strings.Repeat("0", upperM-1)
		upperResult, err := computeBetweenOfSameLength(uppersLower, upper)
		if err != nil {
			return 0, err
		}
		nextLowestRange := lowerM + 2
		nextHighestRange := upperM - 2
		betweenResult := int64(0)
		if nextLowestRange <= nextHighestRange {
			betweenResult, err = computeEntiretyOfLengths(nextLowestRange, nextHighestRange)
			if err != nil {
				return 0, err
			}
		}

		fmt.Printf("Lower: %v\n", lower)
		fmt.Printf("Upper: %v\n", upper)
		fmt.Printf("lowersUpper: %v\n", lowersUpper)
		fmt.Printf("uppersLower: %v\n", uppersLower)
		fmt.Printf("nextLowestRange: %v\n", nextLowestRange)
		fmt.Printf("nextHighestRange: %v\n", nextHighestRange)
		fmt.Printf("lowerResult: %v\n", lowerResult)
		fmt.Printf("upperResult: %v\n", upperResult)
		fmt.Printf("betweenResult: %v\n", betweenResult)
		return int64(lowerResult + upperResult + betweenResult), nil
	}
}

func sanitizeLower(lower string) string {
	if len(lower)%2 == 1 {
		// odd numbered length - can't possibly have a twice repeated pattern
		// so go to the beginning of the range at the next highest length
		lowerLength := len(lower)
		newLower := "1" + strings.Repeat("0", lowerLength)
		lower = newLower
	}
	return lower
}

func sanitizeUpper(upper string) string {
	if len(upper)%2 == 1 {
		// odd numbered length - can't possibly have a twice repeated pattern
		// so go to the end of the range at the next lowest length
		upperLength := len(upper)
		newUpper := strings.Repeat("9", upperLength-1)
		upper = newUpper
	}
	return upper
}

func computeBetweenOfSameLength(lower, upper string) (int64, error) {
	lowerM := len(lower)
	upperM := len(upper)
	if lowerM != upperM {
		return 0, errors.New("lower and upper must be the same length")
	}
	lowerH1 := lower[:lowerM/2]
	lowerH2 := lower[lowerM/2:]
	upperH1 := upper[:upperM/2]
	upperH2 := upper[upperM/2:]

	lowerH1Num, err := strconv.Atoi(lowerH1)
	if err != nil {
		return 0, err
	}
	lowerH2Num, err := strconv.Atoi(lowerH2)
	if err != nil {
		return 0, err
	}
	upperH1Num, err := strconv.Atoi(upperH1)
	if err != nil {
		return 0, err
	}
	upperH2Num, err := strconv.Atoi(upperH2)
	if err != nil {
		return 0, err
	}
	lowestCrosses := lowerH2Num <= lowerH1Num
	upperCrosses := upperH1Num <= upperH2Num

	// Debug output
	fmt.Printf("L H1: %v\n", lowerH1)
	fmt.Printf("L H2: %v\n", lowerH2)
	fmt.Printf("U H1: %v\n", upperH1)
	fmt.Printf("U H2: %v\n", upperH2)
	fmt.Printf("lowestCrosses: %v\n", lowestCrosses)
	fmt.Printf("upperCrosses: %v\n", upperCrosses)
	result := int64(0)
	for i := lowerH1Num + 1; i < upperH1Num; i++ {
		numResult, err := getRepeatedNum(i)
		if err != nil {
			return 0, err
		}
		result += numResult
	}
	fmt.Printf("result: %v\n", result)

	if lowerH1 == upperH1 && lowestCrosses && upperCrosses {
		// They both cross, but are the same. Result = 1 in this case
		lowestNum, err := getRepeatedNum(lowerH1Num)
		if err != nil {
			return 0, err
		}
		result += lowestNum
	} else if lowestCrosses && lowerH1 != upperH1 {
		// The lowest crosses and upper is a different prefix. Result ++
		lowestNum, err := getRepeatedNum(lowerH1Num)
		if err != nil {
			return 0, err
		}
		result += lowestNum
	}
	if upperCrosses && lowerH1 != upperH1 {
		// The "upper cross" case is already covered above if lowerH1 == upperH1
		highestNum, err := getRepeatedNum(upperH1Num)
		if err != nil {
			return 0, err
		}
		result += highestNum
	}
	fmt.Printf("final result: %v\n", result)

	return result, nil
}

func getRepeatedNum(half int) (int64, error) {
	halfStr := strconv.Itoa(half)
	wholeStr := halfStr + halfStr
	whole, err := strconv.Atoi(wholeStr)
	if err != nil {
		return 0, err
	}
	return int64(whole), nil
}

func computeEntiretyOfLengths(lowerM, upperM int) (int64, error) {
	if lowerM%2 == 1 {
		return 0, errors.New("Lower length is odd")
	}
	if upperM%2 == 1 {
		return 0, errors.New("Upper length is odd")
	}

	result := int64(0)
	fmt.Printf("lowerM: %v\n", lowerM)
	fmt.Printf("upperM: %v\n", upperM)
	for m := lowerM; m <= upperM; m += 2 {
		lower := getLowerOfM(m)
		upper := getUpperOfM(m)
		mResult, err := computeBetweenOfSameLength(lower, upper)
		if err != nil {
			return 0, err
		}
		result += mResult
	}

	return result, nil
}

func getLowerOfM(m int) string {
	return "1" + strings.Repeat("0", m-1)
}
func getUpperOfM(m int) string {
	return strings.Repeat("9", m)
}
