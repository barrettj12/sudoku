package main

import (
	"os"

	"github.com/jedib0t/go-sudoku/generator"
)

func main() {
	gen := generator.BackTrackingGenerator()
	grid, err := gen.Generate(nil)
	if err != nil {
		panic(err)
	}
	// Number of filled positions - smaller is harder
	grid.ApplyDifficulty(35)

	// Write CSV to file
	os.WriteFile("sudoku", []byte(grid.String()), os.ModePerm)
}
