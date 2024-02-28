package main

import (
	"github.com/gdamore/tcell/v2"
)

type Snake struct {
	Segments  [gameWinHeight * gameWinWidth]Point
	Direction int
	Length    int
}

func (s *Snake) Head() *Point {
	return &s.Segments[0]
}

func (s *Snake) Update() {
	head := *s.Head()

	switch s.Direction {
	case DirectionLeft:
		head.X--
	case DirectionRight:
		head.X++
	case DirectionUp:
		head.Y--
	case DirectionDown:
		head.Y++
	}

	for i := s.Length - 1; i > 0; i-- {
		s.Segments[i] = s.Segments[i - 1]
	}

	s.Segments[0] = head
}

func (snk *Snake) Render(w *Window, s tcell.Screen) {
	w.SetContentAtPoints(snk.Segments[:snk.Length], SnakeChar, s)
}

func (s *Snake) Increment() {
	newPoint := s.Segments[s.Length - 1]
	newPoint.X--
	s.Segments[s.Length] = newPoint
	s.Length++
}
