package game

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"math/rand"
	"sync"
	"time"
)

// Shape encapsulates the color and rotations for a given shape
type Shape struct {
	rotations [][][]int
	color     string
}

// NewShape creates a new shape with color and rotations
func NewShape(color string, rotations [][][]int) *Shape {
	return &Shape{
		rotations: rotations,
		color:     color,
	}
}

// Shapes slice of all shapes in each possible rotation
var Shapes = []*Shape{
	NewShape("#ffff00", [][][]int{
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
	}),
	NewShape("#00ffff", [][][]int{
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
	}),
	NewShape("#BF40BF", [][][]int{
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
	}),
	NewShape("#00ff00", [][][]int{
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
	}),
	NewShape("#ff0000", [][][]int{
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
	}),
	NewShape("#0096FF", [][][]int{
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
	}),
	NewShape("#ff7ff0", [][][]int{
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
	}),
}

// setShapeId every time a new piece is created an auto-incrementing ID is assigned to that piece, so that
// collisions can be detected after adding the piece to the grid
//
// This function assigns the specified ID to the shape matrix by iterating over the rows and columns
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

// clearStdout clear the content displayed in stdout
func clearStdout() {
	fmt.Print("\033[H\033[2J")
}

// Grid encapsulates the layout of the current grid
type Grid struct {
	layout [][]int
}

// NewGrid create a new grid with the specified layout
func NewGrid(layout [][]int) *Grid {
	return &Grid{
		layout: layout,
	}
}

// Piece encapsulates the id, x/y co-ordinates, shape and rotation of a given piece
type Piece struct {
	id       int
	x        int
	y        int
	index    int
	rotation int
	shape    [][]int
}

// NewPiece creates a new piece with specified id, index, rotation and shape
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
	return piece
}

// Height return the height of this piece
func (p *Piece) Height() int {
	return len(p.shape)
}

// Width return the width of this piece
func (p *Piece) Width() int {
	return len(p.shape[0])
}

type Tetris struct {
	activeScore int
	scores      []int
	height      int
	width       int
	colors      map[int]string
	activePiece *Piece
	archive     []*Piece
	grid        *Grid
	shapes      []*Shape
	lock        sync.RWMutex
}

// StartNewGame starts a new game by initializing an empty grid, resetting the active score, and
// creating a new active piece
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
	tetris.colors = make(map[int]string)
	tetris.newActivePiece()
}

// NewTetris creates a new instance of the Tetris game with specified grid size and game speed
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
		colors:  make(map[int]string),
	}
	StartNewGame(tetris)
	go func() {
		for range time.NewTicker(gameSpeed).C {
			tetris.MoveDown()
		}
	}()
	return tetris
}

// gameOver clears stdout and prints the game over message
//
// After a short delay, we start a new game
func (t *Tetris) gameOver() {
	clearStdout()
	fmt.Println("GAME OVER! Score =", t.activeScore)
	fmt.Println()
	time.Sleep(time.Second * 2)
	fmt.Printf("New game in... ")
	for i := 0; i < 3; i++ {
		fmt.Printf("%d... ", 3-i)
		time.Sleep(time.Second * 1)
	}
	StartNewGame(t)
}

// newActivePiece randomly select a new active piece at a random x co-ordinate and with a random rotation
//
// If the creation of this piece results in a collision then the grid is full and a new game is started
func (t *Tetris) newActivePiece() {
	id := len(t.archive) + 1
	index := rand.Intn(len(t.shapes) - 1)
	t.colors[id] = t.shapes[index].color
	rotation := rand.Intn(3)
	shape := t.shapes[index].rotations[rotation]
	x := rand.Intn(len(t.grid.layout[0]) - len(shape[0]) - 1)
	t.activePiece = NewPiece(id, index, rotation, shape)
	t.activePiece.x = x
	if t.isCollisionDetected(x, 0, t.activePiece.shape, t.activePiece.id) {
		t.gameOver()
	}
}

// clearCompletedRows removes any rows in the grid that have pieces in all x co-ordinates, updates
// the active score for this game, and then adds new empty rows to the top of the grid to replace
// the completed rows that were removed
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

// Rotate rotates the active piece by 90 degrees as long as rotating the piece would not result in a collision
// with the grid boundaries or with another piece that is currently in place
//
// If the width of the piece increases after rotation, then it may be possible to rotate the piece after shifting
// to the left
//
// If a collision is detected for the current x and y co-ordinates, then we try to shift the piece to the left
// and check if rotating at the new co-ordinates is free from collisions
//
// When it is impossible to rotate the piece without causing a collision, then this method returns so that the
// piece is not rotated
func (t *Tetris) Rotate() {
	t.lock.Lock()
	defer t.lock.Unlock()
	rotation := t.activePiece.rotation + 1
	if rotation > 3 {
		rotation = 0
	}
	newShape := t.shapes[t.activePiece.index].rotations[rotation]
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
	t.printGrid()
}

// MoveDown moves the active piece down, so long as there are no collisions with the
// grid boundaries or with other pieces that are already in place
//
// If a collision is detected then the active piece can be considered to have reached the
// bottom of the grid, and the grid is updated and any full rows are removed from the grid
//
// # After clearing the completed rows, a new active piece is generated at the top of the grid
//
// Finally, print the updated grid to stdout
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

// MoveLeft moves the active piece to the left, so long as there are no collisions with the
// grid boundaries or with other pieces that are already in place
//
// Finally, print the updated grid to stdout
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

// MoveRight moves the active piece to the right, so long as there are no collisions with the
// grid boundaries or with other pieces that are already in place
//
// Finally, print the updated grid to stdout
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

// addPieceToGrid adds a piece to the grid
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

// UpdateGrid adds the active piece to the grid and then prints the updated grid to stdout
func (t *Tetris) UpdateGrid() {
	t.addPieceToGrid(t.grid, t.activePiece)
	t.printGrid()
	t.archive = append(t.archive, t.activePiece)
}

// printGrid makes a copy of the grid layout, adds the active piece to the grid, and then
// prints this grid to stdout
func (t *Tetris) printGrid() {
	clearStdout()
	fmt.Printf("   ╔╦╗╔═╗╔╦╗╦═╗╦╔═╗\n    ║ ║╣  ║ ╠╦╝║╚═╗\n    ╩ ╚═╝ ╩ ╩╚═╩╚═╝\n")
	hr := " ―"
	for i := 0; i < (t.width * 2); i++ {
		hr = hr + "―"
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
	for _, row := range grid.layout {
		fmt.Printf("| ")
		for _, cell := range row {
			if cell > 0 {
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(t.colors[cell])).Bold(true)
				fmt.Printf(style.Render("■"))
			} else {
				fmt.Printf(" ")
			}
			fmt.Printf(" ")
		}
		fmt.Printf("|\n")
	}
	fmt.Println(hr)
	fmt.Println(fmt.Sprintf(" Score: %d", t.activeScore))
	fmt.Println(" --")
	fmt.Println(" Press Esc to exit")
}

// isCollisionDetected checks if a shape at the given co-ordinates would collide with the boundaries
// of the grid, or if it collides with another piece that has already been added to the grid
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
