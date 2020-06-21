package main

import (
	"fmt"
	"time"
	"github.com/kindlyfire/go-keylogger"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/go-vgo/robotgo"
)

const (
	MENU = 0
	RECORDING = 1
	PLAYING = 2
	EXIT = 3
)

type click struct {
	x int
	y int
	t int64
	button string
}

var state = MENU
var recording = []click{}
var startTime = time.Now()
var recordingIndex = 0
var leftDown = false
var rightDown = false

func main() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil { panic(err) }

	printMenu()

	go keyboardInput()

	for {
		if state == RECORDING {
			x, y, val := sdl.GetGlobalMouseState()

			if val == 1 && !leftDown {
				recording = append(recording, click{int(x), int(y), time.Now().Sub(startTime).Nanoseconds(), "left"})
				leftDown = true
			} else if val == 4 && !rightDown {
				recording = append(recording, click{int(x), int(y), time.Now().Sub(startTime).Nanoseconds(), "right"})
				rightDown = true
			}

			if val != 1 {
				leftDown = false
			}

			if val != 4 {
				rightDown = false
			}
		} else if state == PLAYING {
			delta := time.Now().Sub(startTime).Nanoseconds()

			r := recording[recordingIndex]

			if delta > r.t {
				robotgo.MovesClick(r.x, r.y, r.button, false)
				recordingIndex++
			}

			if recordingIndex >= len(recording) {
				recordingIndex = 0
				startTime = time.Now()
			}
		} else if state == EXIT {
			break
		}
	}
}

func keyboardInput() {
	kl := keylogger.NewKeylogger()

	for {

		key := kl.GetKey()
		if key.Empty {
			continue
		}

		char := key.Rune

		if char == 'e' {
			state = EXIT
			break
		}

		if state == MENU {
			menuInput(char)
		} else if state == RECORDING {
			recordingInput(char)
		} else if state == PLAYING {
			playingInput(char)
		}
	}
}

func menuInput(char rune) {
	if char == 115 { // s
		recording = []click{}
		startTime = time.Now()
		state = RECORDING
		printRecording()
	} else if char == 112 { // p
		state = PLAYING
		startTime = time.Now()
		recordingIndex = 0
		printPlaying()
	}
}

func recordingInput(char rune) {
	if char == 120 {
		state = MENU
		printMenu()
	}
}

func playingInput(char rune) {
	if char == 120 {
		state = MENU
		printMenu()
	}
}

func printMenu() {
	fmt.Println()
	fmt.Println("[s] Start recording")
	fmt.Println("[p] Play recording")
	fmt.Println("[e] Exit")
}

func printRecording() {
	fmt.Println()
	fmt.Println("Recording...")
	fmt.Println("[x] Stop recording")
	fmt.Println("[e] Exit")
}

func printPlaying() {
	fmt.Println()
	fmt.Println("Playing...")
	fmt.Println("[x] Stop playing")
	fmt.Println("[e] Exit")
}