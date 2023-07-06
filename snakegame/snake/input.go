package snake

type Input interface {
	HandleInput(game *Game, snake *Snake)
}
