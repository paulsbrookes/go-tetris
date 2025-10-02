package input

import (
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Handler manages screen initialization and input event handling
type Handler struct {
	screen        tcell.Screen
	actionCh      chan Action
	quit          chan struct{}
	once          sync.Once
	lastInputTime map[Action]time.Time
	inputMutex    sync.Mutex
	debounceTime  time.Duration
}

// NewHandler creates a new input handler
func NewHandler() *Handler {
	return &Handler{
		actionCh:      make(chan Action, 10),
		quit:          make(chan struct{}),
		lastInputTime: make(map[Action]time.Time),
		debounceTime:  50 * time.Millisecond,
	}
}

// InitScreen initializes the tcell screen
func (h *Handler) InitScreen() error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	
	if err := screen.Init(); err != nil {
		return err
	}
	
	// Enable mouse (optional, but good practice)
	screen.EnableMouse()
	
	screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	screen.Clear()
	
	h.screen = screen
	return nil
}

// GetScreen returns the tcell screen instance
func (h *Handler) GetScreen() tcell.Screen {
	return h.screen
}

// Start begins listening for input events in a goroutine
func (h *Handler) Start() {
	go h.inputLoop()
}

// inputLoop continuously polls for events and converts them to actions
func (h *Handler) inputLoop() {
	for {
		select {
		case <-h.quit:
			return
		default:
			ev := h.screen.PollEvent()
			if ev == nil {
				continue
			}
			
			action := h.mapEventToAction(ev)
			if action != ActionNone && h.shouldProcessAction(action) {
				select {
				case h.actionCh <- action:
				default:
					// Drop action if channel is full
				}
			}
		}
	}
}

// shouldProcessAction implements debouncing by checking if enough time has passed
// since the last occurrence of this action
func (h *Handler) shouldProcessAction(action Action) bool {
	h.inputMutex.Lock()
	defer h.inputMutex.Unlock()
	
	now := time.Now()
	lastTime, exists := h.lastInputTime[action]
	
	// If this action hasn't been seen before or enough time has passed
	if !exists || now.Sub(lastTime) >= h.debounceTime {
		h.lastInputTime[action] = now
		return true
	}
	
	return false
}

// mapEventToAction converts tcell events to game actions based on key mappings
func (h *Handler) mapEventToAction(ev tcell.Event) Action {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		return h.mapKeyToAction(ev)
	case *tcell.EventResize:
		h.screen.Sync()
		return ActionNone
	default:
		return ActionNone
	}
}

// mapKeyToAction maps keyboard input to game actions
func (h *Handler) mapKeyToAction(ev *tcell.EventKey) Action {
	// Handle special keys
	switch ev.Key() {
	case tcell.KeyEscape, tcell.KeyCtrlC:
		return ActionQuit
	case tcell.KeyLeft:
		return ActionMoveLeft
	case tcell.KeyRight:
		return ActionMoveRight
	case tcell.KeyDown:
		return ActionMoveDown
	case tcell.KeyUp:
		return ActionRotateClockwise
	case tcell.KeyEnter:
		return ActionPause
	case tcell.KeyRune:
		// Handle character keys
		switch ev.Rune() {
		case 'q', 'Q':
			return ActionQuit
		case 'a', 'A':
			return ActionMoveLeft
		case 'd', 'D':
			return ActionMoveRight
		case 's', 'S':
			return ActionMoveDown
		case 'w', 'W':
			return ActionRotateClockwise
		case 'z', 'Z':
			return ActionRotateCounterClockwise
		case ' ':
			return ActionHardDrop
		case 'p', 'P':
			return ActionPause
		case 'r', 'R':
			return ActionRestart
		}
	}
	
	return ActionNone
}

// Actions returns the channel for receiving input actions
func (h *Handler) Actions() <-chan Action {
	return h.actionCh
}

// Stop stops the input handler and cleans up resources
func (h *Handler) Stop() {
	h.once.Do(func() {
		close(h.quit)
		if h.screen != nil {
			h.screen.Fini()
		}
		close(h.actionCh)
	})
}
