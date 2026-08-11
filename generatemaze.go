package main

import (
	"math"
	"math/rand/v2"
	"slices"

	"github.com/imacks/bitflags-go"
)

type genmodel struct {
	current    *cell
	backstack  []vec2
	grid       [][]cell
	gridWidth  int
	gridHeight int
	isFinished bool
}

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

func initializeMaze(gridWidth int, gridHeight int) *genmodel {
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

	return &genmodel{
		current:   current,
		backstack: backstack,
		grid:      grid,
		gridWidth: gridWidth,
		gridHeight: gridHeight,
	}
}

// generate maze using DFS with backtracking to visit every
// cell in a grid.
// param `iterations`: how many loops of the algorithm to run before returning. 0 value will render the whole maze before returning.
func (m *genmodel) generateMaze(iterations int) [][]cell {
	// iterations allow the maze to be generated bit-by-bit which
	// allows rendering the maze as it's being built
	if iterations == 0 {
		iterations = math.MaxInt
	}

	for range iterations {
		// pick a random un-visited neighbor of the last m.backstack element
		// collect unvisited neighbors
		leftPos := vec2{m.current.pos.x - 1, m.current.pos.y}
		rightPos := vec2{m.current.pos.x + 1, m.current.pos.y}
		topPos := vec2{m.current.pos.x, m.current.pos.y - 1}
		bottomPos := vec2{m.current.pos.x, m.current.pos.y + 1}
		var neighbors []vec2

		if leftPos.isInBounds(m.gridWidth, m.gridHeight) && !m.grid[leftPos.x][leftPos.y].visited {
			neighbors = append(neighbors, leftPos)
		}
		if rightPos.isInBounds(m.gridWidth, m.gridHeight) && !m.grid[rightPos.x][rightPos.y].visited {
			neighbors = append(neighbors, rightPos)
		}
		if topPos.isInBounds(m.gridWidth, m.gridHeight) && !m.grid[topPos.x][topPos.y].visited {
			neighbors = append(neighbors, topPos)
		}
		if bottomPos.isInBounds(m.gridWidth, m.gridHeight) && !m.grid[bottomPos.x][bottomPos.y].visited {
			neighbors = append(neighbors, bottomPos)
		}

		if len(neighbors) == 0 {
			// no valid neighbors, move backwards
			lastIndex := len(m.backstack) - 1
			m.backstack = slices.Delete(m.backstack, lastIndex, lastIndex+1)
			if len(m.backstack) < 1 {
				// m.backstack is empty, no more cells to visit
				break
			}
			previousPos := m.backstack[lastIndex-1]
			m.current = &m.grid[previousPos.x][previousPos.y]
			continue
		}

		// break down the wall between the m.backstack element and the neighbor
		pickedNeighborPos := neighbors[rand.IntN(len(neighbors))]
		pickedNeighbor := &m.grid[pickedNeighborPos.x][pickedNeighborPos.y]
		if pickedNeighborPos.x < m.current.pos.x {
			m.current.walls = bitflags.Del(m.current.walls, left)
			pickedNeighbor.walls = bitflags.Del(pickedNeighbor.walls, right)
		} else if pickedNeighborPos.x > m.current.pos.x {
			m.current.walls = bitflags.Del(m.current.walls, right)
			pickedNeighbor.walls = bitflags.Del(pickedNeighbor.walls, left)
		} else if pickedNeighborPos.y < m.current.pos.y {
			m.current.walls = bitflags.Del(m.current.walls, top)
			pickedNeighbor.walls = bitflags.Del(pickedNeighbor.walls, bottom)
		} else if pickedNeighborPos.y > m.current.pos.y {
			m.current.walls = bitflags.Del(m.current.walls, bottom)
			pickedNeighbor.walls = bitflags.Del(pickedNeighbor.walls, top)
		}

		// mark neighbor as visited and append it to the m.backstack
		pickedNeighbor.visited = true
		m.backstack = append(m.backstack, pickedNeighborPos)
		m.current = pickedNeighbor

		if len(m.backstack) == 0 {
			// create exit point in bottom right
			m.grid[m.gridWidth-1][m.gridHeight-1].walls = bitflags.Del(m.grid[m.gridWidth-1][m.gridHeight-1].walls, right)
			break
		}
	}

	return m.grid
}

func (v vec2) isInBounds(width int, height int) bool {
	return v.x > -1 && v.y > -1 && v.x < width && v.y < height
}
