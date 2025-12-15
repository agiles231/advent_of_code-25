package main

import (
	"aoc_25_day7/tachyon"
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	grid := readGrid(scanner)
	tm := initTachyonManifold(grid)
	ts := tm.NewSimulation()
	ts.Tick()
	ts.Tick()
	ts.CompleteSimulation()
	ts.Print()
	fmt.Printf("Total splits: %v\n", ts.CountSplits())
	fmt.Printf("Total realities: %v\n", ts.GetBeamCount())

}

func readGrid(scanner *bufio.Scanner) tachyon.Grid {
	grid := make(tachyon.Grid, 0)
	for scanner.Scan() {
		line := scanner.Text()
		grid = append(grid, []rune(line))
	}
	return grid
}

func initTachyonManifold(grid tachyon.Grid) *tachyon.TachyonManifold {
	tm, err := tachyon.NewTachyonManifold(grid)
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	return tm
}
