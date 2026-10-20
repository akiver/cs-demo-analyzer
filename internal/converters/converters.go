package converters

import (
	"fmt"
	"strconv"
)

func IntToString(value int) string {
	return strconv.Itoa(value)
}

func Float32ToString(value float32) string {
	return fmt.Sprintf("%f", value)
}

func Float64ToString(value float64) string {
	return fmt.Sprintf("%f", value)
}

func Int64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func Uint32ToString(value uint32) string {
	return fmt.Sprint(value)
}

func Uint64ToString(value uint64) string {
	return strconv.FormatUint(value, 10)
}

func BoolToString(value bool) string {
	if value {
		return "1"
	}

	return "0"
}

func ByteToString(value byte) string {
	return fmt.Sprintf("%d", value)
}

// TeamToString converts a side (constants.Side) to its number.
func TeamToString[T ~byte](value T) string {
	return ByteToString(byte(value))
}

// HitgroupToString converts a hit group (constants.HitGroup) to its number.
func HitgroupToString[T ~byte](value T) string {
	return ByteToString(byte(value))
}

// RoundEndReasonToString converts a round end reason (constants.RoundEndReason) to its number.
func RoundEndReasonToString[T ~byte](value T) string {
	return ByteToString(byte(value))
}

// ColorToString converts a color (constants.Color) to its number.
func ColorToString[T ~int](color T) string {
	return IntToString(int(color))
}

func StringToInt(value string) int {
	valueAsInt, err := strconv.Atoi(value)
	if err == nil {
		return valueAsInt
	}

	return 0
}

// BombsiteToString converts a bombsite (A or B as a rune, 0 if unknown) to its name.
func BombsiteToString[T ~rune](bombSite T) string {
	if bombSite == 0 {
		return "Unknown"
	}

	return string(bombSite)
}
