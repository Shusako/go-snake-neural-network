package main

import (
	"math"

	"github.com/shusako/go_snake_neural_network/network"
	"github.com/shusako/go_snake_neural_network/snakegame/snake"
)

type neuralInput struct {
	feedForwardFunc network.FeedForward
}

func NewNeuralInput(feedForwardFunc network.FeedForward) *neuralInput {
	input := &neuralInput{}
	input.feedForwardFunc = feedForwardFunc

	return input
}

func (input *neuralInput) HandleInput(game *snake.Game, snake *snake.Snake) {
	encoding := EncodeGameBoard(game)
	output := input.feedForwardFunc(encoding)

	// find the maximum value of the output and set the direction to that
	maxIndex := 0
	for i := 0; i < len(output); i++ {
		if output[i] > output[maxIndex] {
			maxIndex = i
		}
	}

	// print the outputs in one console line
	// for i := 0; i < len(output); i++ {
	// 	fmt.Printf("%f ", output[i])
	// }
	// fmt.Printf(" %d\n", maxIndex)

	snake.TargetDirection = maxIndex
}

func ScanDirection(game *snake.Game, slopeX, slopeY int) (float64, float64, float64) {
	snakeDistance := math.MaxFloat64
	foodDistance := math.MaxFloat64
	totalDistance := 0.0

	foundFood := false
	foundSnake := false

	scanX := game.Snake.Head.X
	scanY := game.Snake.Head.Y

	// while scan position are in the game board bounds
	for scanX >= 0 && scanY >= 0 && scanX < snake.BoardWidth && scanY < snake.BoardHeight {
		scanX += slopeX
		scanY += slopeY
		totalDistance += 1

		// check if scan is food
		if !foundFood && scanX == game.Food.X && scanY == game.Food.Y {
			foodDistance = totalDistance
			foundFood = true
		}

		// check if scan is snake
		if !foundSnake && game.Snake.ContainsLocation(scanX, scanY, true) {
			snakeDistance = totalDistance
			foundSnake = true
		}
	}

	wallDistance := 1.0 / totalDistance
	snakeDistance = 1.0 / snakeDistance
	foodDistance = 1.0 / foodDistance

	return wallDistance, snakeDistance, foodDistance
}

func EncodeGameBoard(game *snake.Game) []float64 {
	//encoding := make([]float64, snake.BoardWidth*snake.BoardHeight+20)
	// #region encode all the board
	// encode the entire board if the snake is there vs not there
	// for x := 0; x < snake.BoardWidth; x++ {
	// 	for y := 0; y < snake.BoardHeight; y++ {
	// 		if game.Snake.ContainsLocation(x, y, true) {
	// 			encoding[x+y*snake.BoardWidth] = 1
	// 		} else {
	// 			encoding[x+y*snake.BoardWidth] = 0
	// 		}
	// 	}
	// }
	// board offset
	//offset := snake.BoardWidth * snake.BoardHeight
	// #endregion

	encoding := make([]float64, 3*8+20)
	startVision := 0
	encoding[startVision+0], encoding[startVision+1], encoding[startVision+2] = ScanDirection(game, 0, 1)
	encoding[startVision+3], encoding[startVision+4], encoding[startVision+5] = ScanDirection(game, 1, 1)
	encoding[startVision+6], encoding[startVision+7], encoding[startVision+8] = ScanDirection(game, 1, 0)
	encoding[startVision+9], encoding[startVision+10], encoding[startVision+11] = ScanDirection(game, 1, -1)
	encoding[startVision+12], encoding[startVision+13], encoding[startVision+14] = ScanDirection(game, 0, -1)
	encoding[startVision+15], encoding[startVision+16], encoding[startVision+17] = ScanDirection(game, -1, -1)
	encoding[startVision+18], encoding[startVision+19], encoding[startVision+20] = ScanDirection(game, -1, 0)
	encoding[startVision+21], encoding[startVision+22], encoding[startVision+23] = ScanDirection(game, -1, 1)

	offset := startVision

	encoding[offset+0] = float64(game.Food.X)
	encoding[offset+1] = float64(game.Food.Y)
	encoding[offset+2] = float64(game.Snake.Head.X)
	encoding[offset+3] = float64(game.Snake.Head.Y)

	// distance to apple, normalized to 0-1, 0 being on top of the apple, 1 being on the opposite corner
	encoding[offset+4] = float64(game.Snake.Head.X-game.Food.X) / float64(snake.BoardWidth)
	encoding[offset+5] = float64(game.Snake.Head.Y-game.Food.Y) / float64(snake.BoardHeight)

	// distance to left wall
	encoding[offset+6] = float64(game.Snake.Head.X) / float64(snake.BoardWidth)
	// distance to right wall
	encoding[offset+7] = float64(snake.BoardWidth-game.Snake.Head.X) / float64(snake.BoardWidth)
	// distance to top wall
	encoding[offset+8] = float64(game.Snake.Head.Y) / float64(snake.BoardHeight)
	// distance to bottom wall
	encoding[offset+9] = float64(snake.BoardHeight-game.Snake.Head.Y) / float64(snake.BoardHeight)

	// one hot encode the head direction as UP, RIGHT, DOWN, LEFT
	encoding[offset+10] = 0
	encoding[offset+11] = 0
	encoding[offset+12] = 0
	encoding[offset+13] = 0

	encoding[offset+10+game.Snake.Direction] = 1

	// set tail location to end of tail
	encoding[offset+14] = float64(game.Snake.Tail[len(game.Snake.Tail)-1].X)
	encoding[offset+15] = float64(game.Snake.Tail[len(game.Snake.Tail)-1].Y)
	// one hot encode the tail direction as UP, RIGHT, DOWN, LEFT
	encoding[offset+16] = 0
	encoding[offset+17] = 0
	encoding[offset+18] = 0
	encoding[offset+19] = 0
	// TODO: implement

	return encoding
}
