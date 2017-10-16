package engine
import(
	"fmt"
	"time"
	"github.com/pborman/uuid"

)

func (position Coordinate) String() string {
	return fmt.Sprintf("(%d,%d)", position.X, position.Y)
}

func (move Move) String() string {
	return fmt.Sprintf("%d/%s", move.Player, move.Position)
}


func NewMatch(size int, playerBlack string,playWhite string) Match{
	result := Match{}
	result.ID = uuid.New()
	result.GameBoard = newBoard(size)
	result.TurnCount = 0
	result.GridSize = size
	result.StartTime = time.Now()
	result.PlayerBlack = playerBlack
	result.PlayerWhite = playWhite	
	return result
}