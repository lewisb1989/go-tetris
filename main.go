package main

import (
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	term "github.com/nsf/termbox-go"
	"io"
	"os"
	"tetris/game"
	"time"
)

func playAudio() {
	f, err := os.Open("tetris.mp3")
	if err != nil {
		panic(err)
	}
	d, err := mp3.NewDecoder(f)
	if err != nil {
		panic(err)
	}
	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		panic(err)
	}
	p := c.NewPlayer()
	if _, err := io.Copy(p, d); err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
	err = c.Close()
	if err != nil {
		panic(err)
	}
	err = p.Close()
	if err != nil {
		panic(err)
	}
}

func main() {
	err := term.Init()
	if err != nil {
		panic(err)
	}
	defer term.Close()
	tetris := game.NewTetris(10, 20, time.Millisecond*333)
	go playAudio()
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
