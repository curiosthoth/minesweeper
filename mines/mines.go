package mines

import (
	"fmt"
	"math/rand"
	"time"
)

type BlockState uint8

// We use `MaskInvalid` for multiple purposes
// 	- MaskInvalid + MaskFlagged + ~MaskMine = Wrong flagged mine
// 	- MaskInvalid + MaskMine = Stepped on Mine
const (
	MaskEmpty            BlockState = 0b00000000
	MaskNeighboringMines            = 0b00000111
	MaskMine                        = 0b00001000
	MaskRevealed                    = 0b00010000
	MaskFlagged                     = 0b00100000
	MaskQuestioned                  = 0b01000000
	MaskInvalid                     = 0b10000000
)

// NeighborStates holds the states of 8 neighbors
//
type NeighborStates [8]BlockState

type MineField struct {
	Rows   int
	Cols   int
	Mines  int
	States []BlockState
	Flags  int
}

// NewMineField creates a new instance of MineField.
func NewMineField(rows int, cols int, mines int) (*MineField, error) {
	if rows < 2 || rows >= 100 {
		return nil, fmt.Errorf("wrong number of rows (%d). Must be within [2, 100)", rows)
	}
	if cols < 8 || cols >= 100 {
		return nil, fmt.Errorf("wrong number of cols (%d). Must be within [8, 100)", cols)
	}
	if mines < 1 || mines >= 9999 {
		return nil, fmt.Errorf("wrong number of cols (%d). Must be within [1, 9999)", cols)
	}

	mineField := &MineField{
		Rows: rows, Cols: cols, Mines: mines, Flags: mines,
		States: make([]BlockState, rows*cols),
	}

	rand.Seed(time.Now().UnixNano())
	size := rows * cols
	indexes := make([]int, size)
	for i := 0; i < size; i++ {
		indexes[i] = i
	}
	// randomize indexes and cut off the first `mines` to lay mines
	rand.Shuffle(size, func(i, j int) { indexes[i], indexes[j] = indexes[j], indexes[i] })
	for i := 0; i < mines; i++ {
		mineField.States[indexes[i]] = MaskMine
	}

	// Precalculate: walk through all blocks to populate the number of neighboring
	// mines. So when `Reveal`ing, these are already set.
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			currentBlockState := mineField.Get(i, j)
			if (currentBlockState & MaskMine) == MaskMine {
				// Stepped on a mine, skip
				continue
			}
			// 8 neighboring blocks
			neighborStates := mineField.getNeighborStates(i, j)
			numOfMines := 0
			for _, s := range neighborStates {
				if (s & MaskMine) == MaskMine {
					numOfMines++
				}
			}
			if numOfMines > 0 {
				mineField.Set(i, j, MaskNeighboringMines&BlockState(numOfMines))
			}
		}
	}
	return mineField, nil
}

func (m *MineField) Reveal(row, col int) BlockState {
	currentBlockState := m.Get(row, col)
	if (currentBlockState&MaskFlagged) == MaskFlagged ||
		(currentBlockState&MaskRevealed) == MaskRevealed ||
		(currentBlockState&MaskMine) == MaskMine ||
		currentBlockState == MaskInvalid {
		return currentBlockState
	}
	// Marks the current block as `revealed`
	m.Set(row, col, currentBlockState|MaskRevealed)
	neighboringMines := currentBlockState & MaskNeighboringMines
	if neighboringMines > 0 {
		return neighboringMines
	}
	// The current block is empty, we do a recursive `Reveal` (flood-fill)
	neighborStates := m.getNeighborStates(row, col)
	for i := 0; i < len(neighborStates); i++ {
		state := neighborStates[i]
		if (state & MaskMine) == 0 {
			nextRow, nextCol := row, col
			if i <= 2 {
				nextRow--
			} else if i >= 5 {
				nextRow++
			}
			if i == 0 || i == 3 || i == 5 {
				nextCol--
			} else if i == 2 || i == 4 || i == 7 {
				nextCol++
			}
			m.Reveal(nextRow, nextCol)
		}
	}
	return MaskEmpty
}

// RevealAndCheckWinOrLose returns 1 for win, 0 for nothing, -1 for lose
func (m *MineField) RevealAndCheckWinOrLose(row, col int) int {
	// Upon testing if we win or lose, We update the cells with rules
	// If we stepped on a mine, we mark it with an MaskInvalid, then reveal
	// all unflagged mines, while marking the wrongly flagged mines with
	// MaskInvalid
	state := m.Reveal(row, col)
	if state&MaskMine == MaskMine {
		// Lose
		m.RevealAllMines()
		m.Add(row, col, MaskInvalid)
		m.MarkWronglyFlaggedMines()
		return -1
	}

	if m.AreUnrevealedAllMines() {
		// Win
		m.FlagAllMines()
		return 1
	}
	return 0
}

func (m *MineField) Get(row, col int) BlockState {
	if row < 0 || row >= m.Rows || col < 0 || col >= m.Cols {
		return MaskInvalid
	}
	return m.States[m.Cols*row+col]
}

func (m *MineField) Set(row, col int, state BlockState) {
	m.States[m.Cols*row+col] = state
}

// Add adds upon current state (OR'd)
func (m *MineField) Add(row, col int, state BlockState) {
	m.States[m.Cols*row+col] |= state
}

func (m *MineField) RevealAllMines() {
	for i := 0; i < len(m.States); i++ {
		state := m.States[i]
		// Do not reveal those flagged correctly.
		if state&MaskMine == MaskMine && state&MaskFlagged != MaskFlagged {
			m.States[i] |= MaskRevealed
		}
	}
}

// FlagAllMines will mark all mine cells wit Flags. Useful for winning condition.
func (m *MineField) FlagAllMines() {
	for i := 0; i < len(m.States); i++ {
		if m.States[i]&MaskMine == MaskMine {
			m.States[i] |= MaskFlagged
		}
	}
	m.Flags = 0
}

func (m *MineField) AreUnrevealedAllMines() bool {
	for i := 0; i < len(m.States); i++ {
		state := m.States[i]
		if (state&MaskRevealed) == 0 && (state&MaskMine != MaskMine) {
			return false
		}
	}
	return true
}

func (m *MineField) MarkWronglyFlaggedMines() {
	for i := 0; i < len(m.States); i++ {
		state := m.States[i]
		if (state&MaskFlagged) == MaskFlagged && (state&MaskMine != MaskMine) {
			m.States[i] |= MaskInvalid
		}
	}
}

// Remove removes upon current state
func (m *MineField) Remove(row, col int, state BlockState) {
	m.States[m.Cols*row+col] &^= state
}

func (m *MineField) Flag(row, col int) BlockState {
	state := m.Get(row, col)
	if state&(MaskFlagged|MaskQuestioned) == MaskEmpty {
		if m.Flags > 0 {
			m.Flags -= 1
			m.Add(row, col, MaskFlagged)
		}
	} else if state&MaskFlagged == MaskFlagged {
		if m.Flags < m.Mines {
			m.Flags += 1
			m.Remove(row, col, MaskFlagged)
			m.Add(row, col, MaskQuestioned)
		}
	} else if state&MaskQuestioned == MaskQuestioned {
		m.Remove(row, col, MaskQuestioned)
	}
	return m.Get(row, col)
}

func (m *MineField) getNeighborStates(row, col int) NeighborStates {
	return NeighborStates{
		// 0, 1, 2
		m.Get(row-1, col-1), m.Get(row-1, col), m.Get(row-1, col+1),
		// 3, ( ), 4
		m.Get(row, col-1) /*    CurrentBlock    */, m.Get(row, col+1),
		m.Get(row+1, col-1), m.Get(row+1, col), m.Get(row+1, col+1),
		// 5, 6, 7,
	}
}

func DefaultMineField() *MineField {
	m, err := NewMineField(
		9, 9, 10,
	)
	if err != nil {
		panic(err.Error())
	}
	return m
}
