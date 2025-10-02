package game

import (
	"math/rand"
	"time"
)

const (
	// Board dimensions
	BoardWidth  = 10
	BoardHeight = 20
	
	// Game speed (milliseconds between automatic drops)
	InitialDropInterval = 1000
	MinDropInterval     = 100
	SpeedIncrement      = 50
)

// Piece represents a Tetris piece
type Piece struct {
	Shape    [][]int
	X, Y     int
	Color    int
	Rotation int
}

// Game represents the game state
type Game struct {
	board           [][]int
	currentPiece    *Piece
	score           int
	level           int
	lines           int
	paused          bool
	gameOver        bool
	lastDropTime    time.Time
	dropInterval    time.Duration
	rand            *rand.Rand
}

// Tetromino shapes (I, O, T, S, Z, J, L)
var shapes = [][][]int{
	// I piece
	{
		{1, 1, 1, 1},
	},
	// O piece
	{
		{2, 2},
		{2, 2},
	},
	// T piece
	{
		{0, 3, 0},
		{3, 3, 3},
	},
	// S piece
	{
		{0, 4, 4},
		{4, 4, 0},
	},
	// Z piece
	{
		{5, 5, 0},
		{0, 5, 5},
	},
	// J piece
	{
		{6, 0, 0},
		{6, 6, 6},
	},
	// L piece
	{
		{0, 0, 7},
		{7, 7, 7},
	},
}

// New creates a new game instance
func New() *Game {
	g := &Game{
		board:        make([][]int, BoardHeight),
		score:        0,
		level:        1,
		lines:        0,
		paused:       false,
		gameOver:     false,
		lastDropTime: time.Now(),
		dropInterval: GetGravitySpeed(1),
		rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	
	// Initialize board
	for i := range g.board {
		g.board[i] = make([]int, BoardWidth)
	}
	
	// Spawn first piece
	g.spawnPiece()
	
	return g
}

// GetBoard returns the current board state with the current piece merged
func (g *Game) GetBoard() [][]int {
	// Create a copy of the board
	result := make([][]int, BoardHeight)
	for i := range result {
		result[i] = make([]int, BoardWidth)
		copy(result[i], g.board[i])
	}
	
	// Overlay current piece if it exists
	if g.currentPiece != nil {
		for y, row := range g.currentPiece.Shape {
			for x, val := range row {
				if val > 0 {
					boardY := g.currentPiece.Y + y
					boardX := g.currentPiece.X + x
					if boardY >= 0 && boardY < BoardHeight && boardX >= 0 && boardX < BoardWidth {
						result[boardY][boardX] = val
					}
				}
			}
		}
	}
	
	return result
}

// GetScore returns the current score
func (g *Game) GetScore() int {
	return g.score
}

// GetLevel returns the current level
func (g *Game) GetLevel() int {
	return g.level
}

// GetLines returns the number of lines cleared
func (g *Game) GetLines() int {
	return g.lines
}

// IsPaused returns whether the game is paused
func (g *Game) IsPaused() bool {
	return g.paused
}

// IsGameOver returns whether the game is over
func (g *Game) IsGameOver() bool {
	return g.gameOver
}

// TogglePause toggles the pause state
func (g *Game) TogglePause() {
	if !g.gameOver {
		g.paused = !g.paused
	}
}

// Reset resets the game to initial state
func (g *Game) Reset() {
	// Clear board
	for i := range g.board {
		for j := range g.board[i] {
			g.board[i][j] = 0
		}
	}
	
	g.score = 0
	g.level = 1
	g.lines = 0
	g.paused = false
	g.gameOver = false
	g.lastDropTime = time.Now()
	g.dropInterval = GetGravitySpeed(g.level)
	g.spawnPiece()
}

// CalculateGravityInterval returns the drop interval for a given level
// Deprecated: Use GetGravitySpeed instead
func CalculateGravityInterval(level int) time.Duration {
	return GetGravitySpeed(level)
}

// ApplyGravity moves the current piece down one row (called by gravity ticker)
func (g *Game) ApplyGravity() {
	if g.paused || g.gameOver {
		return
	}
	g.MoveDown()
}

// Update updates the game state (called each frame)
func (g *Game) Update() {
	// Frame-based updates can go here if needed
	// Gravity is now handled by the gravity ticker calling ApplyGravity()
}

// MoveLeft moves the current piece left if possible
func (g *Game) MoveLeft() bool {
	if g.paused || g.gameOver || g.currentPiece == nil {
		return false
	}
	
	g.currentPiece.X--
	if g.hasCollision() {
		g.currentPiece.X++
		return false
	}
	return true
}

// MoveRight moves the current piece right if possible
func (g *Game) MoveRight() bool {
	if g.paused || g.gameOver || g.currentPiece == nil {
		return false
	}
	
	g.currentPiece.X++
	if g.hasCollision() {
		g.currentPiece.X--
		return false
	}
	return true
}

// MoveDown moves the current piece down if possible
func (g *Game) MoveDown() bool {
	if g.paused || g.gameOver || g.currentPiece == nil {
		return false
	}
	
	g.currentPiece.Y++
	if g.hasCollision() {
		g.currentPiece.Y--
		g.lockPiece()
		return false
	}
	g.score++ // Small score for soft drop
	return true
}

// HardDrop instantly drops the piece to the bottom
func (g *Game) HardDrop() {
	if g.paused || g.gameOver || g.currentPiece == nil {
		return
	}
	
	dropDistance := 0
	for g.MoveDown() {
		dropDistance++
	}
	g.score += dropDistance * 2 // Double score for hard drop
}

// Rotate rotates the current piece clockwise
func (g *Game) Rotate() bool {
	if g.paused || g.gameOver || g.currentPiece == nil {
		return false
	}
	
	// Save original shape
	originalShape := g.currentPiece.Shape
	
	// Rotate clockwise (transpose and reverse rows)
	height := len(g.currentPiece.Shape)
	width := len(g.currentPiece.Shape[0])
	rotated := make([][]int, width)
	for i := range rotated {
		rotated[i] = make([]int, height)
	}
	
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rotated[x][height-1-y] = g.currentPiece.Shape[y][x]
		}
	}
	
	g.currentPiece.Shape = rotated
	
	// Check for collision and try wall kicks
	if g.hasCollision() {
		// Try simple wall kicks
		kicks := []struct{ dx, dy int }{
			{-1, 0}, {1, 0}, {0, -1}, {-2, 0}, {2, 0},
		}
		
		kicked := false
		for _, kick := range kicks {
			g.currentPiece.X += kick.dx
			g.currentPiece.Y += kick.dy
			if !g.hasCollision() {
				kicked = true
				break
			}
			g.currentPiece.X -= kick.dx
			g.currentPiece.Y -= kick.dy
		}
		
		if !kicked {
			// Restore original shape
			g.currentPiece.Shape = originalShape
			return false
		}
	}
	
	return true
}

// hasCollision checks if the current piece collides with the board or boundaries
func (g *Game) hasCollision() bool {
	if g.currentPiece == nil {
		return false
	}
	
	for y, row := range g.currentPiece.Shape {
		for x, val := range row {
			if val > 0 {
				boardY := g.currentPiece.Y + y
				boardX := g.currentPiece.X + x
				
				// Check boundaries
				if boardX < 0 || boardX >= BoardWidth || boardY >= BoardHeight {
					return true
				}
				
				// Check board collision (ignore if piece is above board)
				if boardY >= 0 && g.board[boardY][boardX] > 0 {
					return true
				}
			}
		}
	}
	
	return false
}

// lockPiece locks the current piece into the board
func (g *Game) lockPiece() {
	if g.currentPiece == nil {
		return
	}
	
	// Copy piece to board
	for y, row := range g.currentPiece.Shape {
		for x, val := range row {
			if val > 0 {
				boardY := g.currentPiece.Y + y
				boardX := g.currentPiece.X + x
				if boardY >= 0 && boardY < BoardHeight && boardX >= 0 && boardX < BoardWidth {
					g.board[boardY][boardX] = val
				}
			}
		}
	}
	
	// Clear completed lines
	g.clearLines()
	
	// Spawn new piece
	g.spawnPiece()
	
	// Check for game over
	if g.hasCollision() {
		g.gameOver = true
	}
}

// clearLines clears completed lines and updates score
func (g *Game) clearLines() {
	// Check for completed lines
	completedRows := CheckLines(g.board)
	
	if len(completedRows) == 0 {
		return
	}
	
	linesCleared := len(completedRows)
	
	// Clear the lines from the board
	ClearLines(g.board, completedRows)
	
	// Update game state
	g.updateGameState(linesCleared)
}

// updateGameState updates score, lines counter, and level based on lines cleared
func (g *Game) updateGameState(linesCleared int) {
	if linesCleared <= 0 {
		return
	}
	
	// Update total lines cleared
	g.lines += linesCleared
	
	// Calculate and add score
	g.score += CalculateScore(linesCleared, g.level)
	
	// Check for level up (every 10 lines)
	newLevel := g.lines/10 + 1
	if newLevel > g.level {
		g.level = newLevel
		// Update drop speed based on new level
		g.dropInterval = GetGravitySpeed(g.level)
	}
}

// spawnPiece spawns a new random piece at the top
func (g *Game) spawnPiece() {
	shapeIndex := g.rand.Intn(len(shapes))
	shape := shapes[shapeIndex]
	
	// Deep copy the shape
	newShape := make([][]int, len(shape))
	for i := range shape {
		newShape[i] = make([]int, len(shape[i]))
		copy(newShape[i], shape[i])
	}
	
	g.currentPiece = &Piece{
		Shape:    newShape,
		X:        BoardWidth/2 - len(newShape[0])/2,
		Y:        0,
		Color:    shapeIndex + 1,
		Rotation: 0,
	}
}
