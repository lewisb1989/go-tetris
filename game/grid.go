package game

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
