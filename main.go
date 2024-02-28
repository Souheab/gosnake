package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	gameWinHeight = 10
	gameWinWidth  = 25

	DirectionUp = iota
	DirectionDown
	DirectionLeft
	DirectionRight

	tickTimeMS = 300
)

type Point struct {
	X int
	Y int
}

func getPellet(w *Window) Point {

	x := rand.Intn(w.Width-1) + 1
	y := rand.Intn(w.Height-1) + 1

	return Point{x, y}
}

func main() {

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	sWidth, sHeight := s.Size()

	s.SetStyle(tcell.StyleDefault)

	quit := func() {
		s.Fini()
		os.Exit(0)
	}

	gameOver := func() {
		s.Fini()
		fmt.Println("GAME OVER!")
		os.Exit(0)
	}

	snk := Snake{[gameWinHeight * gameWinWidth]Point{}, DirectionRight, 3}
	snk.Segments[0] = Point{4, 2}
	snk.Segments[1] = Point{3, 2}
	snk.Segments[2] = Point{2, 2}

	processKeyEvent := func(kev *tcell.EventKey) {
		keyRune := kev.Rune()
		s.SetContent(0, 0, keyRune, nil, tcell.StyleDefault)
		switch keyRune {
		case 'q', 'Q':
			quit()
		case 'j', 'J':
			snk.Direction = DirectionDown
		case 'k', 'K':
			snk.Direction = DirectionUp
		case 'h', 'H':
			snk.Direction = DirectionLeft
		case 'l', 'L':
			snk.Direction = DirectionRight
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

	err = w.RenderBorder(s)
	if err != nil {
		s.Fini()
		panic(err)
	}

	pelletPoint := getPellet(&w)

	collisionHandler := func() {
		head := snk.Head()
		// Collides with walls?
		if head.X == 0 || head.X == w.Width || head.Y == 0 || head.Y == w.Height{
			gameOver()
		}

		if *head == pelletPoint {
			snk.Increment()
		}
	}

	for {
		s.Show()
		if s.HasPendingEvent() {
			ev := s.PollEvent()
			processEvent(ev)
		}
		time.Sleep(tickTimeMS * time.Millisecond)
		w.Clear(s)
		collisionHandler()
		snk.Update()
		w.SetContent(pelletPoint.X, pelletPoint.Y, tcell.RuneDiamond, s)
		snk.Render(&w, s)
	}

}
