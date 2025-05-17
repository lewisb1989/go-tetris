package tetris

import (
	"fmt"
	"math/rand"
)

type Grid struct {
	layout [][]int
}

func NewGrid(layout [][]int) *Grid {
	return &Grid{
		layout: layout,
	}
}

type Piece struct {
	id    int
	x     int
	y     int
	shape [][]int
}

func NewPiece(id int, shape [][]int) *Piece {
	return &Piece{
		id:    id,
		x:     0,
		y:     0,
		shape: shape,
	}
}

func (p *Piece) Height() int {
	return len(p.shape)
}

func (p *Piece) Width() int {
	return len(p.shape[0])
}

type Tetris struct {
	height      int
	width       int
	activePiece *Piece
	archive     []*Piece
	grid        *Grid
	shapes      [][][]int
}

func NewTetris(width int, height int) *Tetris {
	if width < 6 {
		panic("minimum width is 6")
	}
	if height < 10 {
		panic("minimum height is 10")
	}
	shapes := [][][]int{
		{
			{0, 1, 0},
			{1, 1, 1},
		},
		{
			{1, 1},
			{1, 1},
		},
		{
			{0, 1, 1},
			{1, 1, 0},
		},
		{
			{1, 1, 0},
			{0, 1, 1},
		},
		{
			{1, 1, 1, 1},
		},
		{
			{1, 0, 0},
			{1, 1, 1},
		},
		{
			{0, 0, 1},
			{1, 1, 1},
		},
	}
	var layout [][]int
	for i := 0; i < height; i++ {
		layout = append(layout, []int{})
		for j := 0; j < width; j++ {
			layout[i] = append(layout[i], 0)
		}
	}
	tetris := &Tetris{
		height:  height,
		width:   width,
		archive: make([]*Piece, 0),
		grid:    NewGrid(layout),
		shapes:  shapes,
	}
	tetris.NewActivePiece()
	return tetris
}

func (t *Tetris) NewActivePiece() {
	id := len(t.archive) + 1
	shapeId := rand.Intn(len(t.shapes) - 1)
	shape := t.shapes[shapeId]
	x := rand.Intn(len(t.grid.layout[0]) - len(shape[0]) - 1)
	for i, _ := range shape {
		for j, _ := range shape[i] {
			if shape[i][j] > 0 {
				shape[i][j] = id
			}
		}
	}
	t.activePiece = NewPiece(id, shape)
	t.activePiece.x = x
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
	for i := 0; i < newRows; i++ {
		var row []int
		for j := 0; j < t.width; j++ {
			row = append(row, 0)
		}
		newLayout = append([][]int{row}, newLayout...)
	}
	t.grid.layout = newLayout
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
	fmt.Println("*** latest grid ***")
	layout := make([][]int, 0)
	for i, row := range t.grid.layout {
		layout = append(layout, []int{})
		for j, _ := range row {
			layout[i] = append(layout[i], t.grid.layout[i][j])
		}
	}
	grid := NewGrid(layout)
	t.AddPieceToGrid(grid, t.activePiece)
	for _, row := range layout {
		fmt.Println(row)
	}
}

func (t *Tetris) CollisionDetection(x int, y int, width int, height int) bool {
	if y+height > len(t.grid.layout) {
		return true
	}
	subset := t.grid.layout[y+height-1][x : x+width]
	bottomRowPiece := t.activePiece.shape[len(t.activePiece.shape)-1]
	for i, _ := range subset {
		if subset[i]+bottomRowPiece[i] > t.activePiece.id {
			return true
		}
	}
	return false
}
