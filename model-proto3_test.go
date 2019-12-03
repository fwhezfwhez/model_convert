package model_convert

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGoModelToProto3(t *testing.T) {
	type ArrayElement struct {
	}
	type U struct {
		Arr     []ArrayElement
		Config2 []byte
		Config  json.RawMessage

		Username string
		Password string
		Age      int
		Id       int32
	}
	ps, setM, setP := GoModelToProto3(U{}, map[string]string{
		"${pb_pkg_name}":    "userProto",
		"${model_pkg_name}": "userModel",
	})
	fmt.Println(ps)
	fmt.Println(setM)
	fmt.Println(setP)
}

func TestGoModelToProto2(t *testing.T) {
	type PersonalResult struct {
		GameId     int
		GameAreaId int
		Result     int
		Delta      int
		Upgrade    int
		Coin       int
	}
	ps, setM, setP := GoModelToProto2(PersonalResult{}, map[string]string{
		"${pb_pkg_name}":    "gamePb",
		"${model_pkg_name}": "${empty}",
//		"${start_index}": "1",
	})
	fmt.Println(ps)
	fmt.Println(setM)
	fmt.Println(setP)
}
