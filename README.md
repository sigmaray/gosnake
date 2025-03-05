# Go Snake Game

This is a simple snake game implemented in Go using the `goncurses` library for terminal-based graphics.

## Features

* Move the snake using arrow keys
* Automatic movement with a configurable timeout
* Random food generation
* You win when snake fills all game board cells
* After you win, you can restart the game by pressing 'r' key
* The snake can intersect with itself (this was implemented intentionally)
* The snake can go beyond the playing field (this was implemented intentionally)

## Installation

1. Clone the repository
2. Install deps: `apt-get update && apt-get install ncurses-dev`
2. Run the game: `go run cmd/main.go`

## Docker

You can also run the game using Docker. First, build the Docker image:
```
docker build -t gosnake .
```

Then, run the Docker container:
```
docker run -it gosnake
```

Or run with command:
```
docker run -it "$(docker build -q .)"
```

## Usage

- Use the arrow keys to change the direction of the snake.
- The snake will move automatically if `TIMER` is set to true.

# Environment Variables

You can configure the game using the following environment variables:

* TIMER (enabled by default) - setting it to "0", "false", or "off" will disable the timer; other values will enable the timer
* SIZE (default: 5) - size of the game board
* TIMEOUT (default: 300) - timeout for the timer in milliseconds
Examples:
```
TIMER=off SIZE=10 go run cmd/main.go
```
```
TIMEOUT=500 go run cmd/main.go
```

# Web version

There is web version with all functionality except timer (snake is moved only after key press)

How to run it:
```
go run web/main.go
```
```
SIZE=10 go run web/main.go
```
