package demo

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"
)

const (
	dictionary = "ABCDEFGHJKLMNOPQRSTUVWXYZabcdefhijkmnopqrstuvwxyz23456789"
)

type MatchInformation struct {
	MatchId       uint64
	ReservationId uint64
	TvPort        uint32
}

func bytesToShareCode(bytes []byte) string {
	dictionaryLength := big.NewInt(int64(len(dictionary)))
	total := new(big.Int).SetBytes(bytes)
	var builder strings.Builder
	builder.Grow(29) // "CSGO-" + 5*5 chars + 4*"-"

	for i := 0; i < 25; i++ {
		remainder := new(big.Int)
		total.DivMod(total, dictionaryLength, remainder)
		builder.WriteByte(dictionary[remainder.Int64()])
	}

	str := builder.String()

	return fmt.Sprintf("CSGO-%s-%s-%s-%s-%s", str[0:5], str[5:10], str[10:15], str[15:20], str[20:25])
}

func encodeMatchShareCode(match MatchInformation) string {
	bytes := make([]byte, 18)

	binary.LittleEndian.PutUint64(bytes[0:8], match.MatchId)
	binary.LittleEndian.PutUint64(bytes[8:16], match.ReservationId)
	binary.LittleEndian.PutUint16(bytes[16:18], uint16(match.TvPort))

	return bytesToShareCode(bytes)
}
