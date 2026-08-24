package main

const SPACE uint8 = 0
const WALL uint8 = 1
const START uint8 = 2
const END uint8 = 3

type playmodel struct {
	// 0:n, 1:e, 2:s, 3:w
	direction int
	playerLoc vec2
	mazeRunes []rune
	width     int
	height    int
}

func initializePlay(mazeWidth int, mazeHeight int) playmodel {
	mazeGen := initializeMazeGen(mazeWidth, mazeHeight)
	mazeCells := mazeGen.generateMaze(0)
	mazeRunes := []rune(render(mazeCells))

	// place player sprite at entrance
	for i := range len(mazeRunes) {
		if mazeRunes[i] == START_RUNE {
			mazeRunes[i] = rune(PLAYER_SPRITES[1])
			break
		}
	}

	return playmodel{
		direction: 1,
		mazeRunes: mazeRunes,
	}
}

func (m *playmodel) tickPlay() {
	// find and replace player rune
	for i := range len(m.mazeRunes) {
		if runesContainRune(PLAYER_SPRITES, m.mazeRunes[i]) {
			m.mazeRunes[i] = SPRITE_DIRECTION_MAP[m.direction]
			break
		}
	}
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
