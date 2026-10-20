package model

type ChickenPosition struct {
	// Unique ID of the chicken in the match, not exported because it's used only to group positions by chicken.
	chickenID   int
	Frame       int     `json:"frame"`
	Tick        int     `json:"tick"`
	RoundNumber int     `json:"roundNumber"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Z           float64 `json:"z"`
}
