package main

import (
	"log"
	"strconv"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type generatorTickMsg time.Time
type menuTickMsg time.Time
type playTickMsg time.Time

func main() {
	p := tea.NewProgram(&model{iterationsPerFrame: 1})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

const TITLE_STRING_A = "╔═══════╦═══════════╦═══════════╦═══════════════╗\n║       ║           ║           ║               ║\n╠════   ║   ║   ════╣   ╔═══╗   ║   ║   ╔═══╗   ║\n║       ║   ║       ║   ║   ║       ║   ║   ║   ║\n║   ╔═══╣   ╚═══╗   ║   ║   ╚═══╦═══╝   ║   ║   ║\n║   ║   ║       ║               ║       ║   ║   ║\n║   ║   ╚════   ╠═══════════════╣   ════╣   ║   ║\n║   ║           ║   "
const TITLE_STRING_B = "maze  guy"
const TITLE_STRING_C = "   ║       ║   ║   ║\n║   ║   ════════╣   ╔═══════╗   ║   ║   ║   ║   ║\n║   ║           ║   ║"
const TITLE_STRING_D = "║   ║   ║   ║       ║\n║   ╚═══╗   ╔═══╝   ║   ════╣   ╠═══╝   ║   ════╣\n║       ║   ║       ║       ║   ║       ║       ║\n╠═══╗   ╚═══╣   ╔═══╩════   ║   ║   ════╩═══╗   ║\n║   ║       ║   ║           ║   ║           ║   ║\n║   ╚════   ║   ║   ║   ════╝   ╚════════   ║   ║\n║               ║   ║                       ║   ║\n╚═══════════════╩═══╩═══════════════════════╩═══╝\n"
const SCREEN_MENU = 0
const SCREEN_GENERATOR = 1
const SCREEN_PLAY = 2

// sprite for each direction: n, e, s, w
// var PLAYER_SPRITES = []rune{'⯊', '◗', '⯋', '◖'}
var PLAYER_SPRITES = []rune{'🠵', '🠶', '🠷', '🠴'}
var KEY_DIRECTION_MAP = map[string]int{
	"up":    0,
	"right": 1,
	"down":  2,
	"left":  3,
}
var KEY_WASD_MAP = map[string]int{
	"w": 0,
	"d": 1,
	"s": 2,
	"a": 3,
}
var DIRECTION_TO_SPRITE_MAP = map[int]rune{
	0: PLAYER_SPRITES[0],
	1: PLAYER_SPRITES[1],
	2: PLAYER_SPRITES[2],
	3: PLAYER_SPRITES[3],
}
var SPRITE_TO_DIRECTION_MAP = map[rune]int{
	'🠵': 0,
	'🠶': 1,
	'🠷': 2,
	'🠴': 3,
}

var MENU_SPRITES = []string{
	" 🠶     ",
	"  🠶    ",
	"   🠶   ",
	"    🠶  ",
	"     🠶 ",
	"     🠴 ",
	"    🠴  ",
	"   🠴   ",
	"  🠴    ",
	" 🠴     ",
}

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
	generatorTickDelay int

	// menu screen properties
	menuSpriteIndex int

	// play screen
	pmodel playmodel
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
		var old = m.didSetInitalSize
		m.setWindowSize(msg.Width, msg.Height)
		if old == false && m.didSetInitalSize {
			// start menu animation
			cmd = tea.Batch(cmd, menuTickCmd())
		}
	case generatorTickMsg:
		if m.selectedScreen == SCREEN_GENERATOR && !m.isPaused && !m.gmodel.isFinished {
			m.grid = m.gmodel.generateMaze(m.iterationsPerFrame)
			cmd = tea.Batch(cmd, m.generatorTickCmd())
		}
	case menuTickMsg:
		if m.selectedScreen == SCREEN_MENU {
			m.menuSpriteIndex++
			cmd = tea.Batch(cmd, menuTickCmd())
		}
	case playTickMsg:
		if m.selectedScreen == SCREEN_PLAY {
			if m.pmodel.isFinished {
				m.pmodel = initializePlay(m.getMazePixelSize())
			} else {
				m.pmodel.tickPlay()
				cmd = tea.Batch(cmd, playTickCmd())
			}
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
		m.selectedScreen = SCREEN_PLAY
		m.pmodel = initializePlay(m.getMazePixelSize())
	case "g":
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
			m.generatorTickCmd(),
		)
	case "space":
		if len(m.grid) > 0 {
			m.isPaused = !m.isPaused
			if !m.isPaused {
				cmd = tea.Batch(
					cmd,
					m.generatorTickCmd(),
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
	case "down":
		switch m.generatorTickDelay {
		case 0, 1:
			m.generatorTickDelay = 0
		case 10:
			m.generatorTickDelay = 1
		default:
			if m.generatorTickDelay > 50 {
				m.generatorTickDelay -= 50
			} else {
				m.generatorTickDelay -= 10
			}
		}
	case "up":
		switch m.generatorTickDelay {
		case 0:
			m.generatorTickDelay = 1
		case 1:
			m.generatorTickDelay = 10
		default:
			if m.generatorTickDelay >= 50 {
				m.generatorTickDelay += 50
			} else {
				m.generatorTickDelay += 10
			}
		}
	}

	return cmd
}

func (m *model) handlePlayKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd
	var didManeuver bool
	switch msg.String() {
	case "esc":
		m.selectedScreen = SCREEN_MENU
	case "w", "a", "s", "d":
		m.pmodel.direction = KEY_WASD_MAP[msg.String()]
		didManeuver = true
	case "up", "down", "left", "right":
		m.pmodel.direction = KEY_DIRECTION_MAP[msg.String()]
		didManeuver = true
	}

	if !m.pmodel.isStarted && didManeuver {
		m.pmodel.isStarted = true
		cmd = tea.Batch(
			cmd,
			playTickCmd(),
		)
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
		v = m.pmodel.playView(m.width, m.height)
	}

	v.AltScreen = true
	return v
}

func (m *model) generatorTickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*time.Duration(m.generatorTickDelay), func(t time.Time) tea.Msg {
		return generatorTickMsg(t)
	})
}

func menuTickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return menuTickMsg(t)
	})
}

func playTickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return playTickMsg(t)
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

func (m *model) menuView() tea.View {
	greenText := lipgloss.NewStyle().Foreground(lipgloss.Color("#1ab753"))
	blueText := lipgloss.NewStyle().Foreground(lipgloss.Color("#6cacd4"))
	sprite := MENU_SPRITES[m.menuSpriteIndex%len(MENU_SPRITES)]
	faintStyle := lipgloss.NewStyle().Faint(true)
	titleString := TITLE_STRING_A +
		greenText.Render(TITLE_STRING_B) +
		TITLE_STRING_C +
		greenText.Render(sprite) +
		TITLE_STRING_D

	buttons := blueText.Render("<enter>") +
		faintStyle.Render(" to play") +
		" / " +
		blueText.Render("<g>") +
		faintStyle.Render(" for maze generator")

	vert := lipgloss.JoinVertical(
		lipgloss.Center,
		titleString,
		buttons,
	)

	return tea.NewView(lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		vert,
	))
}

func (m *model) generatorMazeView() string {
	width, height := m.getMazePixelSize()
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#6cacd4"))
	maze, _ := render(m.grid)
	return style.Width(width).MaxWidth(width).Height(height).MaxHeight(height).AlignVertical(lipgloss.Center).AlignHorizontal(lipgloss.Center).Render(maze)
}

func (m *model) generatorControlView() string {
	cmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6cacd4"))
	typeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#1ab753"))
	style := lipgloss.NewStyle().Faint(true)
	str := "iterations per tick: " +
		typeStyle.Render(strconv.Itoa(m.iterationsPerFrame)) +
		style.Render(" | ") +
		"tick delay (ms): " +
		typeStyle.Render(strconv.Itoa(m.generatorTickDelay)) +
		"\n" +
		cmdStyle.Render("<left/right>") +
		style.Render(" adjusts iterations per tick") +
		"\n" +
		cmdStyle.Render("<up/down>") +
		style.Render(" adjusts tick delay") +
		"\n" +
		cmdStyle.Render("<enter>") +
		style.Render(" regenerates") +
		"\n" +
		cmdStyle.Render("<space>") +
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
