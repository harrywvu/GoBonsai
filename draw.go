package main

import "fmt"

const (
	WIDTH = 100
	HEIGHT = 30
)

func DrawGrid(){

	// Create the buffer.
	grid := make([][]rune, HEIGHT)
	for y := range grid {
		grid[y] = make([]rune, WIDTH)
		for x := range grid[y] {
			grid[y][x] = '.' // empty cell
		}
	}

	// Place some characters.
	grid[0][0] = '@'
	grid[0][1] = '@'

	// Draw the whole grid.
	for _, row := range grid {
		fmt.Println(string(row))
	}
}