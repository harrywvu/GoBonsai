package main

import "math"

type Point struct {
	X, Y float64
}

func (p Point) Add(o Point) Point {
	return Point{p.X + o.X, p.Y + o.Y}
}

func (p Point) Scale(s float64) Point {
	return Point{p.X * s, p.Y * s}
}

func (p Point) Mag() float64 {
	return math.Hypot(p.X, p.Y)
}

func (p *Point) Normalise() {
	m := p.Mag()
	if m > 0 {
		p.X /= m
		p.Y /= m
	}
}

type Line struct {
	Start, End  Point
	M, C        float64
	IsVertical  bool
}

func (l *Line) SetEndPoints(start, end Point) {
	l.Start = start
	l.End = end

	if start.X == end.X {
		l.IsVertical = true
		l.C = start.X
		return
	}

	l.M = (start.Y - end.Y) / (start.X - end.X)
	l.C = start.Y - l.M*start.X
}

func (l *Line) GetY(x float64) float64 {
	if l.IsVertical {
		return math.NaN()
	}
	return l.M*x + l.C
}

func (l *Line) GetX(y float64) float64 {
	if l.IsVertical {
		return l.C
	}
	if l.M != 0 {
		return (y - l.C) / l.M
	}
	return math.NaN()
}

func (l *Line) GetTheta() float64 {
	if l.IsVertical {
		return math.Pi / 2
	}
	return math.Atan(l.M)
}

func roundEven(x float64) int {
	return int(math.RoundToEven(x))
}