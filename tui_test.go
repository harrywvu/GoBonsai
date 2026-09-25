package main

import (
	"math"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTUIStyleSelectionPreservesRandomAndSetsExplicitStyle(t *testing.T) {
	m := newOptionsTUIModel(&Options{Type: 2})
	if got := m.treeStyleSelection(); got != 0 {
		t.Fatalf("random selection = %d, want 0", got)
	}

	m.setTreeStyle(3)
	if !m.options.UserSetType || m.options.Type != 2 {
		t.Fatalf("explicit TUI selection = (%d, %t), want (2, true)", m.options.Type, m.options.UserSetType)
	}
}

func TestTUIAdjustsAndBoundsNumericValues(t *testing.T) {
	m := newOptionsTUIModel(&Options{NumLayers: 1, AngleMean: 40 * math.Pi / 180})
	m.cursor = layersField
	m.adjust(-1)
	if m.options.NumLayers != 1 {
		t.Fatalf("layers = %d, want lower bound of 1", m.options.NumLayers)
	}

	m.cursor = angleField
	m.adjust(1)
	if got := math.Round(m.options.AngleMean * 180 / math.Pi); got != 41 {
		t.Fatalf("angle = %g, want 41", got)
	}
}

func TestTUITextInputRejectsEmptyValue(t *testing.T) {
	m := newOptionsTUIModel(&Options{BranchChars: "~", LeafChars: "&"})
	m.cursor = branchCharsField
	m.input = ""
	m.editing = true
	updated, _ := m.updateTextInput(teaKeyEnter())
	got := updated.(optionsTUIModel)
	if !got.editing || got.error == "" {
		t.Fatal("empty character input should remain in edit mode with an error")
	}
}

func TestTUIArgumentHelpers(t *testing.T) {
	args := []string{"--tui", "--layers", "10"}
	if !hasTUIArg(args) {
		t.Fatal("--tui was not detected")
	}
	if got := withoutTUIArg(args); len(got) != 2 || got[0] != "--layers" || got[1] != "10" {
		t.Fatalf("withoutTUIArg() = %q, want [--layers 10]", got)
	}
	if !hasHelpOrVersionArg([]string{"--tui", "-h"}) {
		t.Fatal("-h was not detected as a help argument")
	}
}

func teaKeyEnter() tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyEnter}
}
