package main

import (
	"os"
)

func main() {
	args := os.Args[1:]
	options := getOptions(args)

	w := NewWindow(options.WindowWidth, options.WindowHeight, options)
	root := treeRoot(options)

	// Get a tree based on the selected type
	switch options.Type {
	case 0:
		t := newClassicTree(w, root, options)
		t.draw()
	case 1:
		t := newFibTree(w, root, options)
		t.draw()
	case 2:
		t := newOffsetFibTree(w, root, options)
		t.draw()
	default:
		t := newRandomOffsetFibTree(w, root, options)
		t.draw()
	}

	// Render to terminal and reset cursor
	w.Draw()
	w.ResetCursor()
}

func treeRoot(options *Options) Point {
	rootY := boxHeight + 4
	rootY += rootY % 2
	return Point{float64(options.WindowWidth / 2), float64(rootY)}
}
