package main

import (
	"fmt"
	"image/color"
	"math"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/shusako/go_snake_neural_network/network"
	"github.com/shusako/go_snake_neural_network/snakegame/snake"
	"golang.org/x/image/font/basicfont"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 720
)

type EvolutionManager struct {
	game          *snake.Game
	nextGame      *snake.Game
	nextGameMutex *sync.Mutex

	nextTick time.Time
	TickMs   int64

	geneticAlgorithm *network.GeneticAlgorithm
}

func NewEvolutionManager() *EvolutionManager {
	manager := &EvolutionManager{}

	populationSize := 1300
	// sizes := []int{120, 18, 18, 4}
	sizes := []int{3*8 + 20, 18, 18, 4}
	// TODO: Lookup what these values actually translate to in production algorithms so we match
	mutationChance := 0.01
	mutationRate := 0.1

	ga := network.NewGeneticAlgoritm(populationSize, sizes, mutationChance, mutationRate, func(feedForward network.FeedForward) float64 {
		game := &snake.Game{}
		game.Reset()
		game.Input = NewNeuralInput(feedForward)
		game.TickMs = -1
		for !game.IsOver {
			game.Update()
		}

		return GetFitness(game)
	})

	manager.geneticAlgorithm = ga
	manager.nextGameMutex = &sync.Mutex{}

	return manager
}

func GetFitness(game *snake.Game) float64 {
	// fitness should be based on snake length minus number of moves

	apples := float64(len(game.Snake.Tail) - 3)
	steps := float64(game.Moves)

	return steps + (math.Pow(2, apples) + math.Pow(apples, 2.1)*500) - (math.Pow(apples, 1.2) * math.Pow(0.25*steps, 1.3))

	//return math.Pow(1.5, float64(len(game.Snake.Tail))) - (math.Max(0, float64(game.Moves-10)) / 10)
}

func DrawSquare(mainImage *ebiten.Image, x, y int, color color.RGBA) {
	vector.DrawFilledRect(mainImage, float32(x), float32(y), 1, 1, color, false)
}

func (g *EvolutionManager) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x21, 0x21, 0x21, 0xff})

	if g.game != nil {
		DrawSnakeGame(g.game, screen)
		text.Draw(screen, fmt.Sprintf("Fitness: %f", GetFitness(g.game)), basicfont.Face7x13, 2, 12, color.White)
	}
}

func DrawSnakeGame(snakegame *snake.Game, screen *ebiten.Image) {
	// draw the game on the right side of the screen
	snakeGameImage := ebiten.NewImage(snakegame.Layout(0, 0))
	snakegame.Draw(snakeGameImage)

	// draw snakeGameImage on right side of screen at 650, 10 to 1270, 710
	op := &ebiten.DrawImageOptions{}
	padding := 10
	topLeftX := ScreenWidth/2 + padding
	topLeftY := padding
	bottomRightX := ScreenWidth - padding
	bottomRightY := ScreenHeight - padding
	sideScale := math.Min(float64(bottomRightX-topLeftX)/float64(snakeGameImage.Bounds().Dx()), float64(bottomRightY-topLeftY)/float64(snakeGameImage.Bounds().Dy()))
	op.GeoM.Scale(sideScale, sideScale)
	op.GeoM.Translate(float64(topLeftX), float64(topLeftY))
	screen.DrawImage(snakeGameImage, op)
}

func (g *EvolutionManager) Update() error {
	g.nextGameMutex.Lock()
	if g.nextGame != nil {
		if (g.game == nil) || (g.game.IsOver) {
			g.game = g.nextGame
			g.nextGame = nil
		}
	}
	g.nextGameMutex.Unlock()

	if g.game != nil {
		// check if current time is past next tick
		if time.Now().After(g.nextTick) {
			g.game.Update() // ignoring error
			g.nextTick = time.Now().Add(time.Duration(g.TickMs) * time.Millisecond)
		}
	}

	return nil
}

func (g *EvolutionManager) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}

func CreateGameFromFeedForward(feedForward network.FeedForward) *snake.Game {
	game := &snake.Game{}
	game.Reset()
	game.Input = NewNeuralInput(feedForward)
	game.TickMs = -1
	//game.AutoRestart = true

	return game
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Snake Evolution")

	manager := NewEvolutionManager()

	// manager.game = &snake.Game{}
	// manager.game.Reset()
	// manager.game.Input = &snake.UserInput{}
	// manager.game.TickMs = 100
	// manager.TickMs = -1

	manager.TickMs = 100

	// blank goroutine with loop
	go func() {
		for {
			manager.geneticAlgorithm.EvaluateGeneration()
			nextGame := CreateGameFromFeedForward(manager.geneticAlgorithm.GetBestIndividual())

			manager.nextGameMutex.Lock()
			manager.nextGame = nextGame
			manager.nextGameMutex.Unlock()

			manager.geneticAlgorithm.EvolveGeneration()
		}
	}()

	if err := ebiten.RunGame(manager); err != nil {
		panic(err)
	}
}
