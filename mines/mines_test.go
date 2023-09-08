package mines

import (
	"fmt"
	"strconv"
	"testing"
)

func TestNewMineField_1Mine(t *testing.T) {
	rows, cols, mines := 10, 10, 1
	p, err := NewMineField(rows, cols, mines)
	if err != nil {
		t.Error(err)
	}
	if p.Rows != rows {
		t.Error("Rows not set.")
	}
	if p.Cols != cols {
		t.Error("Cols not set.")
	}
	if p.Mines != mines {
		t.Error("Mines not set.")
	}
	actualMines := 0
	for _, s := range p.States {
		if (s & MaskMine) == MaskMine {
			actualMines++
		}
	}
	if actualMines != mines {
		t.Errorf("Wrong number of mines set (%d)", actualMines)
	}
}

func TestNewMineField_MostMines(t *testing.T) {
	rows, cols := 99, 99
	mines := rows * cols - 1
	p, err := NewMineField(rows, cols, mines)
	if err != nil {
		t.Error(err)
	}
	if p.Rows != rows {
		t.Error("Rows not set.")
	}
	if p.Cols != cols {
		t.Error("Cols not set.")
	}
	if p.Mines != mines {
		t.Error("Mines not set.")
	}
	actualMines := 0
	for _, s := range p.States {
		if (s & MaskMine) == MaskMine {
			actualMines++
		}
	}
	if actualMines != mines {
		t.Errorf("Wrong number of mines set (%d)", actualMines)
	}
}

func TestMineField_Show_1Mine(t *testing.T) {
	rows, cols := 10, 10
	mines := 1
	p,err := NewMineField(rows, cols, mines)
	if err != nil {
		panic(err)
	}
	for i:=0; i < rows; i++ {
		for j:=0; j< cols; j++ {
			state := p.Get(i, j)
			c := "o"
			if state & MaskMine == MaskMine {
				c = "*"
			} else if state & MaskNeighboringMines != 0 {
				c = strconv.Itoa(int(state & MaskNeighboringMines))
			}
			fmt.Print(c)
		}
		fmt.Println()
	}
}

func TestMineField_Show_MostMines(t *testing.T) {
	rows, cols := 10, 10
	mines := 100
	p,err := NewMineField(rows, cols, mines)
	if err != nil {
		panic(err)
	}
	for i:=0; i < rows; i++ {
		for j:=0; j< cols; j++ {
			state := p.Get(i, j)
			c := "o"
			if state & MaskMine == MaskMine {
				c = "*"
			} else if state & MaskNeighboringMines != 0 {
				c = strconv.Itoa(int(state & MaskNeighboringMines))
			}
			fmt.Print(c)
		}
		fmt.Println()
	}
}

func TestMineField_Show_RandomMines(t *testing.T) {
	rows, cols := 10, 10
	mines := 20
	p,err := NewMineField(rows, cols, mines)
	if err != nil {
		panic(err)
	}
	for i:=0; i < rows; i++ {
		for j:=0; j< cols; j++ {
			state := p.Get(i, j)
			c := "o"
			if state & MaskMine == MaskMine {
				c = "*"
			} else if state & MaskNeighboringMines != 0 {
				c = strconv.Itoa(int(state & MaskNeighboringMines))
			}
			fmt.Print(c)
		}
		fmt.Println()
	}
}

func TestMineField_Reveal_RandomMines(t *testing.T) {
	rows, cols := 10, 10
	mines := 20
	p,err := NewMineField(rows, cols, mines)
	if err != nil {
		panic(err)
	}
	p.Reveal(0, 0)
	for i:=0; i < rows; i++ {
		for j:=0; j< cols; j++ {
			state := p.Get(i, j)
			c := "x"
			if state & MaskRevealed == MaskRevealed {
				nm := int(state & MaskNeighboringMines)
				if nm > 0 {
					c = strconv.Itoa(nm)
				} else {
					c = "o"
				}
			}
			fmt.Print(c)
		}
		fmt.Println()
	}
}

func TestMineField_Reveal_1Mine(t *testing.T) {
	rows, cols := 10, 10
	mines := 1
	p,err := NewMineField(rows, cols, mines)
	if err != nil {
		panic(err)
	}
	p.Reveal(0, 0)
	for i:=0; i < rows; i++ {
		for j:=0; j< cols; j++ {
			state := p.Get(i, j)
			c := "x"
			if state & MaskRevealed == MaskRevealed {
				nm := int(state & MaskNeighboringMines)
				if nm > 0 {
					c = strconv.Itoa(nm)
				} else {
					c = "o"
				}
			}
			fmt.Print(c)
		}
		fmt.Println()
	}
}