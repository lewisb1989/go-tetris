package main

import (
	term "github.com/nsf/termbox-go"
	"tetris/game"
	"time"
)

func main() {
	err := term.Init()
	if err != nil {
		panic(err)
	}
	defer term.Close()
	tetris := game.NewTetris(10, 20, time.Millisecond*333)
	for {
		ev := term.PollEvent()
		if ev.Key == term.KeyEsc {
			break
		} else if ev.Key == term.KeyArrowUp {
			tetris.RotateClockwise()
		} else if ev.Key == term.KeyArrowDown {
			tetris.MoveDown()
		} else if ev.Key == term.KeyArrowLeft {
			tetris.MoveLeft()
		} else if ev.Key == term.KeyArrowRight {
			tetris.MoveRight()
		}
	}
}
