package main

import (
	"math/rand/v2"
	"slices"

	"github.com/imacks/bitflags-go"
)

func createCellGrid(width int, height int) [][]cell {
	grid := make([][]cell, width)
	for x := range width {
		grid[x] = make([]cell, height)
		for y := range height {
			grid[x][y] = cell{
				pos:   vec2{x, y},
				walls: bitflags.Set(0, top, left, right, bottom),
			}
		}
	}

	return grid
}

// generate a maze using a depth-first-search that visits all
// cells
func generateMaze(gridWidth int, gridHeight int) [][]cell {
	if gridHeight < 1 || gridWidth < 1 {
		panic("invalid maze dimensions")
	}

	grid := createCellGrid(gridWidth, gridHeight)

	// entry point is top left
	current := &grid[0][0]
	current.walls = bitflags.Del(current.walls, left)
	current.visited = true
	var backstack []vec2
	backstack = append(backstack, current.pos)

	// while backstack is not empty
	for len(backstack) > 0 {
		// pick a random un-visited neighbor of the last backstack element
		// collect unvisited neighbors
		leftPos := vec2{current.pos.x - 1, current.pos.y}
		rightPos := vec2{current.pos.x + 1, current.pos.y}
		topPos := vec2{current.pos.x, current.pos.y - 1}
		bottomPos := vec2{current.pos.x, current.pos.y + 1}
		var neighbors []vec2

		if leftPos.isInBounds(gridWidth, gridHeight) && !grid[leftPos.x][leftPos.y].visited {
			neighbors = append(neighbors, leftPos)
		}
		if rightPos.isInBounds(gridWidth, gridHeight) && !grid[rightPos.x][rightPos.y].visited {
			neighbors = append(neighbors, rightPos)
		}
		if topPos.isInBounds(gridWidth, gridHeight) && !grid[topPos.x][topPos.y].visited {
			neighbors = append(neighbors, topPos)
		}
		if bottomPos.isInBounds(gridWidth, gridHeight) && !grid[bottomPos.x][bottomPos.y].visited {
			neighbors = append(neighbors, bottomPos)
		}

		if len(neighbors) == 0 {
			// no valid neighbors, move backwards
			lastIndex := len(backstack) - 1
			backstack = slices.Delete(backstack, lastIndex, lastIndex+1)
			if len(backstack) < 1 {
				// backstack is empty, no more cells to visit
				break
			}
			previousPos := backstack[lastIndex-1]
			current = &grid[previousPos.x][previousPos.y]
			continue
		}

		// break down the wall between the backstack element and the neighbor
		pickedNeighborPos := neighbors[rand.IntN(len(neighbors))]
		pickedNeighbor := &grid[pickedNeighborPos.x][pickedNeighborPos.y]
		if pickedNeighborPos.x < current.pos.x {
			current.walls = bitflags.Del(current.walls, left)
			pickedNeighbor.walls = bitflags.Del(pickedNeighbor.walls, right)
		} else if pickedNeighborPos.x > current.pos.x {
			current.walls = bitflags.Del(current.walls, right)
			pickedNeighbor.walls = bitflags.Del(pickedNeighbor.walls, left)
		} else if pickedNeighborPos.y < current.pos.y {
			current.walls = bitflags.Del(current.walls, top)
			pickedNeighbor.walls = bitflags.Del(pickedNeighbor.walls, bottom)
		} else if pickedNeighborPos.y > current.pos.y {
			current.walls = bitflags.Del(current.walls, bottom)
			pickedNeighbor.walls = bitflags.Del(pickedNeighbor.walls, top)
		}

		// mark neighbor as visited and append it to the backstack
		pickedNeighbor.visited = true
		backstack = append(backstack, pickedNeighborPos)
		current = pickedNeighbor
	}

	// create exit point in bottom right
	grid[gridWidth-1][gridHeight-1].walls = bitflags.Del(grid[gridWidth-1][gridHeight-1].walls, right)

	return grid
}

func (v vec2) isInBounds(width int, height int) bool {
	return v.x > -1 && v.y > -1 && v.x < width && v.y < height
}
