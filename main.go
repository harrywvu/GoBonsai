package main

import "fmt"

func main() {
	width, height := 100, 30

	// Create the buffer.
	grid := make([][]rune, height)
	for y := range grid {
		grid[y] = make([]rune, width)
		for x := range grid[y] {
			grid[y][x] = '.' // empty cell
		}
	}

	// Place some characters.
	grid[1][2] = '@'
	grid[3][7] = '#'

	// Draw the whole grid.
	for _, row := range grid {
		fmt.Println(string(row))
	}
}