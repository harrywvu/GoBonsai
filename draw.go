package main

import (
	"fmt"
	"math/rand/v2"	
)

const (
	WIDTH = 100
	HEIGHT = 30
	H_CENTER int = WIDTH / 2 
	V_CEMTER int = HEIGHT / 2

	BASE_FLAT_LENGTH int = 25

	Reset  = "\033[0m"
	Green  = "\033[32m"
	White  = "\033[37m"

	
)

func DrawBase(grid [][]rune){
	//
	//	\_.^._~~..___..^..~~...***-_--/
	//	 \							 /
//		  \_________________________/
//			--					 --

	start := H_CENTER - BASE_FLAT_LENGTH/2
	end := H_CENTER + BASE_FLAT_LENGTH/2

	for x := start; x <= end; x++ {
		grid[27][x] = '_'
	}

	grid[28][39] = '-'
	grid[28][40] = '-'
	grid[28][61] = '-'
	grid[28][60] = '-'

	
	i := 1
	for y:= 27; y >= 25; y--{
		for x:= start - i; x <= end + i; x++{
			if x == start - i {
				grid[y][x] = '\\'
			}
			if x == end + i {
				grid[y][x] = '/'
			}
		}
		i += 1
	}

	chars := []rune{'.', ',', '~', '-', '_', '*', '^'}

	for x := start - 2; x <= end + 2; x++{
		grid[25][x] = chars[rand.IntN(len(chars))]
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