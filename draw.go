package main

import (
	"fmt"
	"math/rand/v2"	
)

type Window struct{
	grid [][]rune
	Width int
	Height int
}

func NewWindow(w, h int) *Window{
	grid := make([][]rune, h)
	for y := range grid {
		grid[y] = make([]rune, w)
		for x:= range grid[y] {
			grid[y][x] = ' '
		}
	}
	return &Window{grid: grid, Width: w, Height: h}
}

const (
	BASE_FLAT_LENGTH int = 25
)

func DrawBase(w *Window, centerCol int){
	half := BASE_FLAT_LENGTH / 2
	start := centerCol - half
	end := centerCol + half

	groundRow := w.Height - 3

	for x := start; x <= end; x++ {
		w.SetChar(groundRow, x, '_')
	}

	w.SetChar(groundRow+1, centerCol-half+1, '-')
	w.SetChar(groundRow+1, centerCol-half+2, '-')
	w.SetChar(groundRow+1, centerCol+half-1, '-')
	w.SetChar(groundRow+1, centerCol+half-2, '-')
	
	w.SetChar(groundRow, centerCol-half-1, '\\')
	w.SetChar(groundRow, centerCol+half+1, '/')
	w.SetChar(groundRow-1, centerCol-half-2, '\\')
	w.SetChar(groundRow-1, centerCol+half+2, '/')
	w.SetChar(groundRow-2, centerCol-half-3, '\\')
	w.SetChar(groundRow-2, centerCol+half+3, '/')

	chars := []rune{'.', ',', '~', '-', '_', '*', '^'}

	for x := start - 2; x <= end + 2; x++{
		w.SetChar(groundRow-2, x, chars[rand.IntN(len(chars))])
	}

}

func (w *Window) SetChar(row, col int, ch rune){
	if row >= 0 && row < w.Height && col >= 0 && col < w.Width{
		w.grid[row][col] = ch
	}
}

func (w *Window) Render(){
	for _, row := range w.grid {
		fmt.Println(string(row))
	}
}