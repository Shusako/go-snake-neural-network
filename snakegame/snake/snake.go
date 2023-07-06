package snake

type Snake struct {
	AteFood         bool
	IsDead          bool
	Head            Location
	Tail            []Location
	Direction       int
	TargetDirection int

	MovesSinceFood int
}

func (snake *Snake) ContainsLocation(x, y int, checkHead bool) bool {
	if checkHead && snake.Head.X == x && snake.Head.Y == y {
		return true
	}

	for _, tail := range snake.Tail {
		if tail.X == x && tail.Y == y {
			return true
		}
	}

	return false
}

func (snake *Snake) Update() {
	// Push head location to the tail
	snake.Tail = append([]Location{snake.Head}, snake.Tail...)

	// Check if snake target direction is invalid (opposite of current direction)
	if snake.TargetDirection == (snake.Direction+2)%4 {
		snake.TargetDirection = snake.Direction
	}

	// Move current head location based on direction
	snake.Direction = snake.TargetDirection
	switch snake.Direction {
	case UpDirection:
		snake.Head.Y--
	case RightDirection:
		snake.Head.X++
	case DownDirection:
		snake.Head.Y++
	case LeftDirection:
		snake.Head.X--
	}

	// Check if snake hit itself
	if snake.ContainsLocation(snake.Head.X, snake.Head.Y, false) {
		snake.IsDead = true
		return
	}

	// Pop last tail location if did not eat food
	if !snake.AteFood {
		snake.Tail = snake.Tail[:len(snake.Tail)-1]
		snake.MovesSinceFood++
	} else {
		snake.MovesSinceFood = 0
		snake.AteFood = false
	}

	if snake.MovesSinceFood > 100 {
		snake.IsDead = true
	}
}
