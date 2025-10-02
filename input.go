package main

import "github.com/gdamore/tcell/v2"

// HandleInput processes keyboard input and returns true if game should quit
func HandleInput(game *Game, event *tcell.EventKey) bool {
	// Handle quit keys
	if event.Key() == tcell.KeyEscape || event.Rune() == 'q' || event.Rune() == 'Q' {
		return true
	}
	
	// Handle pause
	if event.Rune() == 'p' || event.Rune() == 'P' {
		game.TogglePause()
		return false
	}

	// Handle hard drop
	if event.Rune() == ' ' {
		game.HardDrop()
		return false
	}
	
	// Handle game over state - allow restart with 'R'
	if game.State == StateGameOver {
		if event.Rune() == 'r' || event.Rune() == 'R' {
			game.Reset()
		}
		return false
	}
	
	// Don't process game controls if paused or game over
	if game.State != StatePlaying {
		return false
	}
	
	// Handle movement and rotation
	switch event.Key() {
	case tcell.KeyLeft:
		game.MoveLeft()
	case tcell.KeyRight:
		game.MoveRight()
	case tcell.KeyDown:
		game.MoveDown()
	case tcell.KeyUp:
		game.RotatePiece()
	}
	
	return false
}
