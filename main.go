package main

import (
	"log"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/imacks/bitflags-go"
)

type wallFlags int

const (
	top wallFlags = 1 << iota
	right
	bottom
	left
)

type cell struct {
	walls   wallFlags
	visited bool
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	didSetInitalSize bool
	width            int
	height           int
	grid             [][]cell
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m = m.setWindowSize(msg.Width, msg.Height)
	}

	return m, cmd
}

func (m model) setWindowSize(w int, h int) model {
	if !m.didSetInitalSize {
		m.height = ensuringOdd(h)
		m.width = ensuringOdd(w)

		m.didSetInitalSize = true
		m.grid = createCellGrid((m.width-1)/4, (m.height-1)/2)
	}

	return m
}

func ensuringOdd(v int) int {
	if v%2 != 0 {
		return v - 1
	} else {
		return v
	}
}

func (m model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Loading...")
	}

	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#6cacd4"))
	v := tea.NewView(style.Width(m.width).MaxWidth(m.width).Render(m.render2()))
	v.AltScreen = true
	return v
}

func createCellGrid(width int, height int) [][]cell {
	grid := make([][]cell, width)
	for x := range width {
		grid[x] = make([]cell, height)
		for y := range height {
			grid[x][y] = cell{
				walls: 15, // all walls, bitwise flags 1111
			}
		}
	}

	return grid
}

// per-cell render
func (m model) render2() string {
	var result strings.Builder

	// cell height is 3 so we will build 3 lines of text at a time
	var sb1 strings.Builder
	var sb2 strings.Builder
	var sb3 strings.Builder
	gridHeight := len(m.grid[0])
	gridWidth := len(m.grid)

	for y := range gridHeight {
		for x := range gridWidth {
			n, e, s, w := m.getNeighborWalls(x, y)
			c := m.grid[x][y].walls
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
					sb2.WriteRune(rune(' '))
				}
			}
			sb2.WriteString("   ")
			if bitflags.Has(c, right) {
				sb2.WriteRune(rune('║'))
			} else {
				sb2.WriteRune(rune(' '))
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
func (m model) getNeighborWalls(cellX int, cellY int) (n wallFlags, e wallFlags, s wallFlags, w wallFlags) {
	// north neighbor
	if cellY == 0 { // top row
		n = bitflags.Set(n, bottom)
	} else {
		n = m.grid[cellX][cellY-1].walls
	}

	// south neighbor
	if cellY == len(m.grid[0])-1 { // bottom row
		s = bitflags.Set(s, top)
	} else {
		s = m.grid[cellX][cellY+1].walls
	}

	// east neighbor
	if cellX == len(m.grid)-1 { // right edge
		e = bitflags.Set(e, left)
	} else {
		e = m.grid[cellX+1][cellY].walls
	}

	// west neighbor
	if cellX == 0 { // left edge
		w = bitflags.Set(w, right)
	} else {
		w = m.grid[cellX-1][cellY].walls
	}

	return n, e, s, w
}

// Maze generates incrementally, shows status message, "Generating maze..."
// player character appears and status message switches to GO!
// Timer appears and starts counting up (or down?) and player is able to start moving

// maze generates to fill window size at time of initiation
