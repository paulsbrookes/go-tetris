package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

const (
	// Screen offsets for positioning
	BoardOffsetX = 2
	BoardOffsetY = 2
	InfoOffsetX  = 28
	InfoOffsetY  = 2
)

// Renderer handles all drawing to the tcell screen
type Renderer struct {
	screen tcell.Screen
}

// NewRenderer creates a new renderer
func NewRenderer(screen tcell.Screen) *Renderer {
	return &Renderer{screen: screen}
}

// Clear clears the screen
func (r *Renderer) Clear() {
	r.screen.Clear()
}

// Show displays the rendered content
func (r *Renderer) Show() {
	r.screen.Show()
}

// drawText draws text at a specific position
func (r *Renderer) drawText(x, y int, text string, style tcell.Style) {
	for i, ch := range text {
		r.screen.SetContent(x+i, y, ch, nil, style)
	}
}

// drawBlock draws a single colored block
func (r *Renderer) drawBlock(x, y int, color tcell.Color) {
	style := tcell.StyleDefault.Background(color).Foreground(tcell.ColorBlack)
	r.screen.SetContent(x*2, y, ' ', nil, style)
	r.screen.SetContent(x*2+1, y, ' ', nil, style)
}

// drawBorder draws a border around the board
func (r *Renderer) drawBorder() {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	
	// Top border
	for x := 0; x < BoardWidth+2; x++ {
		r.drawText(BoardOffsetX+x*2, BoardOffsetY-1, "▀▀", style)
	}
	
	// Bottom border
	for x := 0; x < BoardWidth+2; x++ {
		r.drawText(BoardOffsetX+x*2, BoardOffsetY+BoardHeight, "▄▄", style)
	}
	
	// Side borders
	for y := 0; y < BoardHeight; y++ {
		r.drawText(BoardOffsetX-2, BoardOffsetY+y, "█", style)
		r.drawText(BoardOffsetX+BoardWidth*2, BoardOffsetY+y, "█", style)
	}
}

// RenderBoard renders the game board
func (r *Renderer) RenderBoard(board *Board) {
	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			if board.Grid[y][x].Occupied {
				r.drawBlock(BoardOffsetX/2+x, BoardOffsetY+y, board.Grid[y][x].Color)
			} else {
				// Draw empty cell with dark background
				r.drawBlock(BoardOffsetX/2+x, BoardOffsetY+y, tcell.ColorBlack)
			}
		}
	}
}

// RenderPiece renders the current piece
func (r *Renderer) RenderPiece(piece *Piece) {
	if piece == nil {
		return
	}

	blocks := piece.GetBlocks()
	color := piece.Color()

	for _, block := range blocks {
		if block.Y >= 0 && block.Y < BoardHeight {
			r.drawBlock(BoardOffsetX/2+block.X, BoardOffsetY+block.Y, color)
		}
	}
}

// RenderGhostPiece renders the ghost piece (shows where the piece will land)
func (r *Renderer) RenderGhostPiece(piece *Piece) {
	if piece == nil {
		return
	}

	blocks := piece.GetBlocks()

	// Use a dim gray style for ghost blocks
	style := tcell.StyleDefault.Background(tcell.ColorDarkGray).Foreground(tcell.ColorBlack)

	for _, block := range blocks {
		if block.Y >= 0 && block.Y < BoardHeight {
			x := BoardOffsetX + block.X*2
			y := BoardOffsetY + block.Y
			// Draw with a dimmed/outline appearance
			r.screen.SetContent(x, y, '▪', nil, style)
			r.screen.SetContent(x+1, y, '▪', nil, style)
		}
	}
}

// RenderNextPiece renders the next piece preview
func (r *Renderer) RenderNextPiece(piece *Piece) {
	if piece == nil {
		return
	}
	
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	r.drawText(InfoOffsetX, InfoOffsetY, "NEXT:", style)
	
	// Draw preview box
	shape := piece.Shape()
	color := piece.Color()
	
	previewX := InfoOffsetX
	previewY := InfoOffsetY + 2
	
	// Clear preview area
	for py := 0; py < 4; py++ {
		for px := 0; px < 4; px++ {
			r.drawBlock(previewX/2+px, previewY+py, tcell.ColorBlack)
		}
	}
	
	// Draw piece
	for y := 0; y < len(shape); y++ {
		for x := 0; x < len(shape[y]); x++ {
			if shape[y][x] {
				r.drawBlock(previewX/2+x, previewY+y, color)
			}
		}
	}
}

// RenderInfo renders score, level, and high score
func (r *Renderer) RenderInfo(game *Game) {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	
	// Score
	r.drawText(InfoOffsetX, InfoOffsetY+8, fmt.Sprintf("SCORE: %d", game.Score), style)
	
	// Level
	r.drawText(InfoOffsetX, InfoOffsetY+10, fmt.Sprintf("LEVEL: %d", game.Level), style)
	
	// Lines
	r.drawText(InfoOffsetX, InfoOffsetY+12, fmt.Sprintf("LINES: %d", game.LinesCleared), style)
	
	// High Score
	r.drawText(InfoOffsetX, InfoOffsetY+14, fmt.Sprintf("HIGH: %d", game.HighScore), style)
	
	// Controls
	style = tcell.StyleDefault.Foreground(tcell.ColorGray)
	r.drawText(InfoOffsetX, InfoOffsetY+18, "CONTROLS:", style)
	r.drawText(InfoOffsetX, InfoOffsetY+19, "← → : Move", style)
	r.drawText(InfoOffsetX, InfoOffsetY+20, "↑   : Rotate", style)
	r.drawText(InfoOffsetX, InfoOffsetY+21, "↓   : Soft Drop", style)
	r.drawText(InfoOffsetX, InfoOffsetY+22, "SPC : Hard Drop", style)
	r.drawText(InfoOffsetX, InfoOffsetY+23, "P   : Pause", style)
	r.drawText(InfoOffsetX, InfoOffsetY+24, "Q   : Quit", style)
}

// RenderPauseScreen renders the pause overlay
func (r *Renderer) RenderPauseScreen() {
	style := tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)
	
	centerX := BoardOffsetX + BoardWidth
	centerY := BoardOffsetY + BoardHeight/2
	
	r.drawText(centerX-4, centerY-1, "═════════", style)
	r.drawText(centerX-2, centerY, "PAUSED", style)
	r.drawText(centerX-4, centerY+1, "═════════", style)

	style = tcell.StyleDefault.Foreground(tcell.ColorWhite)
	r.drawText(centerX-5, centerY+3, "P to Resume", style)
}

// RenderGameOverScreen renders the game over screen
func (r *Renderer) RenderGameOverScreen(game *Game) {
	style := tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
	
	centerX := BoardOffsetX + BoardWidth
	centerY := BoardOffsetY + BoardHeight/2
	
	r.drawText(centerX-5, centerY-2, "═══════════", style)
	r.drawText(centerX-3, centerY-1, "GAME OVER", style)
	r.drawText(centerX-5, centerY, "═══════════", style)
	
	style = tcell.StyleDefault.Foreground(tcell.ColorWhite)
	r.drawText(centerX-6, centerY+2, fmt.Sprintf("Score: %d", game.Score), style)
	
	if game.Score == game.HighScore && game.Score > 0 {
		style = tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)
		r.drawText(centerX-7, centerY+3, "NEW HIGH SCORE!", style)
	}
	
	style = tcell.StyleDefault.Foreground(tcell.ColorGray)
	r.drawText(centerX-7, centerY+5, "R to Restart", style)
	r.drawText(centerX-7, centerY+6, "Q to Quit", style)
}

// Render renders the complete game state
func (r *Renderer) Render(game *Game) {
	r.Clear()

	// Draw border
	r.drawBorder()

	// Draw board
	r.RenderBoard(game.Board)

	// Draw ghost piece and current piece (if playing)
	if game.State == StatePlaying {
		ghostPiece := game.GetGhostPiece()
		r.RenderGhostPiece(ghostPiece)
		r.RenderPiece(game.CurrentPiece)
	}

	// Draw next piece and info
	r.RenderNextPiece(game.NextPiece)
	r.RenderInfo(game)

	// Draw overlays
	if game.State == StatePaused {
		r.RenderPauseScreen()
	} else if game.State == StateGameOver {
		r.RenderGameOverScreen(game)
	}

	r.Show()
}
