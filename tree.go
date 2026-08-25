package main

import (
	"math"
)

const (
	boxHeight           = 3
	maxTopWidth         = 35
	moundThreshold      = 0.1
	soilCharThreshold   = 0.1
	soilChars           = ".~*"
	moundWidthMean      = 2.0
	moundWidthStdDev    = 1.0

	angleStdDev   = 8 * math.Pi / 180
	lenScale      = 0.75
	maxInitialWidth = 6

	meanBranches   = 2.0
	branchesStdDev = 0.5

	numLeaves = 4

	growEndThreshold = 0.5
	nonEndMin        = 0.3
	nonEndMax        = 0.9

	boxFeetChar = '‾'
)

var (
	soilColour      = fixedColour(0, 150, 0)
	boxColour       = fixedColour(200, 200, 200)
	trunkBaseColour = fixedColour(255, 255, 0)
	branchColour    = rangedColour([2]int{200, 255}, [2]int{150, 255}, [2]int{0, 0})
)

type Tree struct {
	window      *Window
	root        Point
	options     *Options
	boxTopWidth int
}

func newTree(window *Window, root Point, options *Options) Tree {
	return Tree{window: window, root: root, options: options, boxTopWidth: boxWidth(window.width)}
}

func boxWidth(windowWidth int) int {
	width := min(windowWidth/3, maxTopWidth)
	if width%2 == 0 {
		width++
	}
	return width
}

func (t *Tree) drawBox() {
	rootRow, rootCol := t.window.PlaneToScreen(t.root.X, t.root.Y)

	for i := range boxHeight {
		row := rootRow + i
		width := t.boxTopWidth - i*2

		for x := range width {
			col := rootCol - width/2 + x

			var ch rune
			colour := boxColour

			switch {
			case x == 0:
				ch = '\\'
			case x == width-1:
				ch = '/'
			case i == 0:
				ch = '_'
				colour = soilColour
			case i == boxHeight-1:
				ch = '_'
			default:
				if rng.Float64() < soilCharThreshold {
					ch = randChoice(soilChars)
				} else {
					ch = ' '
				}
				colour = soilColour
			}

			t.window.SetCharScreen(row, col, ch, colour)
		}
	}

	t.drawBoxFeet(rootRow, rootCol)
	t.drawAllMounds(rootRow, rootCol)
}

func (t *Tree) drawBoxFeet(rootRow, rootCol int) {
	row := rootRow + boxHeight
	offset := t.boxTopWidth/2 - boxHeight - 1

	for _, sign := range [2]int{-1, 1} {
		t.window.SetCharScreen(row, rootCol+sign*offset, boxFeetChar, boxColour)
	}
}

func (t *Tree) drawAllMounds(rootRow, rootCol int) {
	numDrawn := 0

	for i := 1; i < t.boxTopWidth; i++ {
		col := rootCol - t.boxTopWidth/2 + i

		if rng.Float64() < moundThreshold/float64(numDrawn+1) {
			numDrawn++
			t.drawMound(rootRow, col, t.boxTopWidth-i-1)
		}
	}
}

func (t *Tree) drawMound(row, startCol, maxWidth int) {
	topWidth := min(roundEven(rng.NormFloat64()*moundWidthStdDev+moundWidthMean), maxWidth-2)

	if topWidth <= 0 {
		return
	}

	for i := 0; i <= topWidth+1; i++ {
		ch := '-'
		if i == 0 || i == topWidth+1 {
			ch = '.'
		}
		t.window.SetCharScreen(row, startCol+i, ch, soilColour)
	}
}

func (t *Tree) DrawTreeBase(trunkWidth int) {
	rootRow, rootCol := t.window.PlaneToScreen(t.root.X, t.root.Y)

	leftX := rootCol - trunkWidth/2
	rightX := rootCol + trunkWidth/2

	if trunkWidth%2 == 0 {
		rightX--
	}

	t.window.SetCharScreen(rootRow, leftX-2, '.', trunkBaseColour)
	t.window.SetCharScreen(rootRow, leftX-1, '/', trunkBaseColour)
	t.window.SetCharScreen(rootRow, rightX+1, '\\', trunkBaseColour)
	t.window.SetCharScreen(rootRow, rightX+2, '.', trunkBaseColour)
}

type RecursiveTree struct {
	Tree

	angleStdDev   float64
	lenScale      float64
	maxInitialWidth int
}

func newRecursiveTree(window *Window, root Point, options *Options) RecursiveTree {
	return RecursiveTree{
		Tree:          newTree(window, root, options),
		angleStdDev:   angleStdDev,
		lenScale:      lenScale,
		maxInitialWidth: maxInitialWidth,
	}
}

func (t *RecursiveTree) getEndCoords(start Point, length, theta float64) Point {
	return Point{
		X: start.X + length*math.Sin(theta),
		Y: start.Y + length*math.Cos(theta),
	}
}

func (t *RecursiveTree) initialParams() (int, float64) {
	initialWidth := min(max(t.options.InitialLen/5, 0), t.maxInitialWidth)
	initialAngle := rng.NormFloat64() * t.angleStdDev

	return initialWidth, initialAngle
}

type ClassicTree struct {
	RecursiveTree
}

func newClassicTree(window *Window, root Point, options *Options) *ClassicTree {
	return &ClassicTree{RecursiveTree: newRecursiveTree(window, root, options)}
}

func (t *ClassicTree) drawBranch(p Point, layer int, length, width, theta float64) {
	if layer >= t.options.NumLayers {
		drawLeaves(t.window, p, t.options)
		return
	}

	end := t.getEndCoords(p, length, theta)

	t.window.DrawLine(p, end, branchColour, roundEven(width))

	t.drawEndBranches(p, layer, length, width, theta)
}

func (t *ClassicTree) drawEndBranches(start Point, layer int, length, width, theta float64) {
	sign := 1.0
	numBranches := max(roundEven(rng.NormFloat64()*branchesStdDev+meanBranches), 0)

	step := 0.0
	if numBranches != 0 {
		step = length / float64(numBranches)
	}

	newWidth := max(width-1, 1)
	newLength := length * t.lenScale

	for i := range numBranches {
		distUpBranch := float64(i+1) * step
		newTheta := theta + sign*(rng.NormFloat64()*t.angleStdDev+t.options.AngleMean)

		p := t.getEndCoords(start, distUpBranch, theta)

		t.drawBranch(p, layer+1, newLength, newWidth, newTheta)

		sign *= -1
	}
}

type fibTree struct {
	RecursiveTree

	fib        []int
	branchNums [][]int
}

func newFibTree(window *Window, root Point, options *Options) *fibTree {
	t := &fibTree{RecursiveTree: newRecursiveTree(window, root, options)}
	t.fib = fibNums(t.options.NumLayers)
	t.branchNums = generateBranchNums(t.fib, t.options.NumLayers)
	return t
}

func fibNums(numLayers int) []int {
	fib := []int{1, 1}
	for range numLayers {
		fib = append(fib, fib[len(fib)-1]+fib[len(fib)-2])
	}
	return fib
}

func generateBranchNums(fib []int, numLayers int) [][]int {
	branchNums := [][]int{{1}}

	for i := range numLayers {
		numParents := 0
		for _, n := range branchNums[len(branchNums)-1] {
			numParents += n
		}

		numBranches := fib[i+2]
		base := numBranches / numParents
		diff := numBranches - base*numParents

		currentNums := make([]int, numParents)
		for x := range currentNums {
			currentNums[x] = base
			if x < diff {
				currentNums[x] = base + 1
			}
		}

		shuffle(currentNums)

		branchNums = append(branchNums, currentNums)
	}

	return branchNums
}

func shuffle(arr []int) {
	for i := len(arr) - 1; i > 0; i-- {
		j := rng.IntN(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func (t *fibTree) drawBranch(p Point, layerInx, branchInx int, length, width, theta float64) {
	if layerInx > t.options.NumLayers {
		drawLeaves(t.window, p, t.options)
		return
	}

	end := t.getEndCoords(p, length, theta)

	t.window.DrawLine(p, end, branchColour, roundEven(width))

	t.drawEndBranches(p, layerInx, branchInx, length, width, theta)
}

func (t *fibTree) drawEndBranches(start Point, layerInx, branchInx int, length, width, theta float64) {
	sign := 1.0
	numBranches := t.branchNums[layerInx][branchInx]
	newWidth := max(width-1, 1)

	for i := range numBranches {
		angle := rng.NormFloat64()*t.angleStdDev + t.options.AngleMean
		newTheta := theta + sign*angle
		newLength := length * t.lenScale

		p := t.getEndCoords(start, float64(i+1)*length/float64(numBranches), theta)

		t.drawBranch(p, layerInx+1, branchInx+i, newLength, newWidth, newTheta)

		sign *= -1
	}
}

type OffsetFibTree struct {
	fibTree
}

func newOffsetFibTree(window *Window, root Point, options *Options) *OffsetFibTree {
	return &OffsetFibTree{fibTree: *newFibTree(window, root, options)}
}

func (t *OffsetFibTree) draw() {
	initialWidth, initialAngle := t.initialParams()

	t.drawBox()
	t.DrawTreeBase(initialWidth)

	t.drawBranch(t.root, 1, 0, float64(t.options.InitialLen), float64(initialWidth), initialAngle)
}

func (t *OffsetFibTree) drawEndBranches(start Point, layerInx, branchInx int, length, width, theta float64) {
	sign := 1.0
	numBranches := t.branchNums[layerInx][branchInx]

	step := 0.0
	if numBranches != 0 {
		step = length / float64(numBranches)
	}

	newWidth := max(width-1, 1)
	newLength := length * t.lenScale

	for i := range numBranches {
		distUpBranch := float64(i+1) * step
		newTheta := theta + sign*(rng.NormFloat64()*t.angleStdDev+t.options.AngleMean)

		p := t.getEndCoords(start, distUpBranch, theta)

		t.drawBranch(p, layerInx+1, branchInx+i, newLength, newWidth, newTheta)

		sign *= -1
	}
}

type RandomOffsetFibTree struct {
	fibTree
}

func newRandomOffsetFibTree(window *Window, root Point, options *Options) *RandomOffsetFibTree {
	return &RandomOffsetFibTree{fibTree: *newFibTree(window, root, options)}
}

func (t *RandomOffsetFibTree) draw() {
	initialWidth, initialAngle := t.initialParams()

	t.drawBox()
	t.DrawTreeBase(initialWidth)

	t.drawBranch(t.root, 1, 0, float64(t.options.InitialLen), float64(initialWidth), initialAngle)
}

func (t *RandomOffsetFibTree) drawEndBranches(start Point, layerInx, branchInx int, length, width, theta float64) {
	sign := 1.0
	numBranches := t.branchNums[layerInx][branchInx]

	newWidth := max(width-1, 1)
	newLength := length * t.lenScale

	_ = branchInx // silence unused variable check in loop below
	needLeaves := true
	var distUpBranch float64
	for i := range numBranches {
		_ = i // explicitly use loop variable
		growAtEnd := rng.Float64() < growEndThreshold

		if growAtEnd {
			needLeaves = false
		} else {
			distUpBranch = rng.Float64()*(nonEndMax-nonEndMin) + nonEndMin
			distUpBranch *= length
		}

		newTheta := theta + sign*(rng.NormFloat64()*t.angleStdDev+t.options.AngleMean)

		p := t.getEndCoords(start, distUpBranch, theta)

		t.drawBranch(p, layerInx+1, branchInx+i, newLength, newWidth, newTheta)

		sign *= -1
	}

	if needLeaves {
		endPos := t.getEndCoords(start, length, theta)
		drawLeaves(t.window, endPos, t.options)
	}
}

type Leaves struct {
	window *Window
	branch Point
	opts     *Options
}

func drawLeaves(w *Window, branchEnd Point, opts *Options) {
	g := Point{X: 0, Y: -1}

	for k := 0; k < numLeaves; k++ {
		vel := Point{X: rng.Float64()*2 - 1, Y: rng.Float64()*2 - 1}
		vel = vel.Normalise()

		pos := branchEnd

		for j := 0; j < opts.LeafLen; j++ {
			pos = pos.Add(vel)

			rgb := RGB{R: 0, G: 75 + rng.IntN(181), B: 0}
			ch := randChoice(opts.LeafChars)

			if opts.Instant {
				w.SetCharPlane(pos.X, pos.Y, ch, fixedColour(0, rgb.G, 0))
			} else {
				w.SetCharPlaneWait(pos.X, pos.Y, ch, fixedColour(0, rgb.G, 0))
			}

			weight := float64(j) / float64(opts.LeafLen)
			vel = vel.Add(g.Scale(weight))
		}
	}
}