package util

import (
	"math"
	"time"
)

func IfZero(arg interface{}) bool {
	if arg == nil {
		return true
	}
	switch v := arg.(type) {
	case int:
		if v==0 {
			return true
		}
	case int8:
		if int(v)==0 {
			return true
		}
	case int16:
		if v==0 {
			return true
		}
	case int32:
		if v==0 {
			return true
		}
	case int64:
		if v==0 {
			return true
		}
	case uint:
		if v==0 {
			return true
		}
	case uint8:
		if v==uint8(0) {
			return true
		}
	case uint16:
		if v==uint16(0) {
			return true
		}
	case uint32:
		if v==uint32(0) {
			return true
		}
	case uint64:
		if v==uint64(0) {
			return true
		}
	case float32:
		r := float64(v)
		return math.Abs(r-0) < 0.0000001
	case float64:
		return math.Abs(v-0) < 0.0000001
	case string:
		if v == "" || v == "%%" || v == "%" {
			return true
		}
	case *string, *int, *int64, *int32, *int16, *int8, *float32, *float64, *time.Time:
		if v == nil {
			return true
		}
	case time.Time:
		return v.IsZero()
	//case decimal.Decimal:
	//	tmp, _ := v.Float64()
	//	return math.Abs(tmp-0) < 0.0000001
	default:
		return false
	}
	return false
}
