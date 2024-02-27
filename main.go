package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"log"
	"math/rand"
	"os"
)

type Point struct {
	X int
	Y int
}

type Window struct {
	Point  Point
	Width  int
	Height int
}

func (w *Window) Draw(s tcell.Screen) error {
	sWidth, sHeight := s.Size()

	x2 := w.Point.X + w.Width
	y2 := w.Point.Y + w.Height

	if w.Point.X > sWidth || w.Point.Y > sHeight || w.Point.X < 0 || w.Point.Y < 0 || x2 > sWidth || y2 > sHeight {
		return fmt.Errorf("%+v struct is invalid. Screen width, height are %v, %v", w, sWidth, sHeight)
	}

	for col := w.Point.X; col <= x2; col++ {
		s.SetContent(col, w.Point.Y, tcell.RuneHLine, nil, tcell.StyleDefault)
		s.SetContent(col, y2, tcell.RuneHLine, nil, tcell.StyleDefault)
	}

	for row := w.Point.Y; row <= y2; row++ {
		s.SetContent(w.Point.X, row, tcell.RuneVLine, nil, tcell.StyleDefault)
		s.SetContent(x2, row, tcell.RuneVLine, nil, tcell.StyleDefault)
	}

	if w.Point.Y != x2 && w.Point.Y != y2 {
		s.SetContent(w.Point.X, w.Point.Y, tcell.RuneULCorner, nil, tcell.StyleDefault)
		s.SetContent(x2, w.Point.Y, tcell.RuneURCorner, nil, tcell.StyleDefault)
		s.SetContent(w.Point.X, y2, tcell.RuneLLCorner, nil, tcell.StyleDefault)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, tcell.StyleDefault)
	}

	return nil
}

func (w *Window) SetContent(x int, y int, ch rune, s tcell.Screen) {
	s.SetContent(w.Point.X+x, w.Point.Y+y, ch, nil, tcell.StyleDefault)
}

func getPellet(w *Window) Point {

	x := rand.Intn(w.Width) 
	y := rand.Intn(w.Height)

	return Point{x, y}
}

func main() {
	const (
		gameWinHeight = 10
		gameWinWidth  = 25
	)

	const (
		directionUp    = iota
		directionDown  = iota
		directionLeft  = iota
		directionRight = iota
	)

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	sWidth, sHeight := s.Size()

	s.SetStyle(tcell.StyleDefault)

	type snake struct {
		Segments  [gameWinHeight * gameWinWidth]Point
		Direction int
		Length    int
	}

	//snk := snake{[gameWinHeight * gameWinWidth]Point{}, directionRight, 1}

	quit := func() {
		s.Fini()
		os.Exit(0)
	}

	processKeyEvent := func(kev *tcell.EventKey) {
		keyRune := kev.Rune()
		s.SetContent(0, 0, keyRune, nil, tcell.StyleDefault)
		switch keyRune {
		case 'q', 'Q':
			quit()
		}
	}

	processEvent := func(ev tcell.Event) {
		switch ev := ev.(type) {
		case *tcell.EventKey:
			processKeyEvent(ev)
		}
	}

	pX := (sWidth / 2) - (gameWinWidth / 2)
	pY := (sHeight / 2) - (gameWinHeight / 2)
	w := Window{Point{pX, pY}, gameWinWidth, gameWinHeight}
	pelletPoint := getPellet(&w)
	w.SetContent(pelletPoint.X, pelletPoint.Y, tcell.RuneDiamond, s)

	for {
		s.Show()
		ev := s.PollEvent()
		processEvent(ev)
		err := w.Draw(s)
		if err != nil {
			s.Fini()
			panic(err)
		}
	}

}
