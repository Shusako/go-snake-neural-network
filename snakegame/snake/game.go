package snake

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	BoardWidth  = 10
	BoardHeight = 10

	squareSize = 16

	UpDirection    = 0
	RightDirection = 1
	DownDirection  = 2
	LeftDirection  = 3
)

type Location struct {
	X int
	Y int
}

type Game struct {
	Snake  Snake
	Food   Location
	IsOver bool
	Input  Input

	// TODO game probably shouldn't know about this, should allow inherited inputs to control pace?
	// NEED neural input to be able to go as fast as it can compute
	nextTick time.Time
	TickMs   int64

	AutoRestart bool

	Random *rand.Source

	Moves int
}

func (g *Game) Reset() {
	g.Snake = Snake{
		Head: Location{X: 5, Y: 5},
		Tail: []Location{
			{X: 4, Y: 5},
			{X: 3, Y: 5},
			{X: 2, Y: 5},
		},
		TargetDirection: RightDirection,
		Direction:       RightDirection,
	}

	g.PlaceFood()

	// Reset everything
	g.Moves = 0
	g.IsOver = false
	g.Snake.IsDead = false
	g.Snake.AteFood = false
}

func (g *Game) PlaceFood() {
	// Place food at random location not occupied by snake
	for {
		g.Food.X = rand.Intn(BoardWidth)
		g.Food.Y = rand.Intn(BoardHeight)

		if !g.Snake.ContainsLocation(g.Food.X, g.Food.Y, true) {
			break
		}
	}
}

func DrawSquare(mainImage *ebiten.Image, x, y int, color color.RGBA) {
	vector.DrawFilledRect(mainImage, float32(x), float32(y), 1, 1, color, false)
}

func (g *Game) Draw(screen *ebiten.Image) {
	gameBoard := ebiten.NewImage(BoardWidth, BoardHeight)
	gameBoard.Fill(color.RGBA{0, 0, 0, 255})

	// Draw food
	DrawSquare(gameBoard, g.Food.X, g.Food.Y, color.RGBA{0, 255, 0, 255})

	// Draw snake
	DrawSquare(gameBoard, g.Snake.Head.X, g.Snake.Head.Y, color.RGBA{255, 0, 0, 255})
	for index, tail := range g.Snake.Tail {
		DrawSquare(gameBoard, tail.X, tail.Y, color.RGBA{uint8(255 - index*5), 0, uint8(index * 5), 255})
	}

	// Draw gameboard and scale it up
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(squareSize, squareSize)
	screen.DrawImage(gameBoard, op)
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.Reset()
	}

	if g.Snake.IsDead {
		g.IsOver = true
	}

	if g.IsOver {
		if g.AutoRestart {
			g.Reset()
		}
		return nil
	}

	g.Input.HandleInput(g, &g.Snake)

	// check if current time is past next tick
	if time.Now().Before(g.nextTick) {
		return nil
	}
	g.nextTick = time.Now().Add(time.Duration(g.TickMs) * time.Millisecond)

	g.Moves++

	g.Snake.Update()

	// Check if snake hit wall
	if g.Snake.Head.X < 0 || g.Snake.Head.X >= BoardWidth || g.Snake.Head.Y < 0 || g.Snake.Head.Y >= BoardHeight {
		g.Snake.IsDead = true
	}

	// Check if snake ate food
	if g.Snake.Head.X == g.Food.X && g.Snake.Head.Y == g.Food.Y {
		// Generate new food location
		g.Snake.AteFood = true
		g.PlaceFood()
	}

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return BoardWidth * squareSize, BoardHeight * squareSize
}

func RunGame(game *Game) {
	ebiten.SetWindowSize(640, 640*BoardHeight/BoardWidth)
	ebiten.SetWindowTitle("Snake Game")
	// ebiten.SetTPS(10) // Cannot set the TPS because inputs will get dropped
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
