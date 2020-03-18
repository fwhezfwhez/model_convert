package model_convert

import "strings"

func LowerFistLetter(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(string(s[0])) + s[1:]
}
func UpperFirstLetter(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
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

func Format(arg string) string {
	arg = strings.Replace(arg, "\n\n", "\n", -1)
	arg = strings.Replace(arg, "\n\n\n", "\n", -1)
	arg = strings.Replace(arg, "\n    \n", "\n", -1)
	arg = strings.Replace(arg, "//\n\n//", "//\n//", -1)
	return arg
}

func StringMaxLen(max int, buf []byte) string {
	if len(buf) == 0 {
		return ""
	}
	if max >= len(buf) {
		return string(buf)
	}
	return string(buf[:max])
}
