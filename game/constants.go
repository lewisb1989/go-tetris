package game

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
