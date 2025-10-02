# Go Tetris

A classic Tetris game implementation for the terminal, built with Go and the tcell library.

## Features

- **Classic Tetris Gameplay**: All 7 standard tetrominos (I, O, T, S, Z, J, L)
- **Smooth Controls**: Responsive keyboard input with proper collision detection
- **Progressive Difficulty**: Game speed increases with each level
- **Scoring System**: Classic Tetris scoring (1-4 lines with multipliers)
- **High Score Persistence**: Your best score is saved between sessions
- **Next Piece Preview**: See what's coming next
- **Pause Function**: Take a break without losing progress
- **Clean Terminal UI**: Cross-platform rendering with tcell

## Installation

### Prerequisites

- Go 1.21 or higher
- A terminal that supports ANSI colors

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd go-tetris

# Download dependencies
go mod download

# Build the game
go build -o tetris

# Run the game
./tetris
```

Or run directly without building:

```bash
go run .
```

## Controls

| Key | Action |
|-----|--------|
| **←** | Move piece left |
| **→** | Move piece right |
| **↓** | Soft drop (move down faster) |
| **↑** | Rotate piece clockwise |
| **Space** | Pause/unpause game |
| **Q** or **ESC** | Quit game |
| **R** | Restart (when game over) |

## Gameplay

### Objective
Stack falling tetromino pieces to create complete horizontal lines. When a line is completed, it clears and you earn points. The game ends when pieces stack up to the top of the board.

### Scoring
- **1 Line**: 40 × (level + 1) points
- **2 Lines**: 100 × (level + 1) points
- **3 Lines**: 300 × (level + 1) points
- **4 Lines (Tetris)**: 1200 × (level + 1) points

### Level Progression
- Start at Level 1
- Advance one level for every 10 lines cleared
- Pieces fall faster as level increases

## Project Structure

```
.
├── main.go       # Entry point and game loop
├── game.go       # Core game state and logic
├── board.go      # Board representation and operations
├── piece.go      # Tetromino definitions and behaviors
├── renderer.go   # Terminal rendering with tcell
├── input.go      # Keyboard input handling
├── score.go      # Scoring and high score persistence
├── go.mod        # Go module dependencies
└── README.md     # This file
```

## Technical Details

- **Language**: Go
- **Library**: [tcell](https://github.com/gdamore/tcell) for cross-platform terminal rendering
- **Board Size**: 10 columns × 20 rows (classic Tetris dimensions)
- **Refresh Rate**: 60 FPS game loop
- **Platform Support**: Linux, macOS, Windows (via tcell)

## Development

### Running Tests
```bash
go test ./...
```

### Code Organization
- **Separation of Concerns**: Game logic, rendering, and input are cleanly separated
- **No Global State**: All state is encapsulated in the `Game` struct
- **Proper Cleanup**: Terminal state is restored on exit via defer

## License

This project is open source and available for educational purposes.

## Acknowledgments

Built as a demonstration of terminal-based game development in Go using the excellent tcell library.
