package game

import "time"

// Scoring constants for line clears
const (
	ScoreSingle = 100  // 1 line cleared
	ScoreDouble = 300  // 2 lines cleared
	ScoreTriple = 500  // 3 lines cleared
	ScoreTetris = 800  // 4 lines cleared (Tetris)
)

// CalculateScore calculates the score for a given number of lines cleared at a given level
// Base points are multiplied by the current level
func CalculateScore(linesCleared int, level int) int {
	if linesCleared <= 0 {
		return 0
	}
	
	var baseScore int
	switch linesCleared {
	case 1:
		baseScore = ScoreSingle
	case 2:
		baseScore = ScoreDouble
	case 3:
		baseScore = ScoreTriple
	case 4:
		baseScore = ScoreTetris
	default:
		// For more than 4 lines (shouldn't happen in standard Tetris)
		// Award Tetris score
		baseScore = ScoreTetris
	}
	
	return baseScore * level
}

// GetGravitySpeed returns the drop interval (gravity speed) for a given level
// Uses the formula: max(100ms, 1000ms - (level * 80ms))
func GetGravitySpeed(level int) time.Duration {
	if level <= 0 {
		level = 1
	}
	
	// Start at 1000ms, decrease by 80ms per level, minimum 100ms
	interval := 1000 - ((level - 1) * 80)
	if interval < 100 {
		interval = 100
	}
	
	return time.Duration(interval) * time.Millisecond
}
