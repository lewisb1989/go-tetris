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

func loopThemeTune() {
	for {
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
		p.Close()
		c.Close()
		f.Close()
	}
}

func main() {
	go loopThemeTune()
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
			tetris.Rotate()
		} else if ev.Key == term.KeyArrowDown {
			tetris.MoveDown()
		} else if ev.Key == term.KeyArrowLeft {
			tetris.MoveLeft()
		} else if ev.Key == term.KeyArrowRight {
			tetris.MoveRight()
		}
	}
}
