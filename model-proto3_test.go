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
	type ClientAuthPayload struct {
		OpenId  string `json:"open_id"`
		UserId  int    `json:"user_id"`
		GameId  int    `json:"game_id"`
		AppId   string `json:"app_id"`
		Version int    `json:"version"`
		Exp     int64  `json:"exp"`
	}
	ps, setM, setP := GoModelToProto2(ClientAuthPayload{}, map[string]string{
		"${pb_pkg_name}":    "shopPb",
		"${model_pkg_name}": "shopModel",
		//		"${start_index}": "1",
	})
	fmt.Println(ps)
	fmt.Println(setM)
	fmt.Println(setP)
}
