package constants

type Game string

const (
	CSGO  Game = "CSGO"
	CS2   Game = "CS2"
	CS2LT Game = "CS2 LT"
)

func (game Game) String() string {
	return string(game)
}
