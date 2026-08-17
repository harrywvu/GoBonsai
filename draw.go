package main

import (
	"fmt"
	"math/rand/v2"
	"math"	
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
	CHAR_HEIGHT = 2
	CHAR_WIDTH = 1
)

func (w Window) planeToScreen(x,y float64) (row, col int){
	row = round(w.Height - y / CHAR_HEIGHT) // y/2 because 1 unit = 2 rows
	col = round(x / CHAR_WIDTH)				// x/1 because 1 unit = 1 column
}

func (w Window) screenToPlane(row, col int) (x,y float64){
	x = col * CHAR_WIDTH
	y = (w.Height - row) * CHAR_HEIGHT
}

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

// func getEnd(col, row int, length float64, angleDegrees float64) (int, int) {
//     rad := angleDegrees * math.Pi / 180
//     endCol := col + int(math.Round(length * math.Sin(rad)))
//     endRow := row - int(math.Round(length * math.Cos(rad)))
//     return endCol, endRow
// }