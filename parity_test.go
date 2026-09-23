package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestDefaultsMatchPyBonsai(t *testing.T) {
	options := defaultOptions()

	if options.BranchChars != "~;:=" {
		t.Fatalf("BranchChars = %q, want %q", options.BranchChars, "~;:=")
	}
	if options.LeafChars != "&%#@" {
		t.Fatalf("LeafChars = %q, want %q", options.LeafChars, "&%#@")
	}
}

func TestParseArgsMatchesPythonDuplicateHandling(t *testing.T) {
	got := parseArgs([]string{"--seed", "1", "--type", "2", "--seed", "3", "-fi"})
	want := []arg{
		{name: "--seed", value: "3"},
		{name: "--type", value: "2"},
		{name: "-f", value: noValue},
		{name: "-i", value: noValue},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseArgs() = %#v, want %#v", got, want)
	}
}

func TestTreeRootUsesPythonIntegerCoordinates(t *testing.T) {
	options := &Options{WindowWidth: 25}
	if got, want := treeRoot(options), (Point{X: 12, Y: 8}); got != want {
		t.Fatalf("treeRoot() = %#v, want %#v", got, want)
	}
}

func TestBackgroundCellsRemainUncolouredInTerminal(t *testing.T) {
	wasTTY := stdoutIsTTY
	stdoutIsTTY = true
	t.Cleanup(func() { stdoutIsTTY = wasTTY })

	options := &Options{Instant: true}
	w := NewWindow(1, 1, options)
	if got := w.colourChar(w.cells[0][0]); got != " " {
		t.Fatalf("background cell rendered as %q, want a plain space", got)
	}

	w.SetCharScreen(0, 0, 'x', fixedColour(1, 2, 3))
	if got := w.colourChar(w.cells[0][0]); !strings.HasPrefix(got, "\x1b[38;2;1;2;3m") {
		t.Fatalf("drawn cell rendered without its colour: %q", got)
	}
}

func TestClippedLineCellStillConsumesColourRandomness(t *testing.T) {
	rng.Seed(1)
	options := &Options{Instant: true, BranchChars: defaultBranchChars}
	w := NewWindow(1, 1, options)

	w.setLineCell(0, 2, '|', branchColour)
	if got, want := rng.Float64(), 0.7609624449125756; got != want {
		t.Fatalf("next random value = %.17g, want %.17g", got, want)
	}
}

type recordingFibBrancher struct {
	called bool
}

func (r *recordingFibBrancher) drawEndBranches(Point, int, int, float64, float64, float64) {
	r.called = true
}

func TestFibBranchUsesSelectedVariant(t *testing.T) {
	options := &Options{
		NumLayers:   1,
		Instant:     true,
		BranchChars: defaultBranchChars,
	}
	w := NewWindow(20, 12, options)
	tree := newFibTree(w, Point{X: 10, Y: 8}, options)
	recorder := &recordingFibBrancher{}

	tree.drawBranch(recorder, tree.root, 1, 0, 5, 1, 0)
	if !recorder.called {
		t.Fatal("fib branch ignored the selected branching variant")
	}
}

func TestRandomOffsetBranchAtEndUsesFullLength(t *testing.T) {
	rng.Seed(1) // Python's first random() is below growEndThreshold.
	if got, atEnd := randomOffsetDistance(12); !atEnd || got != 12 {
		t.Fatalf("randomOffsetDistance() = (%g, %t), want (12, true)", got, atEnd)
	}
}
