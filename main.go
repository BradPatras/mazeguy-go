package main

import (
	"log"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/imacks/bitflags-go"
)

type wall int

const (
	top wall = 1 << iota
	right
	bottom
	left
)

type cell struct {
	walls   wall
	visited bool
}

var nilcell = cell{-1, false}

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
	m.height = ensuringOdd(h)
	m.width = ensuringOdd(w)

	if !m.didSetInitalSize {
		m.didSetInitalSize = true
		m.grid = createCellGrid((m.width-1)/2, (m.height-1)/2)
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
	v := tea.NewView(style.Width(m.width).MaxWidth(m.width).Render(m.render()))
	v.AltScreen = true
	return v
}

func createCellGrid(width int, height int) [][]cell {
	grid := make([][]cell, width)
	for x := range width {
		grid[x] = make([]cell, height)
		for y := range height {
			grid[x][y] = cell{
				walls: 15, // all walls, binary 1111
			}
		}
	}

	return grid
}

func (m model) render() string {
	var result strings.Builder
	for y := range m.height {
		for x := range m.width {
			cell := m.cellForPixel(x, y)
			if y == 0 && x == 0 { // top left corner
				result.WriteRune(rune('╔'))
			} else if y == 0 && x == m.width-1 { // top right corner
				result.WriteRune(rune('╗'))
			} else if y == m.height-1 && x == 0 { // bottom left corner
				result.WriteRune(rune('╚'))
			} else if y == m.height-1 && x == m.width-1 { // bottom right corner
				result.WriteRune(rune('╝'))
			} else if y == 0 || y == m.height-1 { // top or bottom edge
				if y == 0 && cell != nilcell && bitflags.Has(cell.walls, left, right) {
					result.WriteRune(rune('╦'))
				} else if y == m.height-1 && cell != nilcell && bitflags.Has(cell.walls, left, right) {
					result.WriteRune(rune('╩'))
				} else {
					result.WriteString(strconv.Itoa(int(cell.walls)))
				}
			} else if x == 0 || x == m.width-1 { // left or right edge
				// or ╠ or ╣
				result.WriteRune(rune('║'))
			} else {
				result.WriteRune(' ')
			}
		}
		result.WriteRune('\n')
	}
	return result.String()
}

func (m model) cellForPixel(x int, y int) cell {
	// cells always fall on odd-number coordinates
	if x%2 == 0 || y%2 == 0 {
		return nilcell
	}

	bottomleft := 

	return m.grid[x/2][y/2]
}

// Maze generates incrementally, shows status message, "Generating maze..."
// player character appears and status message switches to GO!
// Timer appears and starts counting up (or down?) and player is able to start moving

// maze generates to fill window size at time of initiation
