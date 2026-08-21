package main

import (
	"log"
	"strconv"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type generatorTickMsg time.Time

func main() {
	p := tea.NewProgram(&model{iterationsPerFrame: 1})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

const TITLE_STRING = "   ════════╦═══════════╦═══════════╦═══════════════╗\n           ║           ║           ║               ║\n   ╔════   ║   ║   ════╣   ╔═══╗   ║   ║   ╔═══╗   ║\n   ║       ║   ║       ║   ║   ║       ║   ║   ║   ║\n   ║   ╔═══╣   ╚═══╗   ║   ║   ╚═══╦═══╝   ║   ║   ║\n   ║   ║   ║       ║               ║       ║   ║   ║\n   ║   ║   ╚════   ╠═══════════════╣   ════╣   ║   ║\n   ║   ║           ║   maze  guy   ║       ║   ║   ║\n   ║   ║   ════════╣   ╔═══════╗   ║   ║   ║   ║   ║\n   ║   ║           ║   ║   ◗   ║   ║   ║   ║       ║\n   ║   ╚═══╗   ╔═══╝   ║   ════╣   ╠═══╝   ║   ════╣\n   ║       ║   ║       ║       ║   ║       ║       ║\n   ╠═══╗   ╚═══╣   ╔═══╩════   ║   ║   ════╩═══╗   ║\n   ║   ║       ║   ║           ║   ║           ║   ║\n   ║   ╚════   ║   ║   ║   ════╝   ╚════════   ║   ║\n   ║               ║   ║                       ║    \n   ╚═══════════════╩═══╩═══════════════════════╩════\n "

const SCREEN_MENU = 0
const SCREEN_GENERATOR = 1
const SCREEN_PLAY = 2

// sprite for each direction: n, e, s, w
var PLAYER_SPRITES = []string{"⯊", "◗", "⯋", "◖"}

type model struct {
	didSetInitalSize bool
	width            int
	height           int
	gmodel           *genmodel
	selectedScreen   int

	// generator screen properties
	grid               [][]cell
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
		switch m.selectedScreen {
		case SCREEN_MENU:
			cmd = tea.Batch(cmd, m.handleMenuKeyPress(msg))
		case SCREEN_GENERATOR:
			cmd = tea.Batch(cmd, m.handleGeneratorKeyPress(msg))
		case SCREEN_PLAY:
			cmd = tea.Batch(cmd, m.handlePlayKeyPress(msg))
		}
	case tea.WindowSizeMsg:
		m.setWindowSize(msg.Width, msg.Height)
	case generatorTickMsg:
		if m.selectedScreen == SCREEN_GENERATOR && !m.isPaused && !m.gmodel.isFinished {
			m.grid = m.gmodel.generateMaze(m.iterationsPerFrame)
			cmd = tea.Batch(cmd, tickCmd())
		}
	}

	return m, cmd
}

func (m *model) handleMenuKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd

	switch msg.String() {
	case "esc":
		cmd = tea.Quit
	case "enter":
		m.selectedScreen = SCREEN_GENERATOR
	}

	return cmd
}

func (m *model) handleGeneratorKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "esc":
		m.gmodel = &genmodel{}
		m.grid = [][]cell{}
		m.isPaused = false
		m.selectedScreen = SCREEN_MENU
	case "enter":
		m.gmodel = initializeMazeGen(m.getMazePixelSize())
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

	return cmd
}

func (m *model) handlePlayKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "esc":
		m.selectedScreen = SCREEN_MENU
	}
	return cmd
}

func (m *model) setWindowSize(w int, h int) {
	if !m.didSetInitalSize {
		m.didSetInitalSize = true
	}
	m.height = h
	m.width = w
}

func (m *model) getMazePixelSize() (int, int) {
	_, controlBarHeight := lipgloss.Size(m.generatorControlView())
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

	var v tea.View
	switch m.selectedScreen {
	case SCREEN_MENU:
		v = m.menuView()
	case SCREEN_GENERATOR:
		v = m.generatorView()
	case SCREEN_PLAY:
		v = m.playView()
	}

	v.AltScreen = true
	return v
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Nanosecond*1, func(t time.Time) tea.Msg {
		return generatorTickMsg(t)
	})
}

func (m *model) generatorView() tea.View {
	cv := m.generatorControlView()

	_, cvHeight := lipgloss.Size(cv)
	var content string
	if len(m.grid) == 0 {
		content = lipgloss.NewStyle().Height(m.height - cvHeight).AlignVertical(lipgloss.Center).Render("Press enter to generate")
	} else {
		content = m.generatorMazeView()
	}

	viewString := lipgloss.JoinVertical(
		lipgloss.Center,
		content,
		m.generatorControlView(),
	)

	return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, viewString))
}

func (m *model) playView() tea.View {
	return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "Play Screen"))
}

func (m *model) menuView() tea.View {
	return tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, TITLE_STRING))
}

func (m *model) generatorMazeView() string {
	width, height := m.getMazePixelSize()
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#6cacd4"))
	return style.Width(width).MaxWidth(width).Height(height).MaxHeight(height).AlignVertical(lipgloss.Center).AlignHorizontal(lipgloss.Center).Render(render(m.grid))
}

func (m *model) generatorControlView() string {
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
