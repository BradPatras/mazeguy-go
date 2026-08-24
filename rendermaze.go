package main

import (
	"strings"

	"github.com/imacks/bitflags-go"
)

const START_RUNE = rune('➬')
const END_RUNE = rune('➫')

// per-cell render
func render(grid [][]cell) string {
	var result strings.Builder

	// cell height is 3 so we will build 3 lines of text at a time
	var sb1 strings.Builder
	var sb2 strings.Builder
	var sb3 strings.Builder
	gridHeight := len(grid[0])
	gridWidth := len(grid)

	for y := range gridHeight {
		for x := range gridWidth {
			n, e, s, w := getNeighborWalls(x, y, grid)
			c := grid[x][y].walls
			// top row
			if y == 0 {
				if x == 0 {
					sb1.WriteRune(getWallIntersection(bitflags.Has(n, left), bitflags.Has(c, top), bitflags.Has(c, left), bitflags.Has(w, top)))
				}
				if bitflags.Has(c, top) {
					sb1.WriteString("═══")
				} else {
					sb1.WriteString("   ")
				}
				sb1.WriteRune(getWallIntersection(bitflags.Has(n, right), bitflags.Has(e, top), bitflags.Has(c, right), bitflags.Has(c, top)))
			}

			// middle row
			if x == 0 {
				if bitflags.Has(c, left) {
					sb2.WriteRune(rune('║'))
				} else {
					sb2.WriteRune(rune(START_RUNE))
				}
			}
			sb2.WriteString("   ")
			if bitflags.Has(c, right) {
				sb2.WriteRune(rune('║'))
			} else {
				if x == gridWidth-1 {
					sb2.WriteRune(rune(END_RUNE))
				} else {
					sb2.WriteRune(rune(' '))
				}
			}

			// bottom row
			if x == 0 {
				sb3.WriteRune(getWallIntersection(bitflags.Has(c, left), bitflags.Has(c, bottom), bitflags.Has(s, left), bitflags.Has(w, bottom)))
			}
			if bitflags.Has(c, bottom) {
				sb3.WriteString("═══")
			} else {
				sb3.WriteString("   ")
			}
			sb3.WriteRune(getWallIntersection(bitflags.Has(c, right), bitflags.Has(e, bottom), bitflags.Has(s, right), bitflags.Has(c, bottom)))
		}

		if sb1.Len() > 0 {
			sb1.WriteRune('\n')
			result.WriteString(sb1.String())
			sb1.Reset()
		}

		sb2.WriteRune('\n')
		result.WriteString(sb2.String())
		sb2.Reset()

		sb3.WriteRune('\n')
		result.WriteString(sb3.String())
		sb3.Reset()
	}

	return result.String()
}

// get the wall intersection rune based on the presence of neighboring walls
func getWallIntersection(n bool, e bool, s bool, w bool) rune {
	if !n && !e && !s && !w {
		return rune(' ')
	} else if !n && !e && s && w {
		return rune('╗')
	} else if !n && e && s && !w {
		return rune('╔')
	} else if !n && e && s && w {
		return rune('╦')
	} else if n && !e && !s && w {
		return rune('╝')
	} else if n && !e && s && w {
		return rune('╣')
	} else if n && e && !s && !w {
		return rune('╚')
	} else if n && e && !s && w {
		return rune('╩')
	} else if n && e && s && !w {
		return rune('╠')
	} else if n && e && s && w {
		return rune('╬')
	} else if n || s {
		return rune('║')
	} else if e || w {
		return rune('═')
	}

	return rune('?')
}

// return the wallsets of the neighboring cells
func getNeighborWalls(cellX int, cellY int, grid [][]cell) (n wallFlags, e wallFlags, s wallFlags, w wallFlags) {
	// north neighbor
	if cellY == 0 { // top row
		n = bitflags.Set(n, bottom)
	} else {
		n = grid[cellX][cellY-1].walls
	}

	// south neighbor
	if cellY == len(grid[0])-1 { // bottom row
		s = bitflags.Set(s, top)
	} else {
		s = grid[cellX][cellY+1].walls
	}

	// east neighbor
	if cellX == len(grid)-1 { // right edge
		e = bitflags.Set(e, left)
	} else {
		e = grid[cellX+1][cellY].walls
	}

	// west neighbor
	if cellX == 0 { // left edge
		w = bitflags.Set(w, right)
	} else {
		w = grid[cellX-1][cellY].walls
	}

	return n, e, s, w
}
