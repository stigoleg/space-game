package main

import (
	"log"

	"stellar-siege/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.NewGame()

	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("STELLAR SIEGE - Defend the Frontier")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
