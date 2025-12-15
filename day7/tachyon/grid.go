package tachyon

import (
	"fmt"
)

type Grid [][]rune

func (g Grid) PrintGrid() {
	g.IterateGrid(func(i, j int, value rune) {
		if i == 0 {
			fmt.Println()
		}
		fmt.Print(string(value))
	})
	fmt.Println()
}

func (g Grid) Height() int {
	return len(g)
}

func (g Grid) Width() int {
	return len(g[0])
}

func (g Grid) IterateGrid(callback func(i, j int, value rune)) {
	gridHeight := g.Height()
	gridWidth := g.Width()
	for i := 0; i < gridHeight; i++ {
		for j := 0; j < gridWidth; j++ {
			callback(i, j, g[i][j])
		}
	}
}

func (g Grid) Duplicate() Grid {
	grid := make(Grid, g.Height())
	for i := 0; i < g.Height(); i++ {
		for j := 0; j < g.Width(); j++ {
			grid[i][j] = g[i][j]
		}
	}
	return grid
}
