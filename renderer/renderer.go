package renderer

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

// GameState represents the minimum interface needed for rendering
type GameState interface {
	GetBoard() [][]int
	GetScore() int
	GetLevel() int
	GetLines() int
	IsPaused() bool
	IsGameOver() bool
}

// Render draws the game state to the terminal screen
func Render(screen tcell.Screen, game GameState) {
	if screen == nil {
		return
	}
	
	screen.Clear()
	
	// Get board dimensions
	board := game.GetBoard()
	if len(board) == 0 {
		return
	}
	
	boardHeight := len(board)
	boardWidth := len(board[0])
	
	// Calculate offsets for centering
	screenWidth, screenHeight := screen.Size()
	offsetX := (screenWidth - boardWidth*2 - 4) / 2
	offsetY := (screenHeight - boardHeight - 2) / 2
	
	if offsetX < 0 {
		offsetX = 0
	}
	if offsetY < 0 {
		offsetY = 0
	}
	
	// Draw top border
	drawHorizontalLine(screen, offsetX, offsetY, boardWidth*2+2, tcell.StyleDefault)
	
	// Draw board with side borders
	for y := 0; y < boardHeight; y++ {
		// Left border
		screen.SetContent(offsetX, offsetY+1+y, '|', nil, tcell.StyleDefault)
		
		// Board cells
		for x := 0; x < boardWidth; x++ {
			cellValue := board[y][x]
			style := getCellStyle(cellValue)
			char := ' '
			if cellValue > 0 {
				char = '█'
			}
			
			// Draw cell twice for better aspect ratio
			screen.SetContent(offsetX+1+x*2, offsetY+1+y, char, nil, style)
			screen.SetContent(offsetX+2+x*2, offsetY+1+y, char, nil, style)
		}
		
		// Right border
		screen.SetContent(offsetX+boardWidth*2+1, offsetY+1+y, '|', nil, tcell.StyleDefault)
	}
	
	// Draw bottom border
	drawHorizontalLine(screen, offsetX, offsetY+boardHeight+1, boardWidth*2+2, tcell.StyleDefault)
	
	// Draw game info (score, level, lines) to the right of the board
	infoX := offsetX + boardWidth*2 + 4
	infoY := offsetY + 2
	
	drawText(screen, infoX, infoY, "SCORE:", tcell.StyleDefault.Bold(true))
	drawText(screen, infoX, infoY+1, fmt.Sprintf("%d", game.GetScore()), tcell.StyleDefault)
	
	drawText(screen, infoX, infoY+3, "LEVEL:", tcell.StyleDefault.Bold(true))
	drawText(screen, infoX, infoY+4, fmt.Sprintf("%d", game.GetLevel()), tcell.StyleDefault)
	
	drawText(screen, infoX, infoY+6, "LINES:", tcell.StyleDefault.Bold(true))
	drawText(screen, infoX, infoY+7, fmt.Sprintf("%d", game.GetLines()), tcell.StyleDefault)
	
	// Draw controls
	controlsY := offsetY + boardHeight - 5
	drawText(screen, infoX, controlsY, "CONTROLS:", tcell.StyleDefault.Bold(true))
	drawText(screen, infoX, controlsY+1, "←/→: Move", tcell.StyleDefault.Dim(true))
	drawText(screen, infoX, controlsY+2, "↓: Drop", tcell.StyleDefault.Dim(true))
	drawText(screen, infoX, controlsY+3, "↑: Rotate", tcell.StyleDefault.Dim(true))
	drawText(screen, infoX, controlsY+4, "P: Pause", tcell.StyleDefault.Dim(true))
	drawText(screen, infoX, controlsY+5, "Q: Quit", tcell.StyleDefault.Dim(true))
	
	// Draw pause overlay
	if game.IsPaused() {
		pauseText := "PAUSED"
		pauseX := offsetX + (boardWidth*2+2-len(pauseText))/2
		pauseY := offsetY + boardHeight/2
		
		// Draw background
		style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow).Bold(true)
		drawText(screen, pauseX, pauseY, pauseText, style)
	}
	
	// Draw game over overlay
	if game.IsGameOver() {
		gameOverText := "GAME OVER"
		restartText := "Press R to restart"
		
		gameOverX := offsetX + (boardWidth*2+2-len(gameOverText))/2
		restartX := offsetX + (boardWidth*2+2-len(restartText))/2
		centerY := offsetY + boardHeight/2
		
		style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed).Bold(true)
		drawText(screen, gameOverX, centerY, gameOverText, style)
		
		style = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
		drawText(screen, restartX, centerY+2, restartText, style)
	}
	
	screen.Show()
}

// drawHorizontalLine draws a horizontal line of characters
func drawHorizontalLine(screen tcell.Screen, x, y, width int, style tcell.Style) {
	for i := 0; i < width; i++ {
		char := '-'
		if i == 0 {
			char = '+'
		} else if i == width-1 {
			char = '+'
		}
		screen.SetContent(x+i, y, char, nil, style)
	}
}

// drawText draws a string at the specified position
func drawText(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	for i, ch := range text {
		screen.SetContent(x+i, y, ch, nil, style)
	}
}

// getCellStyle returns the style for a given cell value
func getCellStyle(cellValue int) tcell.Style {
	baseStyle := tcell.StyleDefault
	
	switch cellValue {
	case 0:
		return baseStyle.Background(tcell.ColorBlack).Foreground(tcell.ColorBlack)
	case 1:
		return baseStyle.Background(tcell.ColorCyan).Foreground(tcell.ColorCyan)
	case 2:
		return baseStyle.Background(tcell.ColorBlue).Foreground(tcell.ColorBlue)
	case 3:
		return baseStyle.Background(tcell.ColorOrange).Foreground(tcell.ColorOrange)
	case 4:
		return baseStyle.Background(tcell.ColorYellow).Foreground(tcell.ColorYellow)
	case 5:
		return baseStyle.Background(tcell.ColorGreen).Foreground(tcell.ColorGreen)
	case 6:
		return baseStyle.Background(tcell.ColorPurple).Foreground(tcell.ColorPurple)
	case 7:
		return baseStyle.Background(tcell.ColorRed).Foreground(tcell.ColorRed)
	default:
		return baseStyle.Background(tcell.ColorWhite).Foreground(tcell.ColorWhite)
	}
}
