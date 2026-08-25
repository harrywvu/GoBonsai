package main

import (
	"os"
	"golang.org/x/term"
)

func main() {
	args := os.Args[1:]
	options := getOptions(args)

	// Set up window sized for the terminal, clamped to defaults
	if term.IsTerminal(int(os.Stdout.Fd())) {
		wTerm, h, err := term.GetSize(int(os.Stdout.Fd()))
		if err == nil {
			options.WindowWidth = min(wTerm, 80)
			options.WindowHeight = min(h, 25)
		}
	}

	w := NewWindow(options.WindowWidth, options.WindowHeight, options)

	// Get a tree based on the selected type
	switch options.Type {
	case 0:
		t := newClassicTree(w, Point{float64(options.WindowWidth)/2, float64(boxHeight+4)}, options)
		t.draw()
	case 1:
		t := newFibTree(w, Point{float64(options.WindowWidth)/2, float64(boxHeight+4)}, options)
		t.draw()
	case 2:
		t := newOffsetFibTree(w, Point{float64(options.WindowWidth)/2, float64(boxHeight+4)}, options)
		t.draw()
	default:
		t := newRandomOffsetFibTree(w, Point{float64(options.WindowWidth)/2, float64(boxHeight+4)}, options)
		t.draw()
	}

	// Render to terminal and reset cursor
	w.Draw()
}