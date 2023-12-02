package constants

type TeamLetter string

func (name TeamLetter) String() string {
	return string(name)
}

const (
	TeamLetterA TeamLetter = "A"
	TeamLetterB TeamLetter = "B"
)
