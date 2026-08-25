package main

import (
	"os"
)

func main() {
	args := os.Args[1:]
	options := getOptions(args)

	width, height := defaultWindowSize()
	width = min(width, options.WindowWidth)
	height = min(height, options.WindowHeight)

	w := NewWindow(options.WindowWidth, options.WindowHeight, options)

	// Draw a simple ground base
	groundRow := w.height - 3
	start := 40 - 25
	end := 40 + 25

	for x := start; x <= end; x++ {
		w.SetCharPlane(float64(groundRow), float64(x), '_', fixedColour(200, 150, 0))
	}

	// Draw a simple trunk column
	for y := float64(groundRow + 1); y <= float64(groundRow+3); y++ {
		w.SetCharPlane(float64(40), y, '|', fixedColour(200, 200, 200))
	}

	w.Draw()
	w.ResetCursor()
}