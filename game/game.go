package game

import (
	"fmt"
	"math/rand"
)

var Shapes = [][][][]int{
	{
		{
			{0, 1, 0},
			{1, 1, 1},
		},
		{
			{1, 0},
			{1, 1},
			{1, 0},
		},
		{
			{1, 1, 1},
			{0, 1, 0},
		},
		{
			{0, 1},
			{1, 1},
			{0, 1},
		},
	},
	{
		{
			{1, 1},
			{1, 1},
		},
		{
			{1, 1},
			{1, 1},
		},
		{
			{1, 1},
			{1, 1},
		},
		{
			{1, 1},
			{1, 1},
		},
	},
	{
		{
			{0, 1, 1},
			{1, 1, 0},
		},
		{
			{1, 0},
			{1, 1},
			{0, 1},
		},
		{
			{0, 1, 1},
			{1, 1, 0},
		},
		{
			{1, 0},
			{1, 1},
			{0, 1},
		},
	},
	{
		{
			{1, 1, 0},
			{0, 1, 1},
		},
		{
			{0, 1},
			{1, 1},
			{1, 0},
		},
		{
			{1, 1, 0},
			{0, 1, 1},
		},
		{
			{0, 1},
			{1, 1},
			{1, 0},
		},
	},
	{
		{
			{1, 1, 1, 1},
		},
		{
			{1},
			{1},
			{1},
			{1},
		},
		{
			{1, 1, 1, 1},
		},
		{
			{1},
			{1},
			{1},
			{1},
		},
	},
	{
		{
			{1, 0, 0},
			{1, 1, 1},
		},
		{
			{1, 1},
			{1, 0},
			{1, 0},
		},
		{
			{1, 1, 1},
			{0, 0, 1},
		},
		{
			{0, 1},
			{0, 1},
			{1, 1},
		},
	},
	{
		{
			{0, 0, 1},
			{1, 1, 1},
		},
		{
			{1, 0},
			{1, 0},
			{1, 1},
		},
		{
			{1, 1, 1},
			{1, 0, 0},
		},
		{
			{1, 1},
			{0, 1},
			{0, 1},
		},
	},
}

type Grid struct {
	layout [][]int
}

func NewGrid(layout [][]int) *Grid {
	return &Grid{
		layout: layout,
	}
}

type Piece struct {
	id       int
	x        int
	y        int
	index    int
	rotation int
	shape    [][]int
}

func NewPiece(id int, index int, rotation int, shape [][]int) *Piece {
	piece := &Piece{
		id:       id,
		x:        0,
		y:        0,
		index:    index,
		rotation: rotation,
		shape:    shape,
	}
	for i, _ := range piece.shape {
		for j, _ := range piece.shape[i] {
			if piece.shape[i][j] > 0 {
				piece.shape[i][j] = id
			}
		}
	}
	return piece
}

func (p *Piece) Height() int {
	return len(p.shape)
}

func (p *Piece) Width() int {
	return len(p.shape[0])
}

type Tetris struct {
	activeScore int
	scores      []int
	height      int
	width       int
	activePiece *Piece
	archive     []*Piece
	grid        *Grid
	shapes      [][][][]int
}

func StartNewGame(tetris *Tetris) {
	var layout [][]int
	for i := 0; i < tetris.height; i++ {
		layout = append(layout, []int{})
		for j := 0; j < tetris.width; j++ {
			layout[i] = append(layout[i], 0)
		}
	}
	tetris.scores = append(tetris.scores, tetris.activeScore)
	tetris.activeScore = 0
	tetris.grid = NewGrid(layout)
	tetris.archive = make([]*Piece, 0)
	tetris.NewActivePiece()
}

func NewTetris(width int, height int) *Tetris {
	if width < 6 {
		panic("minimum width is 6")
	}
	if height < 10 {
		panic("minimum height is 10")
	}
	tetris := &Tetris{
		height:  height,
		width:   width,
		archive: make([]*Piece, 0),
		shapes:  Shapes,
		scores:  make([]int, 0),
	}
	StartNewGame(tetris)
	return tetris
}

func (t *Tetris) NewActivePiece() {
	id := len(t.archive) + 1
	index := rand.Intn(len(t.shapes) - 1)
	rotation := rand.Intn(3)
	shape := t.shapes[index][rotation]
	x := rand.Intn(len(t.grid.layout[0]) - len(shape[0]) - 1)
	t.activePiece = NewPiece(id, index, rotation, shape)
	t.activePiece.x = x
	if t.CollisionDetection(x, 0, t.activePiece.Width(), t.activePiece.Height()) {
		StartNewGame(t)
	}
}

func (t *Tetris) ClearCompletedRows() {
	var newLayout [][]int
	for _, row := range t.grid.layout {
		zeroes := 0
		for _, cell := range row {
			if cell == 0 {
				zeroes++
			}
		}
		if zeroes > 0 {
			newLayout = append(newLayout, row)
		}
	}
	newRows := t.height - len(newLayout)
	t.activeScore += newRows
	for i := 0; i < newRows; i++ {
		var row []int
		for j := 0; j < t.width; j++ {
			row = append(row, 0)
		}
		newLayout = append([][]int{row}, newLayout...)
	}
	t.grid.layout = newLayout
}

func (t *Tetris) RotateClockwise() {
	t.activePiece.rotation += 1
	if t.activePiece.rotation > 3 {
		t.activePiece.rotation = 0
	}
	t.activePiece.shape = t.shapes[t.activePiece.index][t.activePiece.rotation]
}

func (t *Tetris) RotateCounterClockwise() {
	t.activePiece.rotation -= 1
	if t.activePiece.rotation < 0 {
		t.activePiece.rotation = 3
	}
	t.activePiece.shape = t.shapes[t.activePiece.index][t.activePiece.rotation]
}

func (t *Tetris) MoveDown() {
	if t.CollisionDetection(t.activePiece.x, t.activePiece.y+1, t.activePiece.Width(), t.activePiece.Height()) {
		t.UpdateGrid()
		t.ClearCompletedRows()
		t.NewActivePiece()
	} else {
		t.activePiece.y += 1
	}
}

func (t *Tetris) MoveLeft() {
	if t.activePiece.x-1 < 0 {
		return
	}
	t.activePiece.x -= 1
}

func (t *Tetris) MoveRight() {
	if t.activePiece.x+1+t.activePiece.Width() > len(t.grid.layout[0]) {
		return
	}
	t.activePiece.x += 1
}

func (t *Tetris) AddPieceToGrid(grid *Grid, piece *Piece) {
	row := 0
	for i := piece.y; i < piece.y+piece.Height(); i++ {
		col := 0
		for j := piece.x; j < piece.x+piece.Width(); j++ {
			if piece.shape[row][col] > 0 {
				grid.layout[i][j] = piece.shape[row][col]
			}
			col++
		}
		row++
	}
}

func (t *Tetris) UpdateGrid() {
	t.AddPieceToGrid(t.grid, t.activePiece)
	t.PrintGrid()
	t.archive = append(t.archive, t.activePiece)
}

func (t *Tetris) PrintGrid() {
	// FIXME: Make this work in the console with static logging
	fmt.Println("*** latest grid ***")
	layout := make([][]int, 0)
	for i, row := range t.grid.layout {
		layout = append(layout, []int{})
		for j := range row {
			layout[i] = append(layout[i], t.grid.layout[i][j])
		}
	}
	grid := NewGrid(layout)
	t.AddPieceToGrid(grid, t.activePiece)
	for _, row := range layout {
		var charRow []string
		for _, cell := range row {
			char := string(rune('a' - 2 + cell + 1))
			if cell == 0 {
				char = " "
			}
			charRow = append(charRow, char)
		}
		fmt.Println(charRow)
	}
}

func (t *Tetris) CollisionDetection(x int, y int, width int, height int) bool {
	if y+height > len(t.grid.layout) {
		return true
	}
	//subset := t.grid.layout[y+height-1][x : x+width]
	var subset [][]int
	for i := y; i < y+height; i++ {
		subset = append(subset, []int{})
		for j := x; j < x+width; j++ {
			subset[i-y] = append(subset[i-y], t.grid.layout[i][j])
		}
	}
	for j := range t.activePiece.shape {
		activeRow := t.activePiece.shape[j]
		subsetRow := subset[j]
		for i := range subsetRow {
			if subsetRow[i]+activeRow[i] > t.activePiece.id {
				return true
			}
		}
	}

	return false
}
