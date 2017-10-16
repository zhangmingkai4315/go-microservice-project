package engine

import(
	"fmt"
)

func genMove(x int, y int, player byte) Move {
	return Move{Player: player, Position: Coordinate{X: x, Y: y}}
}

// Copy copies the state of an existing board into a new board.
func (gameboard GameBoard) copy() GameBoard {
	outBoard := newBoard(cap(gameboard.Positions))
	copy(outBoard.Positions, gameboard.Positions)
	return outBoard
}

// NewBoard creates a new gameboard of a given size. Gameboards must always be square.
func newBoard(size int) GameBoard {
	outBoard := GameBoard{}
	a := make([][]byte, size)
	for i := range a {
		a[i] = make([]byte, size)
	}
	outBoard.Positions = a
	return outBoard
}

func (gameboard GameBoard) PerformMove(move Move) (outBoard GameBoard, err error) {
	//outBoard = NewBoard(cap(gameboard.Positions))
	outBoard = gameboard.copy()

	// TODO - this should eventually be part of another method that runs through
	// all the game rules in appropriate order.
	if gameboard.Positions[move.Position.X][move.Position.Y] != EmptyPosition {
		return outBoard, fmt.Errorf(RulesFailureSpaceOccupied, move)
	}

	outBoard.Positions[move.Position.X][move.Position.Y] = move.Player
	return outBoard, err
}