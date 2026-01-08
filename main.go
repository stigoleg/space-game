package main

import (
	"log"

	"stellar-siege/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file for configuration (GitHub tokens, etc.)
	// Ignore error if .env doesn't exist - we'll fall back to config file
	_ = godotenv.Load()

	g := game.NewGame()

	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("STELLAR SIEGE - Defend the Frontier")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
