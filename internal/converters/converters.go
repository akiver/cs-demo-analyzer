package converters

import (
	"fmt"
	"strconv"

	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
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

func TeamToString(value common.Team) string {
	return ByteToString(byte(value))
}

func HitgroupToString(value events.HitGroup) string {
	return ByteToString(byte(value))
}

func RoundEndReasonToString(value events.RoundEndReason) string {
	return ByteToString(byte(value))
}

func ColorToString(color common.Color) string {
	return IntToString(int(color))
}

func StringToInt(value string) int {
	valueAsInt, err := strconv.Atoi(value)
	if err == nil {
		return valueAsInt
	}

	return 0
}

func BombsiteToString(bombSite events.Bombsite) string {
	if bombSite == 0 {
		return "Unknown"
	}

	return string(bombSite)
}
