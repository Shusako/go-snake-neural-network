package main

import snakegame "github.com/shusako/go_snake_neural_network/snakegame/snake"

func main() {
	game := snakegame.Game{}
	game.Reset()
	game.Input = &snakegame.UserInput{}
	game.Input.Init()
	game.TickMs = 100

	snakegame.RunGame(&game)
}
