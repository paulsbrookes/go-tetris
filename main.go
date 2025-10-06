package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
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
	var screen tcell.Screen = inputHandler.GetScreen()
	
	// Start input handling
	inputHandler.Start()
	
	// Ensure cleanup on exit with panic recovery
	defer func() {
		if r := recover(); r != nil {
			// Ensure terminal cleanup even on panic
			inputHandler.Stop()
			fmt.Fprintf(os.Stderr, "Panic recovered: %v\n", r)
			os.Exit(1)
		}
		inputHandler.Stop()
	}()
	
	// Create new game
	gameState := game.New()
	
	// Run the game loop
	runGameLoop(gameState, inputHandler, screen)
}

// runGameLoop implements the main game loop with frame and gravity tickers
func runGameLoop(gameState *game.Game, inputHandler *input.Handler, screen tcell.Screen) {
	// Initialize frame ticker for consistent 60 FPS
	frameTicker := time.NewTicker(FrameDuration)
	defer frameTicker.Stop()
	
	// Initialize gravity ticker based on starting level
	gravityInterval := game.CalculateGravityInterval(gameState.GetLevel())
	gravityTicker := time.NewTicker(gravityInterval)
	defer gravityTicker.Stop()
	
	// Track game state for ticker management
	lastLevel := gameState.GetLevel()
	wasPaused := gameState.IsPaused()
	running := true
	
	for running {
		select {
		case <-frameTicker.C:
			// Frame tick: update and render
			gameState.Update()
			renderer.Render(screen, gameState)
			
			// Check if level changed and adjust gravity speed
			currentLevel := gameState.GetLevel()
			if currentLevel != lastLevel {
				lastLevel = currentLevel
				gravityTicker.Stop()
				gravityInterval = game.CalculateGravityInterval(currentLevel)
				gravityTicker = time.NewTicker(gravityInterval)
				defer gravityTicker.Stop()
			}
			
			// Handle pause state changes
			isPaused := gameState.IsPaused()
			if isPaused != wasPaused {
				wasPaused = isPaused
				if isPaused {
					// Stop gravity ticker when paused
					gravityTicker.Stop()
				} else {
					// Resume gravity ticker when unpaused
					gravityInterval = game.CalculateGravityInterval(currentLevel)
					gravityTicker = time.NewTicker(gravityInterval)
					defer gravityTicker.Stop()
				}
			}
			
		case <-gravityTicker.C:
			// Gravity tick: apply automatic piece drop
			if !gameState.IsPaused() && !gameState.IsGameOver() {
				gameState.ApplyGravity()
			}
			
		case action := <-inputHandler.Actions():
			// Input event: handle user action
			running = handleAction(gameState, action)
			if !running {
				return
			}
		}
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
