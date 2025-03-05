package main

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sigmaray/gosnake/shared"
)

// EnvConfig holds configuration parsed from environment variables
type EnvConfig struct {
	// UseTimer bool
	Size int
	// Timeout  int
}

// ParseEnv parses environment variables and returns an EnvConfig
func ParseEnv() EnvConfig {
	// useTimer := true
	// if val := os.Getenv("TIMER"); val != "" {
	// 	val = strings.ToLower(val)
	// 	if val == "0" || val == "false" || val == "off" {
	// 		useTimer = false
	// 	}
	// }

	size := 5
	if val := os.Getenv("SIZE"); val != "" {
		if parsedSize, err := strconv.Atoi(val); err == nil && parsedSize > 1 {
			size = parsedSize
		}
	}

	// timeout := 300
	// if val := os.Getenv("TIMEOUT"); val != "" {
	// 	if parsedTimeout, err := strconv.Atoi(val); err == nil && parsedTimeout > 1 {
	// 		timeout = parsedTimeout
	// 	}
	// }

	return EnvConfig{
		// UseTimer: useTimer,
		Size: size,
		// Timeout:  timeout,
	}
}

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Load HTML templates from the "templates" directory
	router.LoadHTMLGlob("web/templates/*")

	config := ParseEnv()
	state := shared.NewState(false, config.Size)

	// Define a route for the root URL
	router.GET("/", func(c *gin.Context) {
		// Get the value of the "key" parameter from the query string
		key := c.Query("key")

		// Handle the key parameter
		switch key {
		case "up":
			state.OnKeyPress(shared.KeyUp)
		case "down":
			state.OnKeyPress(shared.KeyDown)
		case "left":
			state.OnKeyPress(shared.KeyLeft)
		case "right":
			state.OnKeyPress(shared.KeyRight)
		case "r":
			state.OnKeyPress(shared.KeyReset)
		}

		stateJSON, _ := json.MarshalIndent(state, "", "  ")

		c.HTML(200, "index.html", gin.H{
			"title":     "Hello, World!",
			"board":     state.ToString(),
			"stateJSON": string(stateJSON),
		})
	})

	// Run the server on port 8080
	router.Run(":8080")
}
