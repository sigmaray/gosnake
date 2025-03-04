package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/rthornton128/goncurses"
	"github.com/sigmaray/gosnake/shared"
)

// EnvConfig holds configuration parsed from environment variables
type EnvConfig struct {
	UseTimer bool
	Size     int
	Timeout  int
}

// ParseEnv parses environment variables and returns an EnvConfig
func ParseEnv() EnvConfig {
	useTimer := true
	if val := os.Getenv("TIMER"); val != "" {
		val = strings.ToLower(val)
		if val == "0" || val == "false" || val == "off" {
			useTimer = false
		}
	}

	size := 5
	if val := os.Getenv("SIZE"); val != "" {
		if parsedSize, err := strconv.Atoi(val); err == nil && parsedSize > 1 {
			size = parsedSize
		}
	}

	timeout := 300
	if val := os.Getenv("TIMEOUT"); val != "" {
		if parsedTimeout, err := strconv.Atoi(val); err == nil && parsedTimeout > 1 {
			timeout = parsedTimeout
		}
	}

	return EnvConfig{
		UseTimer: useTimer,
		Size:     size,
		Timeout:  timeout,
	}
}

func main() {
	config := ParseEnv()

	state := shared.NewState(config.UseTimer, config.Size)

	// Initialize ncurses
	stdscr, err := goncurses.Init()
	if err != nil {
		log.Fatal("Init:", err)
	}
	defer goncurses.End()

	// Clean up on interrupt to avoid leaving terminal in a bad state
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		// Wait for the interrupt signal
		<-interrupt

		// Cleanup and exit
		goncurses.End()
		os.Exit(0)
	}()

	// Enable keypad for arrow keys
	stdscr.Keypad(true)

	// Don't echo input
	goncurses.Echo(false)

	// Hide cursor
	goncurses.Cursor(0)

	// Set timeout if timer is enabled
	if config.UseTimer {
		stdscr.Timeout(config.Timeout)
	}

	// Print instructions
	stdscr.MovePrint(0, 0, "Move snake with arrow keys. Press 'q' to quit. Press 'r' to restart")

	// Main loop
	for {
		// Render game board
		stdscr.MovePrint(4, 0, state.ToString())

		// Print game status
		status := "Game in progress"
		if state.DidWin {
			status = "You won         "
		}
		stdscr.MovePrint(config.Size+6, 0, status)

		// Print debug output
		stdscr.MovePrint(config.Size+8, 0, fmt.Sprintf("Debug output: %+v", state))

		// Handle input
		input := stdscr.GetChar()

		switch input {
		case goncurses.KEY_UP, 'w', 'W':
			state.OnKeyPress(shared.KeyUp)
		case goncurses.KEY_DOWN, 's', 'S':
			state.OnKeyPress(shared.KeyDown)
		case goncurses.KEY_LEFT, 'a', 'A':
			state.OnKeyPress(shared.KeyLeft)
		case goncurses.KEY_RIGHT, 'd', 'D':
			state.OnKeyPress(shared.KeyRight)
		case 'q', 'Q':
			return // Exit on 'q'
		case 'r', 'R':
			state.OnKeyPress(shared.KeyReset)
		case 0:
			if config.UseTimer {
				state.MoveSnake()
			}
		}
	}
}
