package main

import (
	"log"
	"strconv"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type tickMsg time.Time

func main() {
	p := tea.NewProgram(&model{iterationsPerFrame: 1})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	didSetInitalSize   bool
	width              int
	height             int
	grid               [][]cell
	gmodel             *genmodel
	isPaused           bool
	iterationsPerFrame int
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
			m.gmodel = initializeMaze(m.getMazeSizeInCells())
			m.grid = m.gmodel.generateMaze(m.iterationsPerFrame)
			cmd = tea.Batch(
				cmd,
				tickCmd(),
			)
		case "space":
			if len(m.grid) > 0 {
				m.isPaused = !m.isPaused
				if !m.isPaused {
					cmd = tea.Batch(
						cmd,
						tickCmd(),
					)
				}
			}
		case "left":
			switch m.iterationsPerFrame {
			case 0, 1:
				m.iterationsPerFrame = 0
			case 5:
				m.iterationsPerFrame = 1
			default:
				m.iterationsPerFrame -= 5
			}
		case "right":
			switch m.iterationsPerFrame {
			case 0:
				m.iterationsPerFrame = 1
			case 1:
				m.iterationsPerFrame = 5
			default:
				m.iterationsPerFrame += 5
			}
		}
	case tea.WindowSizeMsg:
		m.setWindowSize(msg.Width, msg.Height)
	case tickMsg:
		if !m.isPaused && !m.gmodel.isFinished {
			m.grid = m.gmodel.generateMaze(m.iterationsPerFrame)
			cmd = tea.Batch(
				cmd,
				tickCmd(),
			)
		}
	}

	return m, cmd
}

func (m *model) setWindowSize(w int, h int) {
	if !m.didSetInitalSize {
		m.didSetInitalSize = true
	}
	m.height = h
	m.width = w
}

func (m *model) getMazeSizeInCells() (int, int) {
	w, h := m.getMazePixelSize()
	return ((w) / 4) - 1, ((ensuringOdd(h)) / 2) - 1
}

func (m *model) getMazePixelSize() (int, int) {
	_, controlBarHeight := lipgloss.Size(m.controlView())
	hPadding := 4
	vPadding := 1
	return (m.width - hPadding), (ensuringOdd(m.height - controlBarHeight - vPadding))
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

	cv := m.controlView()

	_, cvHeight := lipgloss.Size(cv)

	var content string
	if len(m.grid) == 0 {
		content = lipgloss.NewStyle().Height(m.height - cvHeight).AlignVertical(lipgloss.Center).Render("Press enter to begin")
	} else {
		content = m.mazeView()
	}

	viewString := lipgloss.JoinVertical(
		lipgloss.Center,
		content,
		m.controlView(),
	)

	v := tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, viewString))
	v.AltScreen = true
	return v
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Nanosecond*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) mazeView() string {
	width, height := m.getMazePixelSize()
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#6cacd4"))
	return style.Width(width).MaxWidth(width).Height(height).MaxHeight(height).AlignVertical(lipgloss.Center).AlignHorizontal(lipgloss.Center).Render(render(m.grid))
}

func (m model) controlView() string {
	cmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6cacd4"))
	typeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#1ab753"))
	style := lipgloss.NewStyle().Faint(true)
	str := style.Render("Iterations per frame: ") +
		typeStyle.Render(strconv.Itoa(m.iterationsPerFrame)) +
		style.Render(" | ") +
		cmdStyle.Render("[left/right] arrow") +
		style.Render(" adjusts iterations per frame | ") +
		cmdStyle.Render("[enter]") +
		style.Render(" regenerates | ") +
		cmdStyle.Render("[space]") +
		style.Render(" play/pause")
	v := lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Width(m.width).
		BorderTop(true).
		Border(lipgloss.ThickBorder(), true, false, false, false).
		BorderForeground(lipgloss.Color("#666699")).
		Render(str)

	return v
}
