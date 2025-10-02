package input

// Action represents a user input action in the game
type Action int

const (
	// ActionNone represents no action
	ActionNone Action = iota
	// ActionMoveLeft moves the piece left
	ActionMoveLeft
	// ActionMoveRight moves the piece right
	ActionMoveRight
	// ActionMoveDown moves the piece down (soft drop)
	ActionMoveDown
	// ActionRotateClockwise rotates the piece clockwise
	ActionRotateClockwise
	// ActionRotateCounterClockwise rotates the piece counter-clockwise
	ActionRotateCounterClockwise
	// ActionHardDrop instantly drops the piece to the bottom
	ActionHardDrop
	// ActionPause pauses/unpauses the game
	ActionPause
	// ActionQuit exits the game
	ActionQuit
	// ActionRestart restarts the game
	ActionRestart
)

// String returns the string representation of an Action
func (a Action) String() string {
	switch a {
	case ActionNone:
		return "None"
	case ActionMoveLeft:
		return "MoveLeft"
	case ActionMoveRight:
		return "MoveRight"
	case ActionMoveDown:
		return "MoveDown"
	case ActionRotateClockwise:
		return "RotateClockwise"
	case ActionRotateCounterClockwise:
		return "RotateCounterClockwise"
	case ActionHardDrop:
		return "HardDrop"
	case ActionPause:
		return "Pause"
	case ActionQuit:
		return "Quit"
	case ActionRestart:
		return "Restart"
	default:
		return "Unknown"
	}
}
