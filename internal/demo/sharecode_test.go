package demo

import (
	"testing"
)

type Sample struct {
	shareCode string
	match     MatchInformation
}

func TestEncodeMatchShareCode(t *testing.T) {
	matchSamples := []Sample{
		{
			shareCode: "CSGO-L9spZ-ihuov-cyhtE-kxbqa-FkBAA",
			match: MatchInformation{
				MatchId:       3400360672356205056,
				ReservationId: 3400367402569957763,
				TvPort:        9725,
			},
		},
		{
			shareCode: "CSGO-GADqf-jjyJ8-cSP2r-smZRo-TO2xK",
			match: MatchInformation{
				MatchId:       3230642215713767580,
				ReservationId: 3230647599455273103,
				TvPort:        55788,
			},
		},
		{
			shareCode: "CSGO-bPQEz-PrYTq-u5w8E-ZbUy7-ZeQ3A",
			match: MatchInformation{
				MatchId:       3325408798641750542,
				ReservationId: 3325410334092558852,
				TvPort:        240,
			},
		},
		{
			shareCode: "CSGO-wBrm6-7fkM6-AzBC5-u6GmR-iHLHA",
			match: MatchInformation{
				MatchId:       3302232779302895618,
				ReservationId: 3302241568953467250,
				TvPort:        3085,
			},
		},
		{
			shareCode: "CSGO-TKDTJ-YrAXs-sDNfL-HOuKO-i84VH",
			match: MatchInformation{
				MatchId:       3402250361329680757,
				ReservationId: 3402250801563828781,
				TvPort:        61630,
			},
		},
		{
			shareCode: "CSGO-p4X9o-3Mfut-tpe5y-J8K6f-mj5ZJ",
			match: MatchInformation{
				MatchId:       3402249502336221574,
				ReservationId: 3402252092201501292,
				TvPort:        14119,
			},
		},
	}

	for _, sample := range matchSamples {
		shareCode := encodeMatchShareCode(sample.match)
		if shareCode != sample.shareCode {
			t.Errorf("Expected share code %s, got %s", sample.shareCode, shareCode)
		}
	}
}
