package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)
type Window struct {
	Point  Point
	Width  int
	Height int
}

func (w *Window) RenderBorder(s tcell.Screen) error {
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

func (w *Window) SetContent(x , y int, ch rune, s tcell.Screen) {
	s.SetContent(w.Point.X+x, w.Point.Y+y, ch, nil, tcell.StyleDefault)
}

func (w *Window) SetContentAtPoint(p *Point, ch rune, s tcell.Screen) {
	w.SetContent(p.X, p.Y, ch, s)
}

func (w *Window) SetContentAtPoints(p []Point, ch rune, s tcell.Screen) {
	for _, value := range p {
		w.SetContentAtPoint(&value, ch, s)
	}
}

func (w *Window) GetContent(x, y int, s tcell.Screen) (primary rune, combining []rune, style tcell.Style, width int) {
	return s.GetContent(w.Point.X + x, w.Point.Y + y)
}

func (w *Window) GetContentAtPoint(p *Point, s tcell.Screen) (primary rune, combining []rune, style tcell.Style, width int) {
	return w.GetContent(p.X, p.Y, s )
}

func (w *Window) Clear(s tcell.Screen) {
	for x := range w.Width - 1 {
		for y := range w.Height - 1 {
			w.SetContent(x+1, y+1, ' ', s)
		}
	}
}
