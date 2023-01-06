package main

import (
	"bytes"
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
		time.Sleep(500 * time.Millisecond)

		// Fill in squares with only one possibility
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if grid[i][j] != Blank {
					continue
				}

				if ps.possible[i][j].Size() == 1 {
					var sq Square
					for b := range *ps.possible[i][j] {
						sq = b
					}

					ps.FillSquare(i, j, sq)
					deadlock = false
				}
			}
		}

		// If a given row/col/box has only one place that a number X can go
		// then we know X has to go there.
		for n := 0; n < 9; n++ {
			ps.FillUniqueSquares(row(n))
			ps.FillUniqueSquares(col(n))
			ps.FillUniqueSquares(box(n))
		}

		screen.Clear()
		printGrid(*ps.grid)
		// printPartialSoln(ps)
		screen.Update()
	}

	// Solver has done as much as it can
	// Either we've solved it, or we've reached a deadlock.
	solved := true
	for _, row := range ps.grid {
		for _, sq := range row {
			if sq == Blank {
				solved = false
				break
			}
		}
	}

	if solved {
		screen.Println("sudoku solved :)")
	} else {
		screen.Println("deadlock! can't solve")
		printPartialSoln(ps)
	}
	screen.Update()

	// Write sudoku back to file
	os.WriteFile("sudoku", ps.grid.CSV(), os.ModePerm)
}

type SudokuGrid [9][9]Square

// Write the grid as a CSV-formatted byte slice.
func (g SudokuGrid) CSV() []byte {
	buf := &bytes.Buffer{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if j > 0 {
				fmt.Fprint(buf, ",")
			}
			fmt.Fprint(buf, g[i][j])
		}
		fmt.Fprintln(buf)
	}
	return buf.Bytes()
}

type Square int

const Blank Square = 0

func ParseSquare(s string) Square {
	switch s {
	case " ":
		return Blank
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		return Square(E(strconv.Atoi(s)))
	default:
		panic(fmt.Sprintf("invalid Square %q", s))
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
		row := strings.Split(rows[i], ",")
		for j := 0; j < 9; j++ {
			grid[i][j] = ParseSquare(row[j])
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

func (ps *PartialSoln) FillSquare(x, y int, sq Square) {
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

func (ps *PartialSoln) FillUniqueSquares(group []Pos) {
	// Map each number to positions in the group where it could appear
	positions := make(map[Square][]Pos, 9)

	for _, pos := range group {
		for sq := range *ps.possible[pos.x][pos.y] {
			positions[sq] = append(positions[sq], pos)
		}
	}

	// If positions[i] has len 1, then we can fill it
	for sq, poss := range positions {
		if len(poss) == 1 {
			p := poss[0]
			ps.FillSquare(p.x, p.y, sq)
		}
	}
}

type Pos struct {
	x, y int
}

func row(y int) []Pos {
	r := make([]Pos, 9)
	for i := 0; i < 9; i++ {
		r = append(r, Pos{i, y})
	}
	return r
}

func col(x int) []Pos {
	c := make([]Pos, 9)
	for j := 0; j < 9; j++ {
		c = append(c, Pos{x, j})
	}
	return c
}

func box(n int) []Pos {
	sqx := (3 * n % 9)
	sqy := n - n%3

	b := make([]Pos, 9)
	for i := sqx; i < sqx+3; i++ {
		for j := sqy; j < sqy+3; j++ {
			b = append(b, Pos{i, j})
		}
	}
	return b
}

func E[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
