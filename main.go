package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	c "github.com/barrettj12/collections"
	"github.com/barrettj12/screen"
)

func main() {
	grid := loadSudokuGrid("sudoku")
	ps := NewPartialSoln(&grid)

	screen.Clear()
	printGrid(grid)
	// printPartialSoln(ps)
	screen.Update()

	for deadlock := false; !deadlock; {
		deadlock = true
		time.Sleep(2 * time.Second)

		// Fill in squares with only one possibility
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if grid[i][j] != Blank {
					continue
				}

				if ps.possible[i][j].Size() == 1 {
					ps.FillSquare(i, j)
					deadlock = false
				}
			}
		}

		// If a given row/col/box has only one place that a number X can go
		// then we know X has to go there.
		// TODO
		// ps.FillUniqueSquares()

		screen.Clear()
		printGrid(*ps.grid)
		// printPartialSoln(ps)
		screen.Update()
	}

	screen.Println("deadlock! can't solve")
	printPartialSoln(ps)
	screen.Update()
}

type SudokuGrid [9][9]Square

type Square int

const Blank Square = 0

func ParseSquare(b byte) Square {
	switch b {
	case ' ':
		return Blank
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return Square(E(strconv.Atoi(string([]byte{b}))))
	default:
		panic(fmt.Sprintf("invalid Square %q", b))
	}
}

func (s Square) String() string {
	switch s {
	case Blank:
		return " "
	default:
		return fmt.Sprint(int(s))
	}
}

func loadSudokuGrid(filename string) SudokuGrid {
	data := E(os.ReadFile(filename))
	rows := strings.Split(string(data), "\n")

	grid := SudokuGrid{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			grid[i][j] = ParseSquare(rows[i][j])
		}
	}

	return grid
}

var horizLines = []string{
	"╔═══╤═══╤═══╦═══╤═══╤═══╦═══╤═══╤═══╗",
	"╟───┼───┼───╫───┼───┼───╫───┼───┼───╢",
	"╟───┼───┼───╫───┼───┼───╫───┼───┼───╢",
	"╠═══╪═══╪═══╬═══╪═══╪═══╬═══╪═══╪═══╣",
	"╟───┼───┼───╫───┼───┼───╫───┼───┼───╢",
	"╟───┼───┼───╫───┼───┼───╫───┼───┼───╢",
	"╠═══╪═══╪═══╬═══╪═══╪═══╬═══╪═══╪═══╣",
	"╟───┼───┼───╫───┼───┼───╫───┼───┼───╢",
	"╟───┼───┼───╫───┼───┼───╫───┼───┼───╢",
	"╚═══╧═══╧═══╩═══╧═══╧═══╩═══╧═══╧═══╝",
}

func printGrid(grid SudokuGrid) {
	for i := 0; i < 10; i++ {
		screen.Println(horizLines[i])
		if i == 9 {
			break
		}

		// Print row
		screen.Print("║")
		for j, sq := range grid[i] {
			div := "│"
			if j%3 == 2 {
				div = "║"
			}
			screen.Printf(" %s %s", sq, div)
		}
		screen.Println()
	}
}

// Lists possible values for each square.
type PartialSoln struct {
	grid     *SudokuGrid
	possible [9][9]*c.Set[Square]
}

func NewPartialSoln(grid *SudokuGrid) PartialSoln {
	ps := PartialSoln{
		grid: grid,
	}

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if grid[i][j] == Blank {
				ps.possible[i][j] = getPossible(grid, i, j)
			} else {
				ps.possible[i][j] = c.NewSet[Square](0)
			}
		}
	}

	return ps
}

// Only valid for blank squares.
func getPossible(grid *SudokuGrid, x, y int) *c.Set[Square] {
	possible := c.NewSet[Square](9)
	for i := 1; i <= 9; i++ {
		possible.Add(Square(i))
	}

	// Remove all in row
	for i := 0; i < 9; i++ {
		possible.Remove(grid[i][y])
	}

	// Remove all in col
	for j := 0; j < 9; j++ {
		possible.Remove(grid[x][j])
	}

	// Remove all in box
	sqx := 3 * (x / 3)
	sqy := 3 * (y / 3)
	for i := sqx; i < sqx+3; i++ {
		for j := sqy; j < sqy+3; j++ {
			possible.Remove(grid[i][j])
		}
	}

	return possible
}

func printPartialSoln(ps PartialSoln) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if ps.grid[i][j] != Blank {
				continue
			}

			screen.Printf("(%d, %d): possible values are %v\n", i, j, ps.possible[i][j])
		}
	}
}

func (ps *PartialSoln) FillSquare(x, y int) {
	var sq Square
	for b := range *ps.possible[x][y] {
		sq = b
	}
	ps.grid[x][y] = sq

	// Update possible solutions to related positions
	// Remove all in row
	for i := 0; i < 9; i++ {
		ps.possible[i][y].Remove(sq)
	}

	// Remove all in col
	for j := 0; j < 9; j++ {
		ps.possible[x][j].Remove(sq)
	}

	// Remove all in box
	sqx := 3 * (x / 3)
	sqy := 3 * (y / 3)
	for i := sqx; i < sqx+3; i++ {
		for j := sqy; j < sqy+3; j++ {
			ps.possible[i][j].Remove(sq)
		}
	}
}

func E[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
