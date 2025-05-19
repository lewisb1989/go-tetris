package game

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
