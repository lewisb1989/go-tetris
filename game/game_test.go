package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTetris_Rotate(t *testing.T) {
	tetris := NewTetris(6, 10, time.Minute)
	tetris.activePiece.shape = tetris.shapes[4][0]
	tetris.activePiece.index = 0
	tetris.activePiece.rotation = 0
	tetris.Rotate()
	assert.Equal(t, tetris.activePiece.shape, tetris.shapes[0][1])
	tetris.Rotate()
	assert.Equal(t, tetris.activePiece.shape, tetris.shapes[0][2])
	tetris.Rotate()
	assert.Equal(t, tetris.activePiece.shape, tetris.shapes[0][3])
	tetris.Rotate()
	assert.Equal(t, tetris.activePiece.shape, tetris.shapes[0][0])
}

func TestTetris_RotateWithoutShift(t *testing.T) {
	tetris := NewTetris(6, 10, time.Minute)
	tetris.activePiece.shape = [][]int{
		{6},
		{6},
		{6},
		{6},
	}
	tetris.activePiece.index = 4
	tetris.activePiece.rotation = 1
	tetris.activePiece.x = 5
	tetris.activePiece.y = 4
	tetris.grid.layout = [][]int{
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 4, 4, 0},
		{0, 0, 0, 4, 4, 0},
		{0, 0, 0, 3, 3, 0},
		{0, 0, 0, 3, 3, 0},
		{0, 0, 0, 2, 2, 0},
		{0, 0, 0, 2, 2, 0},
		{0, 0, 0, 1, 1, 0},
		{0, 0, 0, 1, 1, 0},
	}
	tetris.printGrid()
	tetris.Rotate()
	assert.Equal(t, tetris.activePiece.shape, [][]int{
		{6},
		{6},
		{6},
		{6},
	})
}

func TestTetris_RotateWithShift(t *testing.T) {
	tetris := NewTetris(6, 10, time.Minute)
	tetris.activePiece.shape = [][]int{
		{6},
		{6},
		{6},
		{6},
	}
	tetris.activePiece.id = 6
	tetris.activePiece.index = 4
	tetris.activePiece.rotation = 1
	tetris.activePiece.x = 5
	tetris.activePiece.y = 4
	tetris.grid.layout = [][]int{
		{0, 0, 0, 0, 0, 0},
		{5, 5, 5, 5, 0, 0},
		{2, 4, 0, 0, 0, 0},
		{2, 4, 0, 0, 0, 0},
		{2, 4, 0, 0, 0, 0},
		{2, 4, 0, 0, 0, 0},
		{1, 3, 0, 0, 0, 0},
		{1, 3, 0, 0, 0, 0},
		{1, 3, 0, 0, 0, 0},
		{1, 3, 0, 0, 0, 0},
	}
	tetris.printGrid()
	tetris.Rotate()
	assert.Equal(t, tetris.activePiece.shape, [][]int{
		{6, 6, 6, 6},
	})
}

func TestTetris_MinGridWidth(t *testing.T) {
	assert.PanicsWithValue(t, "minimum width is 6", func() {
		_ = NewTetris(5, 10, time.Minute)
	})
}

func TestTetris_MinGridHeight(t *testing.T) {
	assert.PanicsWithValue(t, "minimum height is 10", func() {
		_ = NewTetris(10, 9, time.Minute)
	})
}

func TestTetris_ClearCompletedRows(t *testing.T) {
	tetris := NewTetris(6, 10, time.Minute)
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
	tetris.clearCompletedRows()
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

func TestTetris_AddPieceToGrid(t *testing.T) {
	tetris := NewTetris(6, 10, time.Minute)
	shape := tetris.shapes[1][1]
	piece := NewPiece(1, 1, 1, shape)
	piece.x = 0
	piece.y = 0
	tetris.addPieceToGrid(tetris.grid, piece)
	assert.Equal(t, tetris.grid.layout, [][]int{
		{1, 1, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
	})
}

func TestTetris_IsCollisionDetected(t *testing.T) {
	tetris := NewTetris(6, 10, time.Minute)
	shape := tetris.shapes[1][1]
	assert.Equal(t, shape, [][]int{
		{1, 1},
		{1, 1},
	})
	piece := NewPiece(1, 1, 1, shape)
	piece.x = 0
	piece.y = 0
	tetris.activePiece = piece
	result := tetris.isCollisionDetected(0, 9, tetris.activePiece.shape, tetris.activePiece.id)
	assert.True(t, result)
	result = tetris.isCollisionDetected(0, 7, tetris.activePiece.shape, tetris.activePiece.id)
	assert.False(t, result)
	tetris.grid.layout = [][]int{
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{2, 2, 0, 0, 0, 0},
		{2, 2, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0},
	}
	result = tetris.isCollisionDetected(1, 5, tetris.activePiece.shape, tetris.activePiece.id)
	assert.True(t, result)
	result = tetris.isCollisionDetected(1, 4, tetris.activePiece.shape, tetris.activePiece.id)
	assert.False(t, result)
	result = tetris.isCollisionDetected(2, 5, tetris.activePiece.shape, tetris.activePiece.id)
	assert.False(t, result)
}

func TestTetris_MoveRight_MoveLeft(t *testing.T) {
	tetris := NewTetris(6, 10, time.Minute)
	tetris.activePiece = NewPiece(1, 0, 0, tetris.shapes[1][1])
	assert.Equal(t, tetris.activePiece.shape, [][]int{
		{1, 1},
		{1, 1},
	})
	for i := 0; i < 10; i++ {
		tetris.MoveRight()
	}
	tetris.addPieceToGrid(tetris.grid, tetris.activePiece)
	assert.Equal(t, tetris.grid.layout, [][]int{
		{0, 0, 0, 0, 1, 1},
		{0, 0, 0, 0, 1, 1},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
	})
	StartNewGame(tetris)
	tetris.activePiece = NewPiece(1, 0, 0, tetris.shapes[1][1])
	for i := 0; i < 10; i++ {
		tetris.MoveRight()
	}
	for i := 0; i < 10; i++ {
		tetris.MoveLeft()
	}
	tetris.addPieceToGrid(tetris.grid, tetris.activePiece)
	assert.Equal(t, tetris.grid.layout, [][]int{
		{1, 1, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
	})
}

func TestTetris_MoveDown(t *testing.T) {
	tetris := NewTetris(6, 10, time.Minute)
	tetris.activePiece = NewPiece(1, 0, 0, tetris.shapes[1][1])
	assert.Equal(t, tetris.activePiece.shape, [][]int{
		{1, 1},
		{1, 1},
	})
	for i := 0; i < 9; i++ {
		tetris.MoveDown()
	}
	assert.Equal(t, tetris.grid.layout, [][]int{
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0},
	})
}

func TestTetris_StartNewGame(t *testing.T) {
	tetris := NewTetris(6, 10, time.Minute)
	for i := 0; i < 100; i++ {
		tetris.MoveDown()
	}
	assert.Greater(t, len(tetris.scores), 0)
}
