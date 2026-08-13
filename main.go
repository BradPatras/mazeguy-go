package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type tickMsg struct{}

func main() {
	p := tea.NewProgram(&model{})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	didSetInitalSize bool
	width            int
	height           int
	grid             [][]cell
	gmodel           *genmodel
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			if !m.gmodel.isFinished {
				m.grid = m.gmodel.generateMaze(0)
			}
		}
	case tea.WindowSizeMsg:
		m.setWindowSize(msg.Width, msg.Height)
	case tickMsg:

	}

	return m, cmd
}

func (m *model) setWindowSize(w int, h int) {
	if !m.didSetInitalSize {
		m.height = ensuringOdd(h)
		m.width = ensuringOdd(w)

		m.didSetInitalSize = true
		m.gmodel = initializeMaze(((m.width)/4)-1, ((m.height)/2)-1)
		m.grid = m.gmodel.generateMaze(1)
	}
}

func ensuringOdd(v int) int {
	if v%2 != 0 {
		return v - 1
	} else {
		return v
	}
}

func (m *model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Loading...")
	}

	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#6cacd4"))
	v := tea.NewView(style.Width(m.width).MaxWidth(m.width).Render(render(m.grid)))
	v.AltScreen = true
	return v
}
