package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if hasTUIArg(args) && hasHelpOrVersionArg(args) {
		getOptions(withoutTUIArg(args))
		return
	}

	if shouldUseOptionsTUI(args) {
		if !isInteractiveTerminal() {
			fmt.Fprintln(os.Stderr, "--tui requires an interactive terminal")
			os.Exit(1)
		}

		options := getOptions(withoutTUIArg(args))
		var ok bool
		var err error
		options, ok, err = runOptionsTUI(options)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not start the options interface: %v\n", err)
			os.Exit(1)
		}
		if !ok {
			return
		}

		drawTree(options)
		return
	}

	drawTree(getOptions(args))
}

func hasHelpOrVersionArg(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" || arg == "--version" {
			return true
		}
	}
	return false
}

func drawTree(options *Options) {
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
