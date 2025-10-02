package main

import (
	"os"
	"strconv"
	"strings"
)

const highScoreFile = "highscore.txt"

// CalculateScore calculates points for lines cleared
// Classic Tetris scoring: 1 line = 40, 2 lines = 100, 3 lines = 300, 4 lines = 1200
// Multiplied by (level + 1)
func CalculateScore(linesCleared, level int) int {
	baseScore := 0
	
	switch linesCleared {
	case 1:
		baseScore = 40
	case 2:
		baseScore = 100
	case 3:
		baseScore = 300
	case 4:
		baseScore = 1200
	default:
		return 0
	}
	
	return baseScore * (level + 1)
}

// LoadHighScore loads the high score from file
func LoadHighScore() int {
	data, err := os.ReadFile(highScoreFile)
	if err != nil {
		return 0 // No high score file yet
	}
	
	scoreStr := strings.TrimSpace(string(data))
	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		return 0
	}
	
	return score
}

// SaveHighScore saves the high score to file
func SaveHighScore(score int) error {
	data := strconv.Itoa(score)
	return os.WriteFile(highScoreFile, []byte(data), 0644)
}
