package constants

type EconomyType string

func (name EconomyType) String() string {
	return string(name)
}

const (
	EconomyTypePistol   EconomyType = "pistol"
	EconomyTypeEco      EconomyType = "eco"
	EconomyTypeSemi     EconomyType = "semi"
	EconomyTypeForceBuy EconomyType = "force-buy"
	EconomyTypeFull     EconomyType = "full"
)
