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

	SnakeChar = 'o'

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

	gameOverFlag := false

	snk := Snake{[gameWinHeight * gameWinWidth]Point{}, DirectionRight, 3}
	snk.Segments[0] = Point{4, 2}
	snk.Segments[1] = Point{3, 2}
	snk.Segments[2] = Point{2, 2}

	processKeyEvent := func(kev *tcell.EventKey) {
		keyRune := kev.Rune()
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
		headCopy := *snk.Head()
		// Collides with walls?
		if headCopy.X == 0 || headCopy.X == w.Width || headCopy.Y == 0 || headCopy.Y == w.Height{
			gameOver()
		}

		// Collides with pellet?
		if headCopy == pelletPoint {
			snk.Increment()
			pelletPoint = getPellet(&w)
		}

		// Collides with itself?
		switch snk.Direction {
		case DirectionLeft:
			headCopy.X--
		case DirectionRight:
			headCopy.X++
		case DirectionUp:
			headCopy.Y--
		case DirectionDown:
			headCopy.Y++
		}
		ch, _, _ , _ := w.GetContentAtPoint(&headCopy, s)
		if ch == SnakeChar {
			gameOverFlag = true
		}
	}

	for {
		s.Show()
		if s.HasPendingEvent() {
			ev := s.PollEvent()
			processEvent(ev)
		}
		time.Sleep(tickTimeMS * time.Millisecond)
		collisionHandler()
		w.Clear(s)
		w.SetContent(pelletPoint.X, pelletPoint.Y, tcell.RuneDiamond, s)
		snk.Update()
		snk.Render(&w, s)
		if gameOverFlag {
			gameOver()
		}
	}

}
