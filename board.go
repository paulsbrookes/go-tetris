package main

import "github.com/gdamore/tcell/v2"

const (
	BoardWidth  = 10
	BoardHeight = 20
)

// Cell represents a single cell on the board
type Cell struct {
	Occupied bool
	Color    tcell.Color
}

// Board represents the game board
type Board struct {
	Grid [BoardHeight][BoardWidth]Cell
}

// NewBoard creates a new empty board
func NewBoard() *Board {
	return &Board{}
}

// IsValid checks if a position is within board bounds
func (b *Board) IsValid(x, y int) bool {
	return x >= 0 && x < BoardWidth && y >= 0 && y < BoardHeight
}

// IsOccupied checks if a cell is occupied
func (b *Board) IsOccupied(x, y int) bool {
	if !b.IsValid(x, y) {
		return true // Out of bounds counts as occupied
	}
	return b.Grid[y][x].Occupied
}

// CanPlacePiece checks if a piece can be placed at its current position
func (b *Board) CanPlacePiece(piece *Piece) bool {
	blocks := piece.GetBlocks()
	
	for _, block := range blocks {
		// Check bounds
		if block.X < 0 || block.X >= BoardWidth || block.Y >= BoardHeight {
			return false
		}
		// Allow negative Y (above board) during spawn
		if block.Y < 0 {
			continue
		}
		// Check if cell is occupied
		if b.Grid[block.Y][block.X].Occupied {
			return false
		}
	}
	
	return true
}

// LockPiece locks a piece onto the board
func (b *Board) LockPiece(piece *Piece) {
	blocks := piece.GetBlocks()
	color := piece.Color()
	
	for _, block := range blocks {
		if b.IsValid(block.X, block.Y) {
			b.Grid[block.Y][block.X] = Cell{
				Occupied: true,
				Color:    color,
			}
		}
	}
}

// IsLineFull checks if a line is completely filled
func (b *Board) IsLineFull(y int) bool {
	if y < 0 || y >= BoardHeight {
		return false
	}
	
	for x := 0; x < BoardWidth; x++ {
		if !b.Grid[y][x].Occupied {
			return false
		}
	}
	
	return true
}

// ClearLine clears a specific line
func (b *Board) ClearLine(y int) {
	if y < 0 || y >= BoardHeight {
		return
	}
	
	// Shift all lines above down by one
	for row := y; row > 0; row-- {
		b.Grid[row] = b.Grid[row-1]
	}
	
	// Clear the top line
	b.Grid[0] = [BoardWidth]Cell{}
}

// ClearFullLines clears all full lines and returns the number cleared
func (b *Board) ClearFullLines() int {
	linesCleared := 0
	
	// Check from bottom to top
	y := BoardHeight - 1
	for y >= 0 {
		if b.IsLineFull(y) {
			b.ClearLine(y)
			linesCleared++
			// Don't decrement y, check same line again (since lines shifted down)
		} else {
			y--
		}
	}
	
	return linesCleared
}

// Clear resets the board to empty state
func (b *Board) Clear() {
	b.Grid = [BoardHeight][BoardWidth]Cell{}
}
