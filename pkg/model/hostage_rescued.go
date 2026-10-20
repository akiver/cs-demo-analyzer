package model

type HostageRescued struct {
	Frame                  int     `json:"frame"`
	Tick                   int     `json:"tick"`
	RoundNumber            int     `json:"roundNumber"`
	PlayerSteamID64        uint64  `json:"playerSteamId"`
	IsPlayerControllingBot bool    `json:"isPlayerControllingBot"`
	HostageEntityId        int     `json:"hostageEntityId"`
	X                      float64 `json:"x"`
	Y                      float64 `json:"y"`
	Z                      float64 `json:"z"`
}
