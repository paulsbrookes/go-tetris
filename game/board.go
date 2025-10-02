package game

// CheckLines scans the board from bottom to top and returns indices of completed rows
func CheckLines(board [][]int) []int {
	if len(board) == 0 {
		return nil
	}
	
	width := len(board[0])
	completedRows := []int{}
	
	// Scan from bottom to top
	for y := len(board) - 1; y >= 0; y-- {
		full := true
		for x := 0; x < width; x++ {
			if board[y][x] == 0 {
				full = false
				break
			}
		}
		
		if full {
			completedRows = append(completedRows, y)
		}
	}
	
	return completedRows
}

// ClearLines removes the specified rows and collapses the board
// Rows should be provided as indices (sorted or unsorted)
func ClearLines(board [][]int, rows []int) {
	if len(board) == 0 || len(rows) == 0 {
		return
	}
	
	width := len(board[0])
	height := len(board)
	
	// Create a map for quick lookup of rows to clear
	rowsToRemove := make(map[int]bool)
	for _, row := range rows {
		if row >= 0 && row < height {
			rowsToRemove[row] = true
		}
	}
	
	// Build new board without cleared lines
	newBoard := [][]int{}
	for y := 0; y < height; y++ {
		if !rowsToRemove[y] {
			newBoard = append(newBoard, board[y])
		}
	}
	
	// Add empty rows at the top
	for i := 0; i < len(rows); i++ {
		emptyRow := make([]int, width)
		newBoard = append([][]int{emptyRow}, newBoard...)
	}
	
	// Copy new board back to original
	for y := 0; y < height; y++ {
		copy(board[y], newBoard[y])
	}
}
