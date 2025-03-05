package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"image/color"
	imageColor "image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/sigmaray/gosnake/shared"
	"golang.org/x/image/font/basicfont"
)

const cellSize = 40

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

type Game struct {
	state  *shared.State
	config EnvConfig
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.state.OnKeyPress(shared.KeyUp)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.state.OnKeyPress(shared.KeyDown)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.state.OnKeyPress(shared.KeyLeft)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.state.OnKeyPress(shared.KeyRight)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		os.Exit(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.state.OnKeyPress(shared.KeyReset)
	}

	if g.config.UseTimer && ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.state.MoveSnake()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	board := g.state.To2DArray()
	// Draw the game board
	for y := 0; y < g.config.Size; y++ {
		for x := 0; x < g.config.Size; x++ {
			cell := board[y][x]
			color := imageColor.RGBA{R: 255, G: 255, B: 255, A: 255}
			if cell == shared.CELL_SNAKE_HEAD {
				color = imageColor.RGBA{R: 0, G: 255, B: 0, A: 255}
			} else if cell == shared.CELL_SNAKE_TAIL {
				color = imageColor.RGBA{R: 173, G: 255, B: 47, A: 255}
			} else if cell == shared.CELL_FOOD {
				color = imageColor.RGBA{R: 255, G: 0, B: 0, A: 255}
			}
			ebitenutil.DrawRect(screen, float64(x*cellSize), float64(y*cellSize), cellSize, cellSize, color)
		}
	}

	// Draw game status
	status := "Game in progress"
	if g.state.DidWin {
		status = "You won"
	}
	text.Draw(screen, status, basicfont.Face7x13, 10, g.config.Size*cellSize+20, color.White)

	// Draw debug output
	// stateJSON, _ := json.MarshalIndent(g.state, "", "  ")
	// debugText := fmt.Sprintf("Debug output: %+v", string(stateJSON))
	// text.Draw(screen, debugText, basicfont.Face7x13, 10, g.config.Size*cellSize+40, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.calcWidthHeight()
}

func (g *Game) calcWidthHeight() (int, int) {
	return g.config.Size*cellSize + 100, g.config.Size*cellSize + 100
}

func main() {
	config := ParseEnv()

	state := shared.NewState(config.UseTimer, config.Size)

	game := &Game{
		state:  state,
		config: config,
	}

	screenWidth, screenHeight := game.calcWidthHeight()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake Game")

	if game.config.UseTimer {
		go func() {
			for range time.Tick(time.Millisecond * time.Duration(game.config.Timeout)) {
				game.state.MoveSnake()
			}
		}()
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
