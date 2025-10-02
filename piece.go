package main

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// PieceType represents the type of tetromino
type PieceType int

const (
	PieceI PieceType = iota
	PieceO
	PieceT
	PieceS
	PieceZ
	PieceJ
	PieceL
)

// Piece represents a tetromino piece
type Piece struct {
	Type     PieceType
	X, Y     int // Position on board (top-left of bounding box)
	Rotation int // Current rotation state (0-3)
}

// Shape returns the current shape of the piece based on rotation
func (p *Piece) Shape() [][]bool {
	return pieceShapes[p.Type][p.Rotation]
}

// Color returns the color for this piece type
func (p *Piece) Color() tcell.Color {
	return pieceColors[p.Type]
}

// Rotate rotates the piece clockwise
func (p *Piece) Rotate() {
	p.Rotation = (p.Rotation + 1) % len(pieceShapes[p.Type])
}

// UnRotate rotates the piece counter-clockwise (for undo)
func (p *Piece) UnRotate() {
	p.Rotation = (p.Rotation - 1 + len(pieceShapes[p.Type])) % len(pieceShapes[p.Type])
}

// GetBlocks returns the absolute positions of all blocks in the piece
func (p *Piece) GetBlocks() []struct{ X, Y int } {
	shape := p.Shape()
	blocks := []struct{ X, Y int }{}
	
	for y := 0; y < len(shape); y++ {
		for x := 0; x < len(shape[y]); x++ {
			if shape[y][x] {
				blocks = append(blocks, struct{ X, Y int }{p.X + x, p.Y + y})
			}
		}
	}
	
	return blocks
}

// NewRandomPiece creates a random tetromino piece at the spawn position
func NewRandomPiece() *Piece {
	pieceType := PieceType(rand.Intn(7))
	return &Piece{
		Type:     pieceType,
		X:        3, // Spawn in middle-ish of 10-wide board
		Y:        0,
		Rotation: 0,
	}
}

// pieceColors defines the color for each piece type
var pieceColors = map[PieceType]tcell.Color{
	PieceI: tcell.ColorCyan,
	PieceO: tcell.ColorYellow,
	PieceT: tcell.ColorPurple,
	PieceS: tcell.ColorGreen,
	PieceZ: tcell.ColorRed,
	PieceJ: tcell.ColorBlue,
	PieceL: tcell.ColorOrange,
}

// pieceShapes defines all rotation states for each piece type
// Each shape is a 2D array where true = block, false = empty
var pieceShapes = map[PieceType][][][]bool{
	// I piece - 4 blocks in a line
	PieceI: {
		{ // Rotation 0 (horizontal)
			{false, false, false, false},
			{true, true, true, true},
			{false, false, false, false},
			{false, false, false, false},
		},
		{ // Rotation 1 (vertical)
			{false, false, true, false},
			{false, false, true, false},
			{false, false, true, false},
			{false, false, true, false},
		},
		{ // Rotation 2 (horizontal)
			{false, false, false, false},
			{false, false, false, false},
			{true, true, true, true},
			{false, false, false, false},
		},
		{ // Rotation 3 (vertical)
			{false, true, false, false},
			{false, true, false, false},
			{false, true, false, false},
			{false, true, false, false},
		},
	},
	
	// O piece - 2x2 square (no rotation needed, but we define 1 state)
	PieceO: {
		{
			{true, true},
			{true, true},
		},
	},
	
	// T piece - T shape
	PieceT: {
		{ // Rotation 0
			{false, true, false},
			{true, true, true},
			{false, false, false},
		},
		{ // Rotation 1
			{false, true, false},
			{false, true, true},
			{false, true, false},
		},
		{ // Rotation 2
			{false, false, false},
			{true, true, true},
			{false, true, false},
		},
		{ // Rotation 3
			{false, true, false},
			{true, true, false},
			{false, true, false},
		},
	},
	
	// S piece - S shape
	PieceS: {
		{ // Rotation 0
			{false, true, true},
			{true, true, false},
			{false, false, false},
		},
		{ // Rotation 1
			{false, true, false},
			{false, true, true},
			{false, false, true},
		},
		{ // Rotation 2
			{false, false, false},
			{false, true, true},
			{true, true, false},
		},
		{ // Rotation 3
			{true, false, false},
			{true, true, false},
			{false, true, false},
		},
	},
	
	// Z piece - Z shape
	PieceZ: {
		{ // Rotation 0
			{true, true, false},
			{false, true, true},
			{false, false, false},
		},
		{ // Rotation 1
			{false, false, true},
			{false, true, true},
			{false, true, false},
		},
		{ // Rotation 2
			{false, false, false},
			{true, true, false},
			{false, true, true},
		},
		{ // Rotation 3
			{false, true, false},
			{true, true, false},
			{true, false, false},
		},
	},
	
	// J piece - J shape
	PieceJ: {
		{ // Rotation 0
			{true, false, false},
			{true, true, true},
			{false, false, false},
		},
		{ // Rotation 1
			{false, true, true},
			{false, true, false},
			{false, true, false},
		},
		{ // Rotation 2
			{false, false, false},
			{true, true, true},
			{false, false, true},
		},
		{ // Rotation 3
			{false, true, false},
			{false, true, false},
			{true, true, false},
		},
	},
	
	// L piece - L shape
	PieceL: {
		{ // Rotation 0
			{false, false, true},
			{true, true, true},
			{false, false, false},
		},
		{ // Rotation 1
			{false, true, false},
			{false, true, false},
			{false, true, true},
		},
		{ // Rotation 2
			{false, false, false},
			{true, true, true},
			{true, false, false},
		},
		{ // Rotation 3
			{true, true, false},
			{false, true, false},
			{false, true, false},
		},
	},
}
