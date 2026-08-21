package main

const SPACE uint8 = 0
const WALL uint8 = 1
const START uint8 = 2
const END uint8 = 3

type playmodel struct {
	// 0:n, 1:e, 2:s, 3:w
	direction int
	playerLoc vec2
	mazeView  string
}

// func initializePlay(mazeWidth int, mazeHeight int) playmodel {
// 	return playmodel{}
// }
