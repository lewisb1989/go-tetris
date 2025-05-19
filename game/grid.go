package game

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

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

// ClearCompletedRows removes any rows in the grid that have pieces in all x co-ordinates, updates
// the active score for this game, and then adds new empty rows to the top of the grid to replace
// the completed rows that were removed
func (g *Grid) ClearCompletedRows(
	updateScore func(int),
) {
	var newLayout [][]int
	for _, row := range g.layout {
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
	newRows := len(g.layout) - len(newLayout)
	updateScore(newRows)
	for i := 0; i < newRows; i++ {
		var row []int
		for j := 0; j < len(g.layout[0]); j++ {
			row = append(row, 0)
		}
		newLayout = append([][]int{row}, newLayout...)
	}
	g.layout = newLayout
}

// AddPiece adds a piece to the grid
func (g *Grid) AddPiece(
	piece *Piece,
) {
	row := 0
	for i := piece.y; i < piece.y+piece.Height(); i++ {
		col := 0
		for j := piece.x; j < piece.x+piece.Width(); j++ {
			if piece.shape[row][col] > 0 {
				g.layout[i][j] = piece.shape[row][col]
			}
			col++
		}
		row++
	}
}

// Print makes a copy of the grid layout, adds the active piece to the grid, and then
// prints this grid to stdout
func (g *Grid) Print(
	activePiece *Piece,
	colors map[int]string,
	score int,
) {
	clearStdout()
	fmt.Printf("   ╔╦╗╔═╗╔╦╗╦═╗╦╔═╗\n    ║ ║╣  ║ ╠╦╝║╚═╗\n    ╩ ╚═╝ ╩ ╩╚═╩╚═╝\n")
	hr := " ―"
	for i := 0; i < (len(g.layout[0]) * 2); i++ {
		hr = hr + "―"
	}
	fmt.Println(hr)
	layout := make([][]int, 0)
	for i, row := range g.layout {
		layout = append(layout, []int{})
		for j := range row {
			layout[i] = append(layout[i], g.layout[i][j])
		}
	}
	grid := NewGrid(layout)
	if activePiece != nil {
		grid.AddPiece(activePiece)
	}
	for _, row := range grid.layout {
		fmt.Printf("| ")
		for _, cell := range row {
			if cell > 0 {
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(colors[cell])).Bold(true)
				fmt.Printf(style.Render("■"))
			} else {
				fmt.Printf(" ")
			}
			fmt.Printf(" ")
		}
		fmt.Printf("|\n")
	}
	fmt.Println(hr)
	fmt.Println(fmt.Sprintf(" Score: %d", score))
	fmt.Println(" --")
	fmt.Println(" Press Esc to exit")
}
