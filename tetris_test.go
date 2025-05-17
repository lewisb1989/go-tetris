package tetris

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTetris_MinGridWidth(t *testing.T) {
	assert.PanicsWithValue(t, "minimum width is 6", func() {
		_ = NewTetris(5, 10)
	})
}

func TestTetris_MinGridHeight(t *testing.T) {
	assert.PanicsWithValue(t, "minimum height is 10", func() {
		_ = NewTetris(10, 9)
	})
}

func TestTetris_ClearCompletedRows(t *testing.T) {
	tetris := NewTetris(6, 10)
	tetris.grid.layout = [][]int{
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 0},
		{0, 0, 0, 0, 2, 2},
		{7, 7, 7, 7, 2, 2},
		{0, 0, 0, 3, 3, 0},
		{0, 0, 0, 3, 3, 0},
		{6, 6, 4, 4, 4, 4},
		{6, 6, 5, 5, 5, 5},
	}
	tetris.ClearCompletedRows()
	assert.Equal(t, tetris.grid.layout, [][]int{
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 0},
		{0, 0, 0, 0, 2, 2},
		{0, 0, 0, 3, 3, 0},
		{0, 0, 0, 3, 3, 0},
	})
}

func TestTetris_PlayGame(t *testing.T) {
	tetris := NewTetris(6, 10)
	for i := 0; i < 10; i++ {
		tetris.MoveDown()
	}
	for i := 0; i < 10; i++ {
		tetris.MoveRight()
	}
	for i := 0; i < 10; i++ {
		tetris.MoveLeft()
	}
	for i := 0; i < 10; i++ {
		tetris.MoveDown()
	}
	for i := 0; i < 10; i++ {
		tetris.MoveLeft()
	}
	for i := 0; i < 10; i++ {
		tetris.MoveDown()
	}
	tetris.PrintGrid()
}
