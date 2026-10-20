package model

type BombExploded struct {
	Frame                  int     `json:"frame"`
	Tick                   int     `json:"tick"`
	RoundNumber            int     `json:"roundNumber"`
	Site                   string  `json:"site"`
	PlanterSteamID64       uint64  `json:"defuserSteamId"`
	PlanterName            string  `json:"defuserName"`
	IsPlayerControllingBot bool    `json:"isPlayerControllingBot"`
	X                      float64 `json:"x"`
	Y                      float64 `json:"y"`
	Z                      float64 `json:"z"`
}
