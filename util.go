package model_convert

import "strings"

func LowerFistLetter(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

func GetZeroValue(src interface{}) string {
	switch src.(type) {
	case int, int8, int16, int32, int64, float32, float64, uint, uint8, uint16, uint32, uint64:
		return "0"
	case string:
		return `""`
	}
	return `""`
}
