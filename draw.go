package main

import "fmt"

const (
	WIDTH = 100
	HEIGHT = 30
	H_CENTER int = WIDTH / 2 
	V_CEMTER int = HEIGHT / 2

	BASE_FLAT_LENGTH int = 25
)

func DrawBase(grid [][]rune){
	//
	//	\                             /
	//	 \							 /
//		  \_________________________/
//			--					 --
	

	start := H_CENTER - BASE_FLAT_LENGTH/2
	end := H_CENTER + BASE_FLAT_LENGTH/2

	for x := start; x <= end; x++ {
		grid[27][x] = '_'
	}
}

func DrawGrid(){

	// Create the buffer.
	grid := make([][]rune, HEIGHT)
	for y := range grid {
		grid[y] = make([]rune, WIDTH)
		for x := range grid[y] {
			grid[y][x] = ' ' // empty cell
		}
	}

	// Place some characters.
	// grid[0][0] = '@'
	// grid[0][1] = '@'

	DrawBase(grid)

	// Draw the whole grid.
	for _, row := range grid {
		fmt.Println(string(row))
	}
}