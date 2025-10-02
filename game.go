package main

import (
	"math"
	"time"
)

// GameState represents the current state of the game
type GameState int

const (
	StatePlaying GameState = iota
	StatePaused
	StateGameOver
)

// Game represents the complete game state
type Game struct {
	Board         *Board
	CurrentPiece  *Piece
	NextPiece     *Piece
	State         GameState
	Score         int
	HighScore     int
	Level         int
	LinesCleared  int
	LastFallTime  time.Time
	FallInterval  time.Duration
}

// NewGame creates a new game
func NewGame() *Game {
	g := &Game{
		Board:        NewBoard(),
		State:        StatePlaying,
		HighScore:    LoadHighScore(),
		Level:        1,
		LastFallTime: time.Now(),
	}
	
	g.NextPiece = NewRandomPiece()
	g.SpawnPiece()
	g.UpdateFallInterval()
	
	return g
}

// SpawnPiece spawns the next piece
func (g *Game) SpawnPiece() bool {
	g.CurrentPiece = g.NextPiece
	g.NextPiece = NewRandomPiece()
	
	// Check if piece can be placed (game over if not)
	if !g.Board.CanPlacePiece(g.CurrentPiece) {
		g.State = StateGameOver
		if g.Score > g.HighScore {
			g.HighScore = g.Score
			SaveHighScore(g.HighScore)
		}
		return false
	}
	
	return true
}

// MoveLeft moves the current piece left if possible
func (g *Game) MoveLeft() bool {
	if g.State != StatePlaying || g.CurrentPiece == nil {
		return false
	}
	
	g.CurrentPiece.X--
	if !g.Board.CanPlacePiece(g.CurrentPiece) {
		g.CurrentPiece.X++
		return false
	}
	
	return true
}

// MoveRight moves the current piece right if possible
func (g *Game) MoveRight() bool {
	if g.State != StatePlaying || g.CurrentPiece == nil {
		return false
	}
	
	g.CurrentPiece.X++
	if !g.Board.CanPlacePiece(g.CurrentPiece) {
		g.CurrentPiece.X--
		return false
	}
	
	return true
}

// MoveDown moves the current piece down if possible, locks it if not
func (g *Game) MoveDown() bool {
	if g.State != StatePlaying || g.CurrentPiece == nil {
		return false
	}
	
	g.CurrentPiece.Y++
	if !g.Board.CanPlacePiece(g.CurrentPiece) {
		g.CurrentPiece.Y--
		g.LockPiece()
		return false
	}
	
	return true
}

// HardDrop drops the piece all the way down instantly
func (g *Game) HardDrop() {
	if g.State != StatePlaying || g.CurrentPiece == nil {
		return
	}
	
	for g.MoveDown() {
		// Keep moving down until it can't
	}
}

// RotatePiece rotates the current piece with wall kick attempts
func (g *Game) RotatePiece() bool {
	if g.State != StatePlaying || g.CurrentPiece == nil {
		return false
	}
	
	// Save original position
	origX := g.CurrentPiece.X
	origY := g.CurrentPiece.Y
	
	// Try rotation
	g.CurrentPiece.Rotate()
	
	// If it doesn't fit, try wall kicks
	if !g.Board.CanPlacePiece(g.CurrentPiece) {
		// Try shifting left
		g.CurrentPiece.X--
		if g.Board.CanPlacePiece(g.CurrentPiece) {
			return true
		}
		
		// Try shifting right
		g.CurrentPiece.X = origX + 1
		if g.Board.CanPlacePiece(g.CurrentPiece) {
			return true
		}
		
		// Try shifting right more (for I piece)
		g.CurrentPiece.X = origX + 2
		if g.Board.CanPlacePiece(g.CurrentPiece) {
			return true
		}
		
		// Try shifting left more (for I piece)
		g.CurrentPiece.X = origX - 2
		if g.Board.CanPlacePiece(g.CurrentPiece) {
			return true
		}
		
		// Rotation failed, revert
		g.CurrentPiece.X = origX
		g.CurrentPiece.Y = origY
		g.CurrentPiece.UnRotate()
		return false
	}
	
	return true
}

// LockPiece locks the current piece to the board
func (g *Game) LockPiece() {
	if g.CurrentPiece == nil {
		return
	}
	
	g.Board.LockPiece(g.CurrentPiece)
	
	// Clear full lines
	linesCleared := g.Board.ClearFullLines()
	if linesCleared > 0 {
		g.LinesCleared += linesCleared
		g.Score += CalculateScore(linesCleared, g.Level)
		
		// Level up every 10 lines
		newLevel := 1 + (g.LinesCleared / 10)
		if newLevel > g.Level {
			g.Level = newLevel
			g.UpdateFallInterval()
		}
	}
	
	// Spawn next piece
	g.SpawnPiece()
}

// UpdateFallInterval updates the fall speed based on level
func (g *Game) UpdateFallInterval() {
	// Classic Tetris formula: fall speed increases with level
	// Level 1: ~1 second, Level 10: ~0.2 seconds
	baseInterval := 1000.0 // milliseconds
	speedMultiplier := math.Pow(0.8, float64(g.Level-1))
	interval := baseInterval * speedMultiplier
	
	// Cap minimum at 100ms
	if interval < 100 {
		interval = 100
	}
	
	g.FallInterval = time.Duration(interval) * time.Millisecond
}

// Update updates the game state (called each frame)
func (g *Game) Update() {
	if g.State != StatePlaying {
		return
	}
	
	// Check if it's time for piece to fall
	now := time.Now()
	if now.Sub(g.LastFallTime) >= g.FallInterval {
		g.MoveDown()
		g.LastFallTime = now
	}
}

// TogglePause toggles the pause state
func (g *Game) TogglePause() {
	if g.State == StatePlaying {
		g.State = StatePaused
	} else if g.State == StatePaused {
		g.State = StatePlaying
		g.LastFallTime = time.Now() // Reset fall timer
	}
}

// Reset resets the game to initial state
func (g *Game) Reset() {
	g.Board.Clear()
	g.State = StatePlaying
	g.Score = 0
	g.Level = 1
	g.LinesCleared = 0
	g.NextPiece = NewRandomPiece()
	g.SpawnPiece()
	g.UpdateFallInterval()
	g.LastFallTime = time.Now()
}
