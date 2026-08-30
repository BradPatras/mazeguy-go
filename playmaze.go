package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const SPACE uint8 = 0
const WALL uint8 = 1
const START uint8 = 2
const END uint8 = 3
const SPACE_RUNE = rune(' ')

type playmodel struct {
	// 0:n, 1:e, 2:s, 3:w
	direction  int
	mazeRunes  []rune
	isFinished bool
	isStarted  bool
	pixelWidth int
}

func initializePlay(mazeWidth int, mazeHeight int) playmodel {
	mazeGen := initializeMazeGen(mazeWidth, mazeHeight)
	mazeCells := mazeGen.generateMaze(0)
	mazeString, pixelWidth := render(mazeCells)
	mazeRunes := []rune(mazeString)

	// place player sprite at entrance
	for i := range len(mazeRunes) {
		if mazeRunes[i] == START_RUNE {
			mazeRunes[i] = rune(PLAYER_SPRITES[1])
			break
		}
	}

	return playmodel{
		direction:  1,
		mazeRunes:  mazeRunes,
		pixelWidth: pixelWidth,
	}
}

func (m *playmodel) tickPlay() {
	// find the player rune current location
	var playerRuneIndex int
	for i := range len(m.mazeRunes) {
		if runesContainRune(PLAYER_SPRITES, m.mazeRunes[i]) {
			playerRuneIndex = i
			break
		}
	}

	// try to move in the requested direction
	// if it's not possible, try moving in the
	// player's current sprite direction
	if !m.tryMove(playerRuneIndex, m.direction) {
		m.tryMove(playerRuneIndex, SPRITE_TO_DIRECTION_MAP[m.mazeRunes[playerRuneIndex]])
	}
}

func (m *playmodel) tryMove(playerRuneIndex int, direction int) bool {
	playerRune := DIRECTION_TO_SPRITE_MAP[direction]
	switch direction {
	case 0: // up
		if m.mazeRunes[playerRuneIndex-(m.pixelWidth)] == SPACE_RUNE {
			m.mazeRunes[playerRuneIndex-(m.pixelWidth)] = playerRune
			m.mazeRunes[playerRuneIndex] = SPACE_RUNE
			return true
		}
	case 1: // right
		if m.mazeRunes[playerRuneIndex+1] == SPACE_RUNE && m.mazeRunes[playerRuneIndex+2] == SPACE_RUNE && m.mazeRunes[playerRuneIndex+3] == SPACE_RUNE {
			m.mazeRunes[playerRuneIndex+2] = playerRune
			m.mazeRunes[playerRuneIndex] = SPACE_RUNE
			return true
		} else if m.mazeRunes[playerRuneIndex+2] == END_RUNE {
			m.isFinished = true
			m.mazeRunes[playerRuneIndex+2] = playerRune
			m.mazeRunes[playerRuneIndex] = SPACE_RUNE
		}
	case 2: // down
		if m.mazeRunes[playerRuneIndex+(m.pixelWidth)] == SPACE_RUNE {
			m.mazeRunes[playerRuneIndex+(m.pixelWidth)] = playerRune
			m.mazeRunes[playerRuneIndex] = SPACE_RUNE
			return true
		}
	case 3: // left
		if m.mazeRunes[playerRuneIndex-1] == SPACE_RUNE && m.mazeRunes[playerRuneIndex-2] == SPACE_RUNE && m.mazeRunes[playerRuneIndex-3] == SPACE_RUNE {
			m.mazeRunes[playerRuneIndex-2] = playerRune
			m.mazeRunes[playerRuneIndex] = SPACE_RUNE
			return true
		}
	}

	return false
}

func (m *playmodel) playView(width int, height int) tea.View {
	blueText := lipgloss.NewStyle().Foreground(lipgloss.Color("#6cacd4"))
	faintStyle := lipgloss.NewStyle().Faint(true)
	buttons := blueText.Render("<arrow keys or wasd>") +
		faintStyle.Render(" to manuever")
	vert := lipgloss.JoinVertical(
		lipgloss.Center,
		m.styledMazeString(),
		buttons,
	)

	return tea.NewView(lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, vert))
}

func (m *playmodel) styledMazeString() string {
	var sb strings.Builder
	playerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#1ab753"))
	endStyle := lipgloss.NewStyle().Foreground(lipgloss.BrightRed)
	for i := range len(m.mazeRunes) {
		if runesContainRune(PLAYER_SPRITES, m.mazeRunes[i]) {
			sb.WriteString(playerStyle.Render(string(m.mazeRunes[i])))
		} else if m.mazeRunes[i] == END_RUNE {
			sb.WriteString(endStyle.Render(string(m.mazeRunes[i])))
		} else {
			sb.WriteRune(m.mazeRunes[i])
		}
	}

	return sb.String()
}

func replaceAtIndex(str string, replacement rune, index int) string {
	return str[:index] + string(replacement) + str[index+1:]
}

func locToStringIndex(width int, loc vec2) int {
	return ((width + 1) * loc.y) + loc.x
}

func runesContainRune(runes []rune, r rune) bool {
	for i := range runes {
		if runes[i] == r {
			return true
		}
	}

	return false
}
