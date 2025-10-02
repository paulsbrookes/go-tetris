package main

import (
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
)

func main() {
	// Initialize random seed
	// Note: Go 1.20+ auto-seeds, but for older versions you might need rand.Seed(time.Now().UnixNano())
	
	// Initialize tcell screen
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
	}
	
	if err := screen.Init(); err != nil {
		log.Fatalf("Failed to initialize screen: %v", err)
	}
	
	// Ensure cleanup on exit
	defer screen.Fini()
	
	// Set up screen
	screen.Clear()
	screen.EnableMouse()
	screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	
	// Create game and renderer
	game := NewGame()
	renderer := NewRenderer(screen)
	
	// Initial render
	renderer.Render(game)
	
	// Create a ticker for game updates (60 FPS)
	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()
	
	// Event channel for input
	quit := false
	
	// Main game loop
	for !quit {
		select {
		case <-ticker.C:
			// Update game logic
			game.Update()
			
			// Render
			renderer.Render(game)
			
		default:
			// Poll for events (non-blocking)
			if screen.HasPendingEvent() {
				ev := screen.PollEvent()
				
				switch ev := ev.(type) {
				case *tcell.EventKey:
					// Handle keyboard input
					if HandleInput(game, ev) {
						quit = true
					}
					
					// Re-render immediately after input
					renderer.Render(game)
					
				case *tcell.EventResize:
					// Handle terminal resize
					screen.Sync()
					renderer.Render(game)
				}
			}
			
			// Small sleep to prevent busy-waiting
			time.Sleep(time.Millisecond)
		}
	}
}
