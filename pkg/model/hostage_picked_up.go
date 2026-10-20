package model

type HostagePickedUp struct {
	Frame                  int     `json:"frame"`
	Tick                   int     `json:"tick"`
	RoundNumber            int     `json:"roundNumber"`
	HostageEntityId        int     `json:"hostageEntityId"`
	PlayerSteamID64        uint64  `json:"playerSteamId"`
	IsPlayerControllingBot bool    `json:"isPlayerControllingBot"`
	X                      float64 `json:"x"`
	Y                      float64 `json:"y"`
	Z                      float64 `json:"z"`
}
