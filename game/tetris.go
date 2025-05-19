package game

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Tetris encapsulates the state of the game
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
	startNewGame(tetris)
	go func() {
		for range time.NewTicker(gameSpeed).C {
			tetris.MoveDown()
		}
	}()
	return tetris
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
	newPiece := NewPiece(t.activePiece.id, t.activePiece.index, rotation, newShape)
	if t.isCollisionDetected(t.activePiece.x, t.activePiece.y, newPiece) {
		if t.activePiece.Height() > t.activePiece.Width() {
			shift := t.activePiece.Height() - t.activePiece.Width()
			if t.isCollisionDetected(t.activePiece.x-shift, t.activePiece.y, newPiece) {
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
	t.grid.Print(t.activePiece, t.colors, t.activeScore)
}

// MoveDown moves the active piece down, so long as there are no collisions with the
// grid boundaries or with other pieces that are already in place
//
// If a collision is detected then the active piece can be considered to have reached the
// bottom of the grid, the grid is updated and any full rows are removed from the grid
func (t *Tetris) MoveDown() {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.isCollisionDetected(t.activePiece.x, t.activePiece.y+1, t.activePiece) {
		t.updateGrid()
		t.grid.ClearCompletedRows(func(scoreIncrement int) {
			t.activeScore += scoreIncrement
		})
		t.newActivePiece()
	} else {
		t.activePiece.y += 1
	}
	t.grid.Print(t.activePiece, t.colors, t.activeScore)
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
	if t.isCollisionDetected(t.activePiece.x-1, t.activePiece.y, t.activePiece) {
		return
	}
	t.activePiece.x -= 1
	t.grid.Print(t.activePiece, t.colors, t.activeScore)
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
	if t.isCollisionDetected(t.activePiece.x+1, t.activePiece.y, t.activePiece) {
		return
	}
	t.activePiece.x += 1
	t.grid.Print(t.activePiece, t.colors, t.activeScore)
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
	startNewGame(t)
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
	if t.isCollisionDetected(x, 0, t.activePiece) {
		t.gameOver()
	}
}

// UpdateGrid adds the active piece to the grid and then prints the updated grid to stdout
func (t *Tetris) updateGrid() {
	t.grid.AddPiece(t.activePiece)
	t.grid.Print(t.activePiece, t.colors, t.activeScore)
	t.archive = append(t.archive, t.activePiece)
}

// isCollisionDetected checks if a shape at the given co-ordinates would collide with the boundaries
// of the grid, or if it collides with another piece that has already been added to the grid
func (t *Tetris) isCollisionDetected(x int, y int, piece *Piece) bool {
	if x < 0 {
		return true
	}
	if x+piece.Width() > len(t.grid.layout[0]) {
		return true
	}
	if y+piece.Height() > len(t.grid.layout) {
		return true
	}
	var subset [][]int
	for i := y; i < y+piece.Height(); i++ {
		subset = append(subset, []int{})
		for j := x; j < x+piece.Width(); j++ {
			subset[i-y] = append(subset[i-y], t.grid.layout[i][j])
		}
	}
	for j := range piece.shape {
		activeRow := piece.shape[j]
		subsetRow := subset[j]
		for i := range subsetRow {
			if subsetRow[i]+activeRow[i] > piece.id {
				return true
			}
		}
	}
	return false
}
