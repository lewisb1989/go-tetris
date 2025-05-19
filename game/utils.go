package game

import "fmt"

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

// StartNewGame starts a new game by initializing an empty grid, resetting the active score, and
// creating a new active piece
func startNewGame(tetris *Tetris) {
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
