package constants

// RoundEndReason is why a round ended. The values are the game's round end reason numbers.
type RoundEndReason byte

const (
	RoundEndReasonStillInProgress      RoundEndReason = 0
	RoundEndReasonTargetBombed         RoundEndReason = 1
	RoundEndReasonVIPEscaped           RoundEndReason = 2
	RoundEndReasonVIPKilled            RoundEndReason = 3
	RoundEndReasonTerroristsEscaped    RoundEndReason = 4
	RoundEndReasonCTStoppedEscape      RoundEndReason = 5
	RoundEndReasonTerroristsStopped    RoundEndReason = 6
	RoundEndReasonBombDefused          RoundEndReason = 7
	RoundEndReasonCTWin                RoundEndReason = 8
	RoundEndReasonTerroristsWin        RoundEndReason = 9
	RoundEndReasonDraw                 RoundEndReason = 10
	RoundEndReasonHostagesRescued      RoundEndReason = 11
	RoundEndReasonTargetSaved          RoundEndReason = 12
	RoundEndReasonHostagesNotRescued   RoundEndReason = 13
	RoundEndReasonTerroristsNotEscaped RoundEndReason = 14
	RoundEndReasonVIPNotEscaped        RoundEndReason = 15
	RoundEndReasonGameStart            RoundEndReason = 16
	RoundEndReasonTerroristsSurrender  RoundEndReason = 17
	RoundEndReasonCTSurrender          RoundEndReason = 18
	RoundEndReasonTerroristsPlanted    RoundEndReason = 19
	RoundEndReasonCTsReachedHostage    RoundEndReason = 20
)
