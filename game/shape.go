package game

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
