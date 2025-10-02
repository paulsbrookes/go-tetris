package main

import (
	"fmt"
	"os"
	"time"
	
	"github.com/yourusername/go-tetris/game"
	"github.com/yourusername/go-tetris/input"
	"github.com/yourusername/go-tetris/renderer"
)

const (
	// Target frames per second
	TargetFPS = 60
	// Frame duration
	FrameDuration = time.Second / TargetFPS
)

func main() {
	// Create input handler
	inputHandler := input.NewHandler()
	
	// Initialize screen
	if err := inputHandler.InitScreen(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize screen: %v\n", err)
		os.Exit(1)
	}
	
	// Get screen reference
	screen := inputHandler.GetScreen()
	
	// Start input handling
	inputHandler.Start()
	
	// Ensure cleanup on exit
	defer inputHandler.Stop()
	
	// Create new game
	gameState := game.New()
	
	// Main game loop
	running := true
	lastFrameTime := time.Now()
	
	for running {
		frameStart := time.Now()
		
		// Process input events
		select {
		case action := <-inputHandler.Actions():
			running = handleAction(gameState, action)
		default:
			// No input to process
		}
		
		// Update game state
		gameState.Update()
		
		// Render current state
		renderer.Render(screen, gameState)
		
		// Frame rate limiting
		elapsed := time.Since(frameStart)
		if elapsed < FrameDuration {
			time.Sleep(FrameDuration - elapsed)
		}
		
		lastFrameTime = frameStart
	}
}

// handleAction processes a user action and returns whether to continue running
func handleAction(g *game.Game, action input.Action) bool {
	switch action {
	case input.ActionQuit:
		return false
	case input.ActionMoveLeft:
		g.MoveLeft()
	case input.ActionMoveRight:
		g.MoveRight()
	case input.ActionMoveDown:
		g.MoveDown()
	case input.ActionRotateClockwise:
		g.Rotate()
	case input.ActionRotateCounterClockwise:
		// Counter-clockwise rotation: rotate 3 times clockwise
		g.Rotate()
		g.Rotate()
		g.Rotate()
	case input.ActionHardDrop:
		g.HardDrop()
	case input.ActionPause:
		g.TogglePause()
	case input.ActionRestart:
		g.Reset()
	}
	return true
}
