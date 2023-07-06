package snake

import "github.com/hajimehoshi/ebiten/v2"

type UserInput struct {
}

func (input *UserInput) HandleInput(game *Game, snake *Snake) {
	// Handle updating head direction based on key input
	if ebiten.IsKeyPressed(ebiten.KeyUp) && snake.Direction != DownDirection {
		snake.TargetDirection = UpDirection
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) && snake.Direction != LeftDirection {
		snake.TargetDirection = RightDirection
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) && snake.Direction != UpDirection {
		snake.TargetDirection = DownDirection
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) && snake.Direction != RightDirection {
		snake.TargetDirection = LeftDirection
	}
}
