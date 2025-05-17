package game

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var Characters = []string{
	"a", "b", "c", "d", "e", "f",
	"g", "h", "i", "j", "k", "l",
	"m", "n", "o", "p", "q", "r",
	"s", "t", "u", "v", "w", "x",
	"y", "z",
}

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

func setShapeId(shape [][]int, id int) [][]int {
	for i := range shape {
		for j := range shape[i] {
			if shape[i][j] > 0 {
				shape[i][j] = id
			}
		}
	}
	return shape
}

func clearStdout() {
	fmt.Print("\033[H\033[2J")
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
	piece.shape = setShapeId(piece.shape, id)
	for i := range piece.shape {
		for j := range piece.shape[i] {
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
	lock        sync.RWMutex
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
	tetris.newActivePiece()
}

func NewTetris(
	width int,
	height int,
	gameSpeed time.Duration,
) *Tetris {
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
	go func() {
		for range time.NewTicker(gameSpeed).C {
			tetris.MoveDown()
		}
	}()
	return tetris
}

func (t *Tetris) gameOver() {
	clearStdout()
	fmt.Println("GAME OVER! Score =", t.activeScore)
	fmt.Println()
	time.Sleep(time.Second * 3)
	fmt.Printf("New game in... ")
	for i := 0; i < 3; i++ {
		fmt.Printf("%d... ", 3-i)
		time.Sleep(time.Second * 2)
	}
	StartNewGame(t)
}

func (t *Tetris) newActivePiece() {
	id := len(t.archive) + 1
	index := rand.Intn(len(t.shapes) - 1)
	rotation := rand.Intn(3)
	shape := t.shapes[index][rotation]
	x := rand.Intn(len(t.grid.layout[0]) - len(shape[0]) - 1)
	t.activePiece = NewPiece(id, index, rotation, shape)
	t.activePiece.x = x
	if t.isCollisionDetected(x, 0, t.activePiece.shape, t.activePiece.id) {
		t.gameOver()
	}
}

func (t *Tetris) clearCompletedRows() {
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

func (t *Tetris) Rotate() {
	t.lock.Lock()
	defer t.lock.Unlock()
	rotation := t.activePiece.rotation + 1
	if rotation > 3 {
		rotation = 0
	}
	newShape := t.shapes[t.activePiece.index][rotation]
	newShape = setShapeId(newShape, t.activePiece.id)
	if t.isCollisionDetected(t.activePiece.x, t.activePiece.y, newShape, t.activePiece.id) {
		if t.activePiece.Height() > t.activePiece.Width() {
			shift := t.activePiece.Height() - t.activePiece.Width()
			if t.isCollisionDetected(t.activePiece.x-shift, t.activePiece.y, newShape, t.activePiece.id) {
				return
			} else {
				t.activePiece.x -= shift
			}
		} else {
			return
		}
	}
	t.activePiece.rotation = rotation
	t.activePiece.shape = newShape
	t.activePiece.shape = setShapeId(t.activePiece.shape, t.activePiece.id)
	t.printGrid()
}

func (t *Tetris) MoveDown() {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.isCollisionDetected(t.activePiece.x, t.activePiece.y+1, t.activePiece.shape, t.activePiece.id) {
		t.UpdateGrid()
		t.clearCompletedRows()
		t.newActivePiece()
	} else {
		t.activePiece.y += 1
	}
	t.printGrid()
}

func (t *Tetris) MoveLeft() {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.activePiece.x-1 < 0 {
		return
	}
	if t.isCollisionDetected(t.activePiece.x-1, t.activePiece.y, t.activePiece.shape, t.activePiece.id) {
		return
	}
	t.activePiece.x -= 1
	t.printGrid()
}

func (t *Tetris) MoveRight() {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.activePiece.x+1+t.activePiece.Width() > len(t.grid.layout[0]) {
		return
	}
	if t.isCollisionDetected(t.activePiece.x+1, t.activePiece.y, t.activePiece.shape, t.activePiece.id) {
		return
	}
	t.activePiece.x += 1
	t.printGrid()
}

func (t *Tetris) addPieceToGrid(grid *Grid, piece *Piece) {
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
	t.addPieceToGrid(t.grid, t.activePiece)
	t.printGrid()
	t.archive = append(t.archive, t.activePiece)
}

func (t *Tetris) printGrid() {
	clearStdout()
	hr := "*"
	for i := 0; i < (t.width * 2); i++ {
		hr = hr + "*"
	}
	fmt.Println(hr)
	layout := make([][]int, 0)
	for i, row := range t.grid.layout {
		layout = append(layout, []int{})
		for j := range row {
			layout[i] = append(layout[i], t.grid.layout[i][j])
		}
	}
	grid := NewGrid(layout)
	t.addPieceToGrid(grid, t.activePiece)
	for _, row := range layout {
		var charRow []string
		for _, cell := range row {
			var char string
			if cell == 0 {
				char = " "
			} else {
				char = Characters[(cell-1)%len(Characters)]
			}
			charRow = append(charRow, char)
		}
		fmt.Println(charRow)
	}
	fmt.Println(hr)
	fmt.Println(fmt.Sprintf("Score: %d", t.activeScore))
}

func (t *Tetris) isCollisionDetected(x int, y int, shape [][]int, id int) bool {
	if x < 0 {
		return true
	}
	if x+len(shape[0]) > len(t.grid.layout[0]) {
		return true
	}
	if y+len(shape) > len(t.grid.layout) {
		return true
	}
	var subset [][]int
	for i := y; i < y+len(shape); i++ {
		subset = append(subset, []int{})
		for j := x; j < x+len(shape[0]); j++ {
			subset[i-y] = append(subset[i-y], t.grid.layout[i][j])
		}
	}
	for j := range shape {
		activeRow := shape[j]
		subsetRow := subset[j]
		for i := range subsetRow {
			if subsetRow[i]+activeRow[i] > id {
				return true
			}
		}
	}
	return false
}
