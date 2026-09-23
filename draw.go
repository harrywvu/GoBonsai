package main

import (
	"cmp"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
	"time"

	"golang.org/x/term"
)

const (
	charWidth      = 1.0
	charHeight     = 2.0
	backgroundChar = ' '
	charThreshold  = 0.3

	endColour  = "\033[00m"
	hideCursor = "\033[?25l"
	showCursor = "\033[?25h"
)

var stdoutIsTTY = term.IsTerminal(int(os.Stdout.Fd()))

type RGB struct {
	R, G, B int
}

type Colour struct {
	fixed    RGB
	ranged   [3][2]int
	isRanged bool
}

func fixedColour(r, g, b int) Colour {
	return Colour{fixed: RGB{r, g, b}}
}

func rangedColour(r, g, b [2]int) Colour {
	return Colour{ranged: [3][2]int{r, g, b}, isRanged: true}
}

func (c Colour) pick() RGB {
	if !c.isRanged {
		return c.fixed
	}
	return RGB{
		randInRange(c.ranged[0][0], c.ranged[0][1]),
		randInRange(c.ranged[1][0], c.ranged[1][1]),
		randInRange(c.ranged[2][0], c.ranged[2][1]),
	}
}

func randInRange(lower, upper int) int {
	return lower + rng.IntN(upper-lower+1)
}

type Cell struct {
	Ch       rune
	Fore     RGB
	Coloured bool
}

type Window struct {
	width   int
	height  int
	cells   [][]Cell
	options *Options
}

func NewWindow(width, height int, options *Options) *Window {
	cells := make([][]Cell, height)
	for y := range cells {
		cells[y] = newCellRow(width)
	}
	return &Window{width: width, height: height, cells: cells, options: options}
}

func newCellRow(width int) []Cell {
	row := make([]Cell, width)
	for x := range row {
		row[x] = Cell{Ch: backgroundChar, Fore: RGB{}}
	}
	return row
}

func (w *Window) colourChar(cell Cell) string {
	if !stdoutIsTTY || !cell.Coloured {
		return string(cell.Ch)
	}
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%c%s", cell.Fore.R, cell.Fore.G, cell.Fore.B, cell.Ch, endColour)
}

func (w *Window) PlaneToScreen(x, y float64) (int, int) {
	row := roundEven(float64(w.height) - y/charHeight)
	col := roundEven(x / charWidth)
	return row, col
}

func (w *Window) ScreenToPlane(row, col int) (float64, float64) {
	x := float64(col) * charWidth
	y := float64(w.height-row) * charHeight
	return x, y
}

func (w *Window) IncreaseHeight(delta int) bool {
	if w.options.FixedWindow {
		return false
	}

	for range delta {
		w.cells = append([][]Cell{newCellRow(w.width)}, w.cells...)
	}
	w.height += delta

	return true
}

func (w *Window) SetCharScreen(row, col int, ch rune, colour Colour) {
	if row < 0 && w.IncreaseHeight(-row) {
		row = 0
	}

	if row < 0 || row >= w.height || col < 0 || col >= w.width {
		return
	}

	w.cells[row][col] = Cell{Ch: ch, Fore: colour.pick(), Coloured: true}
}

func (w *Window) SetCharPlane(x, y float64, ch rune, colour Colour) {
	row, col := w.PlaneToScreen(x, y)
	w.SetCharScreen(row, col, ch, colour)
}

func (w *Window) SetCharPlaneWait(x, y float64, ch rune, colour Colour) {
	w.SetCharPlane(x, y, ch, colour)
	w.wait()
}

func (w *Window) wait() {
	if w.options.Instant {
		return
	}
	w.Draw()
	time.Sleep(time.Duration(w.options.WaitTime * float64(time.Second)))
}

func (w *Window) Draw() {
	var out strings.Builder

	if stdoutIsTTY {
		out.WriteString(hideCursor)
	}

	for _, row := range w.cells {
		for _, cell := range row {
			out.WriteString(w.colourChar(cell))
		}
		out.WriteByte('\n')
	}

	if stdoutIsTTY {
		fmt.Fprintf(&out, "\033[%dA", w.height)
		out.WriteString(showCursor)
	}

	fmt.Print(out.String())
}

func (w *Window) ResetCursor() {
	if stdoutIsTTY {
		fmt.Printf("\033[%dB", w.height)
	}
}

func getLineChar(theta float64) rune {
	upper := math.Pi / 2 * (2.0 / 3.0)
	lower := math.Pi / 2 * (1.0 / 3.0)

	switch {
	case math.Abs(theta) > upper:
		return '|'
	case math.Abs(theta) < lower:
		return '_'
	case theta > 0:
		return '/'
	default:
		return '\\'
	}
}

type cellDist struct {
	dist float64
	inx  int
}

func (w *Window) DrawLine(start, end Point, colour Colour, width int) {
	mid := &Line{}
	mid.SetEndPoints(start, end)

	char := getLineChar(mid.GetTheta())

	w.checkLineBounds(start, end)

	if mid.IsVertical || math.Abs(mid.M) >= 1 {
		w.drawSteepLine(start, end, colour, width, char, mid)
	} else {
		w.drawShallowLine(start, end, colour, width, char, mid)
	}
}

func (w *Window) checkLineBounds(start, end Point) {
	h1, _ := w.PlaneToScreen(start.X, start.Y)
	h2, _ := w.PlaneToScreen(end.X, end.Y)

	if roomFromTop := min(h1, h2); roomFromTop < 0 {
		w.IncreaseHeight(-roomFromTop)
	}
}

func (w *Window) drawSteepLine(start, end Point, colour Colour, width int, char rune, mid *Line) {
	startRow, _ := w.PlaneToScreen(start.X, start.Y)
	endRow, _ := w.PlaneToScreen(end.X, end.Y)

	step := 1
	if endRow < startRow {
		step = -1
	}

	for row := startRow; ; row += step {
		dists := make([]cellDist, w.width)

		for col := 0; col < w.width; col++ {
			x, y := w.ScreenToPlane(row, col)
			dists[col] = cellDist{dist: math.Abs(mid.GetX(y) - x), inx: col}
		}

		sortByDist(dists)

		for i := 0; i < width && i < len(dists); i++ {
			w.setLineCell(row, dists[i].inx, char, colour)
		}

		if row == endRow {
			break
		}
	}
}

func (w *Window) drawShallowLine(start, end Point, colour Colour, width int, char rune, mid *Line) {
	_, startCol := w.PlaneToScreen(start.X, start.Y)
	_, endCol := w.PlaneToScreen(end.X, end.Y)

	step := 1
	if endCol < startCol {
		step = -1
	}

	for col := startCol; ; col += step {
		dists := make([]cellDist, w.height)

		for row := 0; row < w.height; row++ {
			x, y := w.ScreenToPlane(row, col)
			dists[row] = cellDist{dist: math.Abs(mid.GetY(x) - y), inx: row}
		}

		sortByDist(dists)

		for i := 0; i < width && i < len(dists); i++ {
			w.setLineCell(dists[i].inx, col, char, colour)
		}

		if col == endCol {
			break
		}
	}
}

func sortByDist(dists []cellDist) {
	slices.SortFunc(dists, func(a, b cellDist) int {
		return cmp.Or(cmp.Compare(a.dist, b.dist), cmp.Compare(a.inx, b.inx))
	})
}

func (w *Window) setLineCell(row, col int, char rune, colour Colour) {
	ch := char
	if rng.Float64() < charThreshold {
		ch = randChoice(w.options.BranchChars)
	}
	picked := colour.pick()
	w.SetCharScreen(row, col, ch, fixedColour(picked.R, picked.G, picked.B))
	w.wait()
}

func randChoice(s string) rune {
	chars := []rune(s)
	return chars[rng.IntN(len(chars))]
}
