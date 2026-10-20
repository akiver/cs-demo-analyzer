package model

type BombDefused struct {
	Frame                      int     `json:"frame"`
	Tick                       int     `json:"tick"`
	RoundNumber                int     `json:"roundNumber"`
	Site                       string  `json:"site"`
	DefuserSteamID64           uint64  `json:"defuserSteamId"`
	DefuserName                string  `json:"defuserName"`
	IsPlayerControllingBot     bool    `json:"isPlayerControllingBot"`
	X                          float64 `json:"x"`
	Y                          float64 `json:"y"`
	Z                          float64 `json:"z"`
	CounterTerroristAliveCount int     `json:"counterTerroristAliveCount"`
	TerroristAliveCount        int     `json:"terroristAliveCount"`
}
