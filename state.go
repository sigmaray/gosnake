package main

import (
	"math/rand"
	"strings"
	"time"
)

const (
	CELL_EMPTY      = "-"
	CELL_SNAKE_HEAD = "o"
	CELL_SNAKE_TAIL = "*"
	CELL_FOOD       = "@"
)

// Point represents coordinates on the game board
type Point struct {
	X int
	Y int
}

// Direction represents the direction of the snake
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

// AllowedKeys represents the keys that can be pressed by the user
type AllowedKeys int

const (
	KeyUp AllowedKeys = iota
	KeyDown
	KeyLeft
	KeyRight
	KeyReset
)

// State represents the game state
type State struct {
	BoardSize int
	UseTimer  bool
	Snake     []Point
	Food      Point
	Direction Direction
	DidWin    bool
}

// NewState initializes a new game state
func NewState(useTimer bool, boardSize int) *State {
	state := &State{
		BoardSize: boardSize,
		UseTimer:  useTimer,
		Snake:     []Point{{X: 0, Y: 0}},
		Direction: Right,
		DidWin:    false,
	}
	state.Food = state.genRandomFoodPosition()
	return state
}

// To2DArray renders the game state into a 2D array that represents the game board
func (s *State) To2DArray() [][]string {
	board := make([][]string, s.BoardSize)
	for i := range board {
		board[i] = make([]string, s.BoardSize)
		for j := range board[i] {
			board[i][j] = CELL_EMPTY
		}
	}

	// Place the food on the board
	board[s.Food.Y][s.Food.X] = CELL_FOOD

	// Place the snake on the board
	for _, segment := range s.Snake {
		board[segment.Y][segment.X] = CELL_SNAKE_TAIL
	}

	// Place the snake head
	head := s.Snake[0]
	board[head.Y][head.X] = CELL_SNAKE_HEAD

	return board
}

// ToString renders the game board into a multiline string that can be printed to the terminal
func (s *State) ToString() string {
	board := s.To2DArray()
	var boardStr strings.Builder

	for _, row := range board {
		for _, cell := range row {
			boardStr.WriteString(cell)
		}
		boardStr.WriteString("\n")
	}

	return boardStr.String()
}

// ChangeDirection changes the direction of the snake after user input if the direction is not forbidden
func (s *State) ChangeDirection(newDirection Direction) bool {
	if s.DidWin {
		return false
	}

	if (s.Direction == Up && newDirection == Down) ||
		(s.Direction == Down && newDirection == Up) ||
		(s.Direction == Left && newDirection == Right) ||
		(s.Direction == Right && newDirection == Left) {
		// These switches are forbidden
		return false
	}

	s.Direction = newDirection
	return true
}

// OnKeyPress handles key press events
func (s *State) OnKeyPress(key AllowedKeys) {
	var canMove bool
	switch key {
	case KeyUp:
		canMove = s.ChangeDirection(Up)
	case KeyDown:
		canMove = s.ChangeDirection(Down)
	case KeyLeft:
		canMove = s.ChangeDirection(Left)
	case KeyRight:
		canMove = s.ChangeDirection(Right)
	case KeyReset:
		*s = *NewState(s.UseTimer, s.BoardSize)
	}

	if canMove && !s.UseTimer {
		s.MoveSnake()
	}
}

// MoveSnake moves the snake per one step
func (s *State) MoveSnake() {
	if s.DidWin {
		return
	}

	head := s.Snake[0]
	newHead := head

	switch s.Direction {
	case Up:
		if head.Y > 0 {
			newHead.Y--
		} else {
			newHead.Y = s.BoardSize - 1
		}
	case Down:
		if head.Y < s.BoardSize-1 {
			newHead.Y++
		} else {
			newHead.Y = 0
		}
	case Left:
		if head.X > 0 {
			newHead.X--
		} else {
			newHead.X = s.BoardSize - 1
		}
	case Right:
		if head.X < s.BoardSize-1 {
			newHead.X++
		} else {
			newHead.X = 0
		}
	}

	s.Snake = append([]Point{newHead}, s.Snake...)

	if newHead.X != s.Food.X || newHead.Y != s.Food.Y {
		s.Snake = s.Snake[:len(s.Snake)-1]
	} else {
		if len(s.Snake) == s.BoardSize*s.BoardSize {
			s.DidWin = true
			return
		}
		s.Food = s.genRandomFoodPosition()
	}
}

// genRandomFoodPosition generates a random position for food on the game board that is not occupied by the snake
func (s *State) genRandomFoodPosition() Point {
	board := s.To2DArray()
	var freeCells []Point

	for y, row := range board {
		for x, cell := range row {
			if cell == CELL_EMPTY {
				freeCells = append(freeCells, Point{X: x, Y: y})
			}
		}
	}

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	randomIndex := rng.Intn(len(freeCells))
	return freeCells[randomIndex]
}
