package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	grid := ReadGrid(scanner)
	gridCopy := CopyGrid(grid)
	PrintGrid(grid)
	yLen, xLen := GetGridDimensions(grid)
	symbol := '@'
	count := 0
	maxNumIterations := 10000000
	gridsNotEqual := true
	iterations := 0
	for gridsNotEqual && iterations < maxNumIterations {
		for j := 0; j < yLen; j++ {
			for i := 0; i < xLen; i++ {
				if grid[j][i] == symbol {
					neighbors := CountNeighbors(grid, j, i, symbol)
					canForkLift := CanForkliftPaperRole(neighbors)
					if canForkLift {
						count++
						gridCopy[j][i] = 'x'
					}
				}
			}
		}
		PrintGrid(gridCopy)
		gridsNotEqual = !GridsEqual(grid, gridCopy)
		if gridsNotEqual {
			grid = gridCopy
			gridCopy = CopyGrid(grid)
		}
		iterations++
	}
	fmt.Printf("Forkliftable paper rolls: %v", count)
}

func GridsEqual(grid1 [][]rune, grid2 [][]rune) bool {
	yLen1, xLen1 := GetGridDimensions(grid1)
	yLen2, xLen2 := GetGridDimensions(grid2)
	if xLen1 != xLen2 || yLen1 != yLen2 {
		return false
	}

	for j := 0; j < yLen1; j++ {
		for i := 0; i < xLen1; i++ {
			if grid1[j][i] != grid2[j][i] {
				return false
			}
		}
	}
	return true
}

func CopyGrid(grid [][]rune) [][]rune {
	yLen, xLen := GetGridDimensions(grid)
	gridCopy := make([][]rune, yLen)
	for j := 0; j < yLen; j++ {
		gridCopy[j] = make([]rune, xLen)
		for i := 0; i < xLen; i++ {
			gridCopy[j][i] = grid[j][i]
		}
	}
	return gridCopy
}

func GetGridDimensions(grid [][]rune) (int, int) {
	yLen := len(grid)
	xLen := len(grid[0])
	return yLen, xLen
}

func ReadGrid(scanner *bufio.Scanner) [][]rune {
	grid := make([][]rune, 0)
	for scanner.Scan() {
		row := scanner.Text()
		runeRow := []rune(row)
		grid = append(grid, runeRow)
	}
	return grid
}

func PrintGrid(grid [][]rune) {
	yLen, _ := GetGridDimensions(grid)
	for j := 0; j < yLen; j++ {
		row := string(grid[j])
		fmt.Println(row)
	}
}

func CountNeighbors(grid [][]rune, j, i int, symbol rune) int {
	count := 0

	yLen, xLen := GetGridDimensions(grid)
	if j != 0 {
		// We can check the row above
		if i != 0 && grid[j-1][i-1] == symbol {
			count++
		}
		if grid[j-1][i] == symbol {
			count++
		}
		if i < xLen-1 && grid[j-1][i+1] == symbol {
			count++
		}
	}
	if i != 0 && grid[j][i-1] == symbol {
		count++
	}
	if i < xLen-1 && grid[j][i+1] == symbol {
		count++
	}
	if j < yLen-1 {
		// We can check the row above
		if i != 0 && grid[j+1][i-1] == symbol {
			count++
		}
		if grid[j+1][i] == symbol {
			count++
		}
		if i < xLen-1 && grid[j+1][i+1] == symbol {
			count++
		}
	}
	return count
}

func CanForkliftPaperRole(count int) bool {
	return count < 4
}
